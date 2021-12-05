package main

import "os"

func main() {
	directory := os.Args[1]

	cmd := make([]string, 0)
	cmd = append(cmd, os.Args[2:]...)

	env, err := ReadDir(directory)
	if err != nil {
		return
	}

	RunCmd(cmd, env)
}
