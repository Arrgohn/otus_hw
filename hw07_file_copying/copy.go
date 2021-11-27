package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrNoOption              = errors.New("option not found")
	ErrLimitExcess           = errors.New("limit exceed")
	ErrWrongFileSize         = errors.New("wrong file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	err := checkOptions(fromPath, toPath)
	if err != nil {
		return err
	}

	fileFrom, err := open(fromPath)
	defer func(fileFrom *os.File) {
		err := fileFrom.Close()
		if err != nil {
			return
		}
	}(fileFrom)

	if err != nil {
		return err
	}

	fileTo, err := create(toPath)
	if err != nil {
		return errors.New("smth wrong")
	}
	defer func(fileTo *os.File) {
		err := fileTo.Close()
		if err != nil {
			return
		}
	}(fileTo)

	fileSize, err := getFileSize(fileFrom)
	if err != nil {
		return err
	}

	err = moveOffset(fileFrom, fileSize, offset)
	if err != nil {
		return err
	}

	err = startCopy(fileSize, fileFrom, fileTo, offset, limit)
	if err != nil {
		return err
	}

	return nil
}

func checkOptions(from, to string) error {
	if from == "" || to == "" {
		return ErrNoOption
	}

	return nil
}

func open(path string) (*os.File, error) {
	file, err := os.Open(path)

	if os.IsNotExist(err) {
		return nil, ErrUnsupportedFile
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

func moveOffset(fileFrom *os.File, fileSize int, offset int64) error {
	if offset != 0 {
		_, err := fileFrom.Seek(offset, io.SeekStart)
		if err != nil || int(offset) >= fileSize {
			return ErrOffsetExceedsFileSize
		}
	}

	return nil
}

func startCopy(fileSize int, fileFrom *os.File, fileTo *os.File, offset, limit int64) error {
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
