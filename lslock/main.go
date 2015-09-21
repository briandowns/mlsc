package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"
)

const (
	// ProcFSLocks holds the path to the /proc/locks file
	ProcFSLocks = "/proc/locks"

	// LsLockDir holds the path to the lslock directory
	LsLockDir = "/tmp/lslock-test"
)

var directoryFlag string

func init() {
	flag.StringVar(&directoryFlag, "d", "", "directory to search through. Hint: "+LsLockDir)
}

// parseLocks reads in all locks and returns them in a string slice
func parseLocks(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func main() {
	os.Exit(realMain())
}

func realMain() int {
	flag.Parse()

	if directoryFlag == "" {
		flag.Usage()
		os.Exit(1)
	}

	lines, err := parseLocks(ProcFSLocks)
	if err != nil {
		log.Fatalln(err)
	}

	files, err := ioutil.ReadDir(directoryFlag)
	if err != nil {
		log.Fatalln(err)
	}

	for file := range files {
		fData, err := os.Lstat(directoryFlag + "/" + files[file].Name())
		if err != nil {
			log.Fatalln(err)
		}

		//fmt.Println(fData.Name(), fData.Size(), fData.ModTime(), fData.IsDir(), fData.Sys())
		inode := int(fData.Sys().(*syscall.Stat_t).Ino)

		for line := range lines {
			lSplit := strings.Fields(lines[line])
			fdData := lSplit[5]
			fInode := strings.Split(fdData, ":")
			ind, err := strconv.Atoi(fInode[2])
			if err != nil {
				log.Fatalln(err)
			}
			if ind == int(inode) {
				fmt.Printf("Path: %s/%s, PID: %s\n", directoryFlag, fData.Name(), lSplit[4])
			}
		}
	}
	return 0
}
