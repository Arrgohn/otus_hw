package main

import (
	"flag"
	"fmt"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()

	err := checkOptions()
	if err != nil {
		fmt.Println("No option error")
		return
	}

	fileFrom, err := open(from)
	defer fileFrom.Close()
	if err != nil {
		fmt.Println("No file error")
		return
	}

	fileTo, err := create(to)
	if err != nil {
		//return 0, err
	}
	defer fileTo.Close()

	fileSize, err := getFileSize(fileFrom)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = moveOffset(fileFrom, fileSize)
	if err != nil {
		fmt.Println("offset error")
		return
	}

	err = startCopy(int(fileSize), fileFrom, fileTo)
	if  err != nil {
		fmt.Println("smth wrong")
		return
	}

	fmt.Println("Success")
}
