package main

import (
	"errors"
	"github.com/cheggaaa/pb/v3"
	"io"
	"os"
	"path/filepath"
)

var ErrNoOption = errors.New("option not found")
var ErrNoFile = errors.New("file not found")
var ErrLimitExcess = errors.New("limit exceed")
var ErrWrongOffset = errors.New("wrong offset")
var ErrWrongFileSize = errors.New("wrong file size")

func Copy(fromPath, toPath string, offset, limit int64) error {
	err := checkOptions()
	if err != nil {
		return err
	}

	fileFrom, err := open(from)
	defer fileFrom.Close()
	if err != nil {
		return err
	}

	fileTo, err := create(to)
	if err != nil {
		return errors.New("smth wrong")
	}
	defer fileTo.Close()

	fileSize, err := getFileSize(fileFrom)
	if err != nil {
		return err
	}

	err = moveOffset(fileFrom, fileSize)
	if err != nil {
		return err
	}

	err = startCopy(fileSize, fileFrom, fileTo)
	if err != nil {
		return err
	}

	return nil
}

func checkOptions() error {
	if from == "" || to == ""{
		return ErrNoOption
	}

	return nil
}

func open(path string) (*os.File, error) {
	file, err := os.Open(path)

	if os.IsNotExist(err) {
		return nil, ErrNoFile
	}

	return file, nil
}

func create(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0776); err != nil {
		return nil, err
	}
	return os.Create(p)
}

func getFileSize(fileFrom *os.File) (int, error) {
	fileInfo, err := fileFrom.Stat()
	fileSize := fileInfo.Size()

	if err != nil && fileSize == 0 {
		return 0, ErrWrongFileSize
	}

	return int(fileSize), nil
}

func moveOffset(fileFrom *os.File, fileSize int) error {
	if offset != 0 {
		_, err := fileFrom.Seek(offset, io.SeekStart)
		if err != nil || int(offset) >= fileSize {
			return ErrWrongOffset
		}
	}

	return nil
}

func startCopy(fileSize int, fileFrom *os.File, fileTo *os.File) error {
	count := fileSize - int(offset)
	bar := pb.StartNew(count)
	writtenTotal := 0

	for i := 0; i < count; i++ {
		bar.Increment()

		written, err := io.CopyN(fileTo, fileFrom, 1)
		writtenTotal += int(written)

		if limit != 0 && writtenTotal > int(limit) {
			return ErrLimitExcess
		}

		if err != nil {
			return err
		}
	}

	bar.Finish()

	return nil
}

