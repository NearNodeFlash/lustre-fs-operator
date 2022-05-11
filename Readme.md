# Lustre File System Operator

## Contributing

Before opening an issue or pull request, please read the [Contributing] guide.

[contributing]: CONTRIBUTING.md

## Bootstrapping

This operator was boostrapped using the operator-sdk

```bash
operator-sdk init --domain cray.hpe.com --repo github.hpe.com/hpe/hpc-rabsw-lustre-fs-operator
operator-sdk create api --version v1alpha1 --kind LustreFileSystem --resource --controller
operator-sdk create webhook --version v1alpha1 --kind LustreFileSystem --programmatic-validation
```
