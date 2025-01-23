# AndroidBox

# Modules

## setprop

原理分析请参考[AndroidBox-01-修改SystemProperty](https://picoorg.github.io/posts/androidbox-01-%E4%BF%AE%E6%94%B9systemproperty/)

```bash
# ./AndroidBox setprop name value
./AndroidBox setprop ro.debuggable 1
```

# libfuzzer

采用`CGO`交叉编译`libfuzzer`，参考`golang`[源码](https://github.com/golang/go/blob/608acff8479640b00c85371d91280b64f5ec9594/src/internal/platform/supported.go#L146)发现，`android/arm64`没有办法静态编译`-buildmode=c-archive`。

```bash
# adb push cmd/fuzz/fuzz_* /data/local/tmp/
cd /data/local/tmp/ && LD_LIBRARY_PATH=. ./fuzz_android33_arm64 > tmp.log
```
