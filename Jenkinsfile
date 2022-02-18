@Library('dst-shared@master') _

// See https://github.hpe.com/hpe/hpc-dst-jenkins-shared-library for all
// the inputs to the dockerBuildPipeline.
// In particular: vars/dockerBuildPipeline.groovy
dockerBuildPipeline {
        repository = "cray"
        imagePrefix = "cray"
        app = "dp-lustre-fs-operator"
        name = "dp-lustre-fs-operator"
        description = "Operator for global lustre filesystem description"
        dockerfile = "Dockerfile"
        autoJira = false
        createSDPManifest = false
        product = "rabsw"
}
