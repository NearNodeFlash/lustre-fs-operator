/*
 * Copyright 2021-2024 Hewlett Packard Enterprise Development LP
 * Other additional copyright holders may be indicated within.
 *
 * The entirety of this work is licensed under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 *
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package controller

import (
	"context"
	"os"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/api/resource"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/DataWorkflowServices/dws/utils/updater"
	lusv1beta1 "github.com/NearNodeFlash/lustre-fs-operator/api/v1beta1"
)

const (
	finalizerLustreFileSystem = "lus.cray.hpe.com/lustre_fs"
)

var (
	// Capacity, or Storage Resource Quantity, is required parameter and must be non-zero. This value is programmed into both the
	// Persistent Volume and Persistent Volume Claim, but remains unused by any of the Lustre CSI.
	persistentVolumeResourceQuantity = resource.MustParse("1")
)

// LustreFileSystemReconciler reconciles a LustreFileSystem object
type LustreFileSystemReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=lus.cray.hpe.com,resources=lustrefilesystems,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=lus.cray.hpe.com,resources=lustrefilesystems/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=lus.cray.hpe.com,resources=lustrefilesystems/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=persistentvolumes,verbs=get;list;update;create;patch;delete;watch
//+kubebuilder:rbac:groups=core,resources=persistentvolumeclaims,verbs=get;list;update;create;patch;delete;watch
//+kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the LustreFileSystem object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *LustreFileSystemReconciler) Reconcile(ctx context.Context, req ctrl.Request) (result ctrl.Result, err error) {
	fs := &lusv1beta1.LustreFileSystem{}
	if err := r.Get(ctx, req.NamespacedName, fs); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	statusUpdater := updater.NewStatusUpdater[*lusv1beta1.LustreFileSystemStatus](fs)
	defer func() { err = statusUpdater.CloseWithStatusUpdate(ctx, r.Client.Status(), err) }()

	// Check if the object is being deleted.
	if !fs.GetDeletionTimestamp().IsZero() {

		containsOtherFinalizers := func(o client.Object) bool {
			for _, f := range o.GetFinalizers() {
				if f != finalizerLustreFileSystem {
					return true
				}
			}
			return false
		}

		// Don't do anything until other finalizers are gone
		if containsOtherFinalizers(fs) {
			return ctrl.Result{}, nil
		}

		// Only remaining finalizer should be ours
		if !controllerutil.ContainsFinalizer(fs, finalizerLustreFileSystem) {
			return ctrl.Result{}, nil
		}

		for namespace := range fs.Spec.Namespaces {
			for _, mode := range fs.Spec.Namespaces[namespace].Modes {
				if err := r.deleteAccess(ctx, fs, namespace, mode); err != nil {
					return ctrl.Result{}, err
				}
			}
		}

		controllerutil.RemoveFinalizer(fs, finalizerLustreFileSystem)
		if err := r.Update(ctx, fs); err != nil {
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, nil
	}

	// Add the finalizer if it doesn't exist
	if !controllerutil.ContainsFinalizer(fs, finalizerLustreFileSystem) {
		controllerutil.AddFinalizer(fs, finalizerLustreFileSystem)
		if err := r.Update(ctx, fs); err != nil {
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, nil
	}

	// Iterate over the access modes in the specification. For each namespace in that mode
	// create a PV/PVC which can be used by pods in the same namespace.
	for namespace := range fs.Spec.Namespaces {
		namespacePresent := true

		// If the namespace doesn't exist, set a flag so that we can appropriately set the status in the mode loop
		ns := &corev1.Namespace{}
		if err := r.Get(ctx, types.NamespacedName{Name: namespace}, ns); err != nil {
			if errors.IsNotFound(err) {
				namespacePresent = false
			} else {
				return ctrl.Result{}, err
			}
		}

		// Create the Status Namespace map if empty
		if fs.Status.Namespaces == nil {
			fs.Status.Namespaces = make(map[string]lusv1beta1.LustreFileSystemNamespaceStatus)
		}

		// For each mode listed for the namespace
		for _, mode := range fs.Spec.Namespaces[namespace].Modes {
			// Create the Status Namespace Mode map if empty
			if fs.Status.Namespaces[namespace].Modes == nil {
				fs.Status.Namespaces[namespace] = lusv1beta1.LustreFileSystemNamespaceStatus{
					Modes: make(map[corev1.PersistentVolumeAccessMode]lusv1beta1.LustreFileSystemNamespaceAccessStatus),
				}
			}

			// Default the status as Pending in case the create/updates fail
			fs.Status.Namespaces[namespace].Modes[mode] = lusv1beta1.LustreFileSystemNamespaceAccessStatus{
				State: lusv1beta1.NamespaceAccessPending,
			}

			// If the namespace is not present or is not active, continue on and the status will be marked as Pending
			if !namespacePresent || ns.Status.Phase != corev1.NamespaceActive {
				continue
			}

			// Attempt to create the PV, if it fails, the status will be marked as Pending
			pv, err := r.createOrUpdatePersistentVolume(ctx, fs, namespace, mode)
			if err != nil {
				return ctrl.Result{}, err
			}

			// Attempt to create the PVC, if it fails, the status will be marked as Pending
			pvc, err := r.createOrUpdatePersistentVolumeClaim(ctx, fs, namespace, mode)
			if err != nil {
				return ctrl.Result{}, err
			}

			// If we got this far, the status is Ready
			fs.Status.Namespaces[namespace].Modes[mode] = lusv1beta1.LustreFileSystemNamespaceAccessStatus{
				State: lusv1beta1.NamespaceAccessReady,
				PersistentVolumeRef: &corev1.LocalObjectReference{
					Name: pv.Name,
				},
				PersistentVolumeClaimRef: &corev1.LocalObjectReference{
					Name: pvc.Name,
				},
			}
		}
	}

	// Remove any resources that are not in the spec
	for namespace := range fs.Status.Namespaces {
		for mode := range fs.Status.Namespaces[namespace].Modes {
			// Check if the provided namespace and mode are present in the specification
			isPresentInSpec := func(namespace string, mode corev1.PersistentVolumeAccessMode) bool {
				if ns, found := fs.Spec.Namespaces[namespace]; found {
					for _, m := range ns.Modes {
						if m == mode {
							return true
						}
					}
				}

				return false
			}

			if !isPresentInSpec(namespace, mode) {
				if err := r.deleteAccess(ctx, fs, namespace, mode); err != nil {
					return ctrl.Result{}, err
				}

				delete(fs.Status.Namespaces[namespace].Modes, mode)

				// Force a requeue because we just modified the modes in place
				return ctrl.Result{Requeue: true}, nil
			}
		}

		if _, found := fs.Spec.Namespaces[namespace]; !found {
			delete(fs.Status.Namespaces, namespace)

			// Force a requeue because we just modified the namespaces in place
			return ctrl.Result{Requeue: true}, nil
		}
	}

	return ctrl.Result{}, nil
}

func (r *LustreFileSystemReconciler) createOrUpdatePersistentVolumeClaim(ctx context.Context, fs *lusv1beta1.LustreFileSystem, namespace string, mode corev1.PersistentVolumeAccessMode) (*corev1.PersistentVolumeClaim, error) {

	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fs.PersistentVolumeClaimName(namespace, mode),
			Namespace: namespace,
		},
	}

	mutateFn := func() error {
		pvc.Spec.StorageClassName = &fs.Spec.StorageClassName
		pvc.Spec.VolumeName = fs.PersistentVolumeName(namespace, mode)

		pvc.Spec.AccessModes = []corev1.PersistentVolumeAccessMode{
			mode,
		}

		pvc.Spec.Resources = corev1.ResourceRequirements{
			Requests: corev1.ResourceList{
				corev1.ResourceStorage: persistentVolumeResourceQuantity,
			},
		}

		return nil
	}

	result, err := ctrl.CreateOrUpdate(ctx, r.Client, pvc, mutateFn)
	if err != nil {
		return nil, err
	}

	if result != controllerutil.OperationResultNone {
		log.FromContext(ctx).Info("PersistentVolumeClaim", "object", client.ObjectKeyFromObject(pvc).String(), "result", result)
	}

	return pvc, nil
}

func (r *LustreFileSystemReconciler) createOrUpdatePersistentVolume(ctx context.Context, fs *lusv1beta1.LustreFileSystem, namespace string, mode corev1.PersistentVolumeAccessMode) (*corev1.PersistentVolume, error) {

	pv := &corev1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: fs.PersistentVolumeName(namespace, mode),
		},
	}

	mutateFn := func() error {
		volumeMode := corev1.PersistentVolumeFilesystem
		pv.Spec.VolumeMode = &volumeMode

		pv.Spec.StorageClassName = fs.Spec.StorageClassName
		pv.Spec.AccessModes = []corev1.PersistentVolumeAccessMode{
			corev1.ReadWriteMany,
		}

		pv.Spec.Capacity = corev1.ResourceList{
			corev1.ResourceStorage: persistentVolumeResourceQuantity,
		}

		if pv.Spec.ClaimRef == nil {
			pv.Spec.ClaimRef = &corev1.ObjectReference{}
		} else if pv.Status.Phase == corev1.VolumeReleased {
			// The PV is being updated, and it was bound to a PVC
			// but now it's released.  Clear the uid of the
			// earlier PVC from the claimRef.
			// This allows it to bind to a new PVC.
			pv.Spec.ClaimRef.UID = ""
		}

		// Reserve this PV for the matching PVC.
		pv.Spec.ClaimRef.Name = fs.PersistentVolumeClaimName(namespace, mode)
		pv.Spec.ClaimRef.Namespace = namespace

		pv.Spec.PersistentVolumeSource = corev1.PersistentVolumeSource{
			CSI: &corev1.CSIPersistentVolumeSource{
				Driver:       os.Getenv("LUSTRE_CSI_SERVICE_NAME"),
				FSType:       "lustre",
				VolumeHandle: fs.Spec.MgsNids + ":/" + fs.Spec.Name,
			},
		}

		return nil
	}

	result, err := ctrl.CreateOrUpdate(ctx, r.Client, pv, mutateFn)
	if err != nil {
		return nil, err
	}

	if result != controllerutil.OperationResultNone {
		log.FromContext(ctx).Info("PersistentVolume", "object", client.ObjectKeyFromObject(pv).String(), "result", result)
	}

	return pv, nil
}

func (r *LustreFileSystemReconciler) deleteAccess(ctx context.Context, fs *lusv1beta1.LustreFileSystem, namespace string, mode corev1.PersistentVolumeAccessMode) error {
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fs.PersistentVolumeClaimName(namespace, mode),
			Namespace: namespace,
		},
	}

	log.FromContext(ctx).Info("Deleting PersistentVolumeClaim", "object", client.ObjectKeyFromObject(pvc).String())
	if err := r.Delete(ctx, pvc); err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	}

	pv := &corev1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: fs.PersistentVolumeName(namespace, mode),
		},
	}

	log.FromContext(ctx).Info("Deleting PersistentVolume", "object", client.ObjectKeyFromObject(pv).String())
	if err := r.Delete(ctx, pv); err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	}

	return nil
}

// For a change with any namespace, we need to go through the list of all the lustrefilesystem resources
func (r *LustreFileSystemReconciler) getLustreFileSystemsHandler(ctx context.Context, o client.Object) []reconcile.Request {
	var res []reconcile.Request

	filesystems := &lusv1beta1.LustreFileSystemList{}
	if err := r.List(ctx, filesystems); err != nil && !meta.IsNoMatchError(err) {
		return res
	}

	for _, lustre := range filesystems.Items {
		res = append(res, reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      lustre.GetName(),
				Namespace: lustre.GetNamespace(),
			},
		})
	}

	return res
}

// SetupWithManager sets up the controller with the Manager.
func (r *LustreFileSystemReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&lusv1beta1.LustreFileSystem{}).
		Watches(
			// Watch all namespaces for changes to ensure lustrefilesystem resources stay current
			&corev1.Namespace{}, handler.EnqueueRequestsFromMapFunc(r.getLustreFileSystemsHandler),
		).
		Complete(r)
}
