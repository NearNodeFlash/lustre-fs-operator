# Lustre File System Operator

## Contributing

Before opening an issue or pull request, please read the [Contributing] guide.

[contributing]: CONTRIBUTING.md

## Installing

Install on a Rabbit system with:
```console
kubectl apply -k 'https://github.com/NearNodeFlash/lustre-fs-operator.git/config/default/?ref=master'
```

Install a specific release with:
```console
kubectl apply -k 'https://github.com/NearNodeFlash/lustre-fs-operator.git/config/default/?ref=v0.0.1'
```

## Making a release

When making a release, set the image tag in the kustomization configuration
with the following and commit the new kustomization.yaml to the release branch.
This will allow the above `kubectl apply -k` commands to work with that
release, ensuring that when installed, the pods will pull the same container
tag.

Note that for a release of "v0.0.2", we specify it without the leading "v":
```console
make installer VERSION=0.0.2
```

The `master` branch should always point to the `master` tag.  This can always be reset with:
```console
make installer VERSION=master
```

