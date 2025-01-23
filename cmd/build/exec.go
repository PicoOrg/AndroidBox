package main

import (
	"fmt"
	"os"
	"os/exec"
)

type Exec interface {
	Run(name string, arg ...string) (err error)
}

type execModel struct {
	cfg       *config
	toolchain Toolchain
}

func NewExec(cfg *config, toolchain Toolchain) Exec {
	return &execModel{
		cfg:       cfg,
		toolchain: toolchain,
	}
}

func (instance *execModel) Run(name string, arg ...string) (err error) {
	cmd := exec.Command(name, arg...)
	cmd.Dir = instance.cfg.BuildPath
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Env = append(cmd.Env, instance.toolchain.GetEnv()...)
	message, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("[error]", err, "[cmd]", cmd.String(), "[message]", string(message))
	} else {
		fmt.Println("[success]", "[cmd]", cmd.String(), "[message]", string(message))
	}
	return
}
