package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	commandString := cmd[0]
	command := exec.Command(commandString, cmd[1:]...)
	command.Stdout = os.Stdout
	processedEnv := getEnv(env)

	command.Env = append(os.Environ(), processedEnv...)
	if err := command.Run(); err != nil {
		println(err.Error())
		return 0
	}

	return 1
}

func getEnv(env Environment) []string {
	vals := make([]string, 0)
	for key, v := range env {
		if v.NeedRemove {
			err := os.Unsetenv(key)
			if err != nil {
				return nil
			}
		} else {
			vals = append(vals, key+"="+v.Value)
		}
	}

	return vals
}
