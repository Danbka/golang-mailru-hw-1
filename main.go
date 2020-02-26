package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
)

func getSize(file os.FileInfo) string {
	var fileSize string

	if file.Size() == 0 {
		fileSize = "empty"
	} else {
		fileSize = strconv.FormatInt(file.Size(), 10) + "b"
	}

	return fileSize
}

func dirLevel(out io.Writer, path string, printFiles bool, level int, prefix string) error {

	// содержимое текущей директории
	dirContent, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err.Error())
	}

	var directories []os.FileInfo

	// выкинуть файлы, если необходимо
	for _, file := range dirContent {
		if !printFiles && !file.IsDir() {
			continue
		}
		directories = append(directories, file)
	}

	level++

	for i, file := range directories {
		fileName := file.Name()
		if !file.IsDir() {
			fileName += " (" + getSize(file) + ")"
		}

		if i == len(directories)-1 {
			_, err = fmt.Fprintln(out, prefix+"└───"+fileName)
		} else {
			_, err = fmt.Fprintln(out, prefix+"├───"+fileName)
		}

		if err != nil {
			return fmt.Errorf(err.Error())
		}

		currentPrefix := prefix

		if file.IsDir() {
			if i != len(directories)-1 {
				currentPrefix = currentPrefix + "│	"
			} else {
				currentPrefix = currentPrefix + "	"
			}
			err = dirLevel(out, path+string(os.PathSeparator)+file.Name(), printFiles, level, currentPrefix)
			if err != nil {
				return fmt.Errorf(err.Error())
			}
		}
	}

	return nil
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	return dirLevel(out, path, printFiles, -1, "")
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
