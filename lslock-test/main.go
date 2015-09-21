package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"syscall"
	"time"

	"code.google.com/p/go-uuid/uuid"
)

var filesFlag int
var directoryFlag string

func init() {
	flag.IntVar(&filesFlag, "f", 0, "number of files to create")
	flag.StringVar(&directoryFlag, "d", "", "directory to create the new files")
}

func cleanUp() {
	fmt.Println("Removing files...")
	files, err := ioutil.ReadDir(directoryFlag)
	if err != nil {
		log.Fatalln(err)
	}
	for file := range files {
		if err := os.Remove(files[file].Name()); err != nil {
			log.Fatalln(err)
		}
	}
}

func main() {
	os.Exit(realMain())
}

func realMain() int {
	flag.Parse()

	if filesFlag == 0 || directoryFlag == "" {
		flag.Usage()
		return 1
	}

	fmt.Println("Creating files...")
	var wg sync.WaitGroup
	for i := 1; i <= filesFlag; i++ {
		wg.Add(1)
		go func() {
			filename := fmt.Sprintf("%s/%s.lock", directoryFlag, uuid.NewUUID().String())
			err := ioutil.WriteFile(filename, []byte("foo"), 0644)
			if err != nil {
				log.Fatalln(err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	defer cleanUp()

	fmt.Println("Getting all files in " + directoryFlag + "...")
	files, err := ioutil.ReadDir(directoryFlag)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Opening all files and applying locks")
	for i, file := range files {
		fi, err := os.Open(directoryFlag + "/" + file.Name())
		if err != nil {
			log.Fatalln(err)
		}
		defer fi.Close()
		fd := fi.Fd()
		if i%2 == 0 {
			syscall.Flock(int(fd), syscall.LOCK_EX|syscall.LOCK_NB)
		} else {
			syscall.Flock(int(fd), syscall.LOCK_SH)
		}
	}

	time.Sleep(time.Second * 3000)

	return 0
}
