package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/go-fsnotify/fsnotify"
)

var watcher *fsnotify.Watcher

var ext = ".go"

func watchDir(path string, fi os.FileInfo, err error) error {

	if !fi.Mode().IsDir() {
		r, err := regexp.MatchString(ext, fi.Name())
		if err == nil && r {
			return watcher.Add(path)
		}
	}

	return nil
}

func watchRemoveDir(path string, fi os.FileInfo, err error) error {

	if !fi.Mode().IsDir() {
		r, err := regexp.MatchString(ext, fi.Name())
		if err == nil && r {
			return watcher.Remove(path)
		}
	}

	return nil
}

func main() {

	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	exPath := pwd

	dirName := filepath.Base(exPath)

	exeName := dirName + ".exe"

	// creates a new file watcher
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("ERROR: ", err)
	}

	defer watcher.Close()

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

	watcherRemoved := false

	go func() {
		for {
			select {

			case event := <-watcher.Events:
				watcher.Remove(event.Name)
				watcher.Add(event.Name)

				if watcherRemoved {
					fmt.Println("watcher removed")
					return
				}

				fmt.Printf("EVENT! %#v\n", event)

				watcherRemoved = true

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

				time.Sleep(time.Second * 2)

				restart <- true

				watcherRemoved = false

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

	if err := filepath.Walk(exPath, watchDir); err != nil {
		fmt.Println("ERROR", err)
	}

	<-done
}
