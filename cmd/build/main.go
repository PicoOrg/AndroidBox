package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func ndkRoot() (string, error) {
	ndkPaths := []string{"NDK", "NDK_HOME", "NDK_ROOT", "ANDROID_NDK_HOME"}
	ndkRoot := ""
	for _, path := range ndkPaths {
		ndkRoot = os.Getenv(path)
		if len(ndkRoot) > 0 {
			if _, err := os.Stat(ndkRoot); err == nil {
				return ndkRoot, nil
			}
		}
	}

	return "", fmt.Errorf("no Android NDK found in $ANDROID_HOME/ndk-bundle, $ANDROID_HOME/ndk, $NDK_HOME, $NDK_ROOT nor in $ANDROID_NDK_HOME")
}

var (
	config_path string
	tag         string
)

func main() {
	flag.StringVar(&config_path, "f", "./config.yaml", "config file path")
	flag.StringVar(&tag, "t", "", "tag")
	flag.Parse()
	if tag != "" {
		tag = fmt.Sprintf(`_%s`, tag)
	}
	// fmt.Println(config_path)
	viper.SetConfigFile(config_path)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	// fmt.Println(os.ReadFile(config_path))
	cfg := &config{}
	if err := viper.Unmarshal(cfg); err != nil {
		panic(err)
	}
	// fmt.Println(cfg)
	if cfg.BuildPath == "" {
		panic(`build_path is empty`)
	}

	if cfg.NdkPath == "" {
		if ndkRoot, err := ndkRoot(); err == nil {
			cfg.NdkPath = ndkRoot
		} else {
			panic(err)
		}
	}

	for _, level := range cfg.ApiLevel {
		for _, arch := range cfg.Arch {
			toolchain := NewNdkToolchain(cfg.NdkPath, arch, level)
			e := NewExec(cfg, toolchain)
			filename := fmt.Sprintf("%s%s_android%s_%s", filepath.Base(cfg.BuildPath), tag, level, arch)
			if cfg.Fuzz {
				err := e.Run("go", "build", "-buildmode=c-shared", "-o", filename+".so", ".")
				if err != nil {
					panic(err)
				}
				_, clangpp, err := toolchain.ClangPath()
				if err != nil {
					panic(err)
				}
				err = e.Run(clangpp, "-g", "-fsanitize=fuzzer,address", "-static-libstdc++", "-Wimplicit-function-declaration", filename+".so", "-o", filename)
				if err != nil {
					panic(err)
				}
			} else {
				err := e.Run("go", "build", "-o", filename, ".")
				if err != nil {
					panic(err)
				}
			}
		}
	}
}
