package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, errors.New("wrong directory")
	}

	env := make(Environment)

	for _, fileInfo := range files {
		buf := make([]byte, fileInfo.Size())

		file, err := open(dir + "/" + fileInfo.Name())
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				return
			}
		}(file)
		if err != nil {
			return env, err
		}

		_, err = file.Read(buf)
		if err != nil {
			return nil, err
		}

		line := extractFirstLine(string(buf))

		env[fileInfo.Name()] = EnvValue{
			line,
			len(line) == 0,
		}
	}

	return env, err
}

func open(path string) (*os.File, error) {
	file, err := os.Open(path)

	if os.IsNotExist(err) {
		return nil, errors.New("file not found")
	}

	return file, nil
}

func extractFirstLine(content string) string {
	lines := strings.Split(content, "\n")

	firstLine := strings.TrimRight(lines[0], " ")

	b := bytes.ReplaceAll([]byte(firstLine), []byte{0x00}, []byte("\n"))

	return string(b)
}
