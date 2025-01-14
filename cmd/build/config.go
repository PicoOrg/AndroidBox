package main

type config struct {
	BuildPath string   `yaml:"build_path" mapstructure:"build_path"`
	NdkPath   string   `yaml:"ndk_path" mapstructure:"ndk_path"`
	ApiLevel  []string `yaml:"api_level" mapstructure:"api_level"`
	Arch      []string `yaml:"arch" mapstructure:"arch"`
}
