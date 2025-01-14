package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Toolchain interface {
	GetEnv() (env []string)
}

type ndkToolchain struct {
	ndkPath     string
	level       string
	arch        string
	abi         string
	toolPrefix  string
	clangPrefix string
}

func NewNdkToolchain(ndkPath, arch, level string) Toolchain {
	tc := &ndkToolchain{
		ndkPath: ndkPath,
		level:   level,
		arch:    arch,
	}
	switch arch {
	case "arm":
		tc.abi = "armeabi-v7a"
		tc.toolPrefix = "arm-linux-androideabi"
		tc.clangPrefix = "armv7a-linux-androideabi"
	case "arm64":
		tc.abi = "arm64-v8a"
		tc.toolPrefix = "aarch64-linux-android"
		tc.clangPrefix = "aarch64-linux-android"
	case "386":
		tc.abi = "x86"
		tc.toolPrefix = "i686-linux-android"
		tc.clangPrefix = "i686-linux-android"
	case "amd64":
		tc.abi = "x86_64"
		tc.toolPrefix = "x86_64-linux-android"
		tc.clangPrefix = "x86_64-linux-android"
	default:
		panic(`unsupported architecture: ` + arch)
	}
	return tc
}

func (tc *ndkToolchain) GetEnv() (env []string) {
	env = make([]string, 0)
	clang, clangpp, err := tc.clangPath()
	if err != nil {
		panic(fmt.Sprintf("no compiler for was found in the NDK. %v", err))
	}
	env = append(env, []string{
		"CGO_CFLAGS=-I" + tc.includePath(),
		"CGO_LDFLAGS=-L" + tc.libraryPath(),
		"GOOS=android",
		"GOARCH=" + tc.arch,
		"CC=" + clang,
		"CXX=" + clangpp,
		"CGO_ENABLED=1",
	}...)
	if tc.arch == "arm" {
		env = append(env, "GOARM=7")
	}
	return
}

func (tc *ndkToolchain) clangPath() (clang, clangpp string, err error) {
	binPath := tc.bin()

	entries, err := os.ReadDir(binPath)
	if err != nil {
		return "", "", err
	}

	clangPrefix := fmt.Sprintf("%s%s-clang", tc.clangPrefix, tc.level)
	clangppPrefix := fmt.Sprintf("%s%s-clang++", tc.clangPrefix, tc.level)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if strings.HasPrefix(entry.Name(), clangPrefix) && !strings.HasPrefix(entry.Name(), clangppPrefix) {
			clang = filepath.Join(binPath, entry.Name())
		}

		if strings.HasPrefix(entry.Name(), clangppPrefix) {
			clangpp = filepath.Join(binPath, entry.Name())
		}
	}

	if clang == "" || clangpp == "" {
		return "", "", errors.New("can't find clang or clang++")
	}
	return clang, clangpp, nil
}

func (tc *ndkToolchain) bin() string {
	return filepath.Join(tc.ndkPath, "toolchains", "llvm", "prebuilt", tc.archNDK(), "bin")
}

func (tc *ndkToolchain) includePath() string {
	return filepath.Join(tc.ndkPath, "toolchains", "llvm", "prebuilt", tc.archNDK(), "sysroot", "usr", "include")
}

func (tc *ndkToolchain) libraryPath() string {
	return filepath.Join(tc.ndkPath, "toolchains", "llvm", "prebuilt", tc.archNDK(), "sysroot",
		"usr", "lib", tc.toolPrefix, tc.level)
}
func (tc *ndkToolchain) archNDK() string {
	if runtime.GOOS == "windows" && runtime.GOARCH == "386" {
		return "windows"
	} else if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" {
		return runtime.GOOS + "-" + "x86_64"
	} else {
		var arch string
		switch runtime.GOARCH {
		case "386":
			arch = "x86"
		case "amd64":
			arch = "x86_64"
		default:
			panic("unsupported GOARCH: " + runtime.GOARCH)
		}
		return runtime.GOOS + "-" + arch
	}
}
