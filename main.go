package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-fsnotify/fsnotify"
)

var watcher *fsnotify.Watcher

// watchDir gets run as a walk func, searching for directories to add watchers to
func watchDir(path string, fi os.FileInfo, err error) error {

	// since fsnotify can watch all the files in a directory, watchers only need
	// to be added to each nested directory
	if fi.Mode().IsDir() {
		return watcher.Add(path)
	}

	return nil
}

func watchRemoveDir(path string, fi os.FileInfo, err error) error {

	// since fsnotify can watch all the files in a directory, watchers only need
	// to be added to each nested directory
	if fi.Mode().IsDir() {
		return watcher.Remove(path)
	}

	return nil
}

// main
func main() {

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	exPath := filepath.Dir(ex)
	fmt.Println(exPath)

	dirName := filepath.Base(exPath)
	fmt.Println("dirname", dirName)

	exeName := dirName + ".exe"

	// creates a new file watcher
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	defer watcher.Close()

	//fc := exec.Command("go", "build", "testexample/test.go")

	fc := exec.Command("go", "build")

	fc.Dir = exPath

	fc.Stdin = os.Stdin

	fc.Stdout = os.Stdout

	fc.Stderr = os.Stderr

	err = fc.Run()
	if err != nil {
		fmt.Println("ERROR: ", err)
		goto buildError
	}

	fc = exec.Command(exeName)

	fc.Dir = exPath

	fc.Stdin = os.Stdin

	fc.Stdout = os.Stdout

	fc.Stderr = os.Stderr

	fc.Start()

buildError:

	done := make(chan bool)

	restart := make(chan bool)

	go func() {
		for {
			select {

			case event := <-watcher.Events:
				fmt.Printf("EVENT! %#v\n", event)

				// watcher.Remove(exPath)

				// starting at the root of the project, walk each file/directory searching for
				// directories
				if err := filepath.Walk(exPath, watchRemoveDir); err != nil {
					fmt.Println("ERROR", err)
				}

				kill := exec.Command("TASKKILL", "/T", "/F", "/PID", strconv.Itoa(fc.Process.Pid))

				kill.Stdin = os.Stdin

				kill.Stdout = os.Stdout

				kill.Stderr = os.Stderr

				err = kill.Run()

				if err != nil {
					fmt.Println("error occurred when killing process")
				}

				time.Sleep(time.Second * 1)

				// fc = exec.Command("go", "build", "testexample/test.go")

				fc = exec.Command("go", "build")

				fc.Dir = exPath

				fc.Stdin = os.Stdin

				fc.Stdout = os.Stdout

				fc.Stderr = os.Stderr

				err = fc.Run()

				if err != nil {
					goto buildErrorInsideLabel
				}

				fc = exec.Command(exeName)

				fc.Dir = exPath

				fc.Stdin = os.Stdin

				fc.Stdout = os.Stdout

				fc.Stderr = os.Stderr

				fc.Start()

			buildErrorInsideLabel:

				time.Sleep(time.Second * 1)

				restart <- true

			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()

	go func() {
		for {
			select {
			case restarted := <-restart:
				if restarted {
					if err := filepath.Walk(exPath, watchDir); err != nil {
						fmt.Println("ERROR", err)
					}
				}
			}
		}
	}()

	// starting at the root of the project, walk each file/directory searching for
	// directories
	if err := filepath.Walk(exPath, watchDir); err != nil {
		fmt.Println("ERROR", err)
	}

	// out of the box fsnotify can watch a single file, or a single directory
	// if err := watcher.Add(exPath); err != nil {
	// 	fmt.Println("error for watcher: ", err)
	// }

	<-done
}
