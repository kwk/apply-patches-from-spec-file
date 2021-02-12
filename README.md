# Purpose

At Red Hat the big [LLVM project](https://github.com/llvm/llvm-project) is packaged into multiple sub-packages (e.g. [`clang`](https://src.fedoraproject.org/rpms/clang.git), [`lld`](https://src.fedoraproject.org/rpms/lld.git), [`lldb`](https://src.fedoraproject.org/rpms/lldb.git)) and before the actual build, the sourecode is patched.

The build as well as the patch process is done with `.spec` files.

Every so often it happens that a patch can no longer be applied to the source code.

Building LLVM can take quite a long time and it is always annoying to find out at some point down the line that a patch cannot be applied. That's where this project comes in handy.

# Usage:

```bash
go run main.go <path-to.spec> <exec-dir>
```

# Example usage for clang

Let's verify that all the packages we maintain in Fedora rawhide for clang can still be applied to the latest LLVM master

We need

1. this project,
1. the rpm repo
1. and the LLVM source tree,

so let's check them out into their own directories. It doesn't matter where you do this.

```bash
git clone -b main https://github.com/kwk/apply-patches-from-spec-file.git ~/apply-patches-from-spec-file

git clone -b rawhide https://src.fedoraproject.org/rpms/clang.git ~/rpm-clang

git clone -b main https://github.com/llvm/llvm-project.git ~/llvm-project
```

Let's apply all the patches

```bash
go run ~/apply-patches-from-spec-file/main.go ~/rpm-clang/clang.spec ~/llvm-project/clang
```



