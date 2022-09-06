package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {
	if fileList, err := execute(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Copied the following files:")
		for _, path := range fileList {
			fmt.Println(path)
		}
		fmt.Println()
		fmt.Printf("%d files have been copied.", len(fileList))
	}
}

func execute() ([]string, error) {
	flagValues := getFlagValues()

	tryCreateOutputDir(flagValues.dest)

	fileList := []string{}
	getAllFilesWithExtension(flagValues.src, flagValues.exts, func(value string) {
		go func() {
			fileList = append(fileList, value)
			copyFile(value, filepath.Join(flagValues.dest, filepath.Base(value)))
		}()
	})
	return fileList, nil
}

type FlagValues struct {
	src  string
	dest string
	exts []string
}

// Parse the values behind the flags and returns a struct
func getFlagValues() FlagValues {
	srcDir := flag.String("src", ".", "Source directory")
	destDir := flag.String("out", "./output", "Destination directory")
	ext := flag.String("exts", "", "List of extensions example: -ext 'pdf' OR -ext 'pdf png jpg'")
	flag.Parse()

	// Setting default values
	if *srcDir == "" {
		*srcDir = "."
	}
	if *destDir == "" {
		*destDir = "./output"
	}

	if *ext == "" {
		fmt.Fprintf(os.Stdout, "\033[0;31m%s\033[0m", "No extensions found! Please provide an extension to look for.")
		flag.Usage()
		os.Exit(-1)
	}

	extList := []string{}
	extList = append(extList, *ext)
	extList = append(extList, flag.Args()...)

	preparedExtList := getExtensionList(extList)

	return FlagValues{src: *srcDir, dest: *destDir, exts: preparedExtList}
}

// Returns prepared extension list in format ".<ext>" like ".pdf"
func getExtensionList(exts []string) []string {
	extensionList := []string{}
	for _, item := range exts {
		extensionList = append(extensionList, "."+item)
	}
	return extensionList
}

// Get all file paths with given extension in folder
func getAllFilesWithExtension(srcDir string, extsList []string, cb func(from string)) {
	err := filepath.Walk(srcDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			ext := filepath.Ext(path)
			if Contains(extsList, ext) {
				cb(path)
			}
			return nil
		})
	if err != nil {
		log.Fatalf("%s", err)
	}
}

// Creates folder if not exists
func tryCreateOutputDir(dirPath string) {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		os.MkdirAll(dirPath, os.ModeDir)
	}
}

// Copies file from source to destination path
func copyFile(src string, dest string) {
	newDest := getNewFilePath(dest)
	input, err := ioutil.ReadFile(src)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ioutil.WriteFile(newDest, input, 0644)
	if err != nil {
		fmt.Println("Error creating", newDest)
		fmt.Println(err)
		return
	}
}

// Checks if the destination path exists
// if so it appends "_copy_<UNIX_TIMESTAMP>" at the end of file
func getNewFilePath(dest string) string {
	if _, err := os.Stat(dest); errors.Is(err, os.ErrNotExist) {
		return dest
	} else {
		ext := filepath.Ext(dest)
		bufferDest := strings.Replace(dest, ext, "", -1)
		return (bufferDest + "_copy_" + strconv.FormatInt(time.Now().Unix(), 10) + ext)
	}
}

// Checks if a given element is in slice
// Returns true if its in slice
func Contains(s []string, e string) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}
