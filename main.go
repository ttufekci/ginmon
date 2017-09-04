package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/go-fsnotify/fsnotify"
)

// main
func main() {

	// creates a new file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("ERROR", err)
	}
	defer watcher.Close()

	fc := exec.Command("cmd", "/C", "go", "build", "testexample/test.go")

	fc.Dir = "C:/gowork/src/github.com/ttufekci/ginmon"

	fc.Stdin = os.Stdin

	fc.Stdout = os.Stdout

	fc.Stderr = os.Stderr

	err = fc.Run()
	if err != nil {
		fmt.Println("error occurred")
		goto buildError
	}

	fc = exec.Command("test.exe")

	fc.Dir = "C:/gowork/src/github.com/ttufekci/ginmon"

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

				watcher.Remove("testexample")

				kill := exec.Command("TASKKILL", "/T", "/F", "/PID", strconv.Itoa(fc.Process.Pid))

				kill.Stdin = os.Stdin

				kill.Stdout = os.Stdout

				kill.Stderr = os.Stderr

				err = kill.Run()

				if err != nil {
					fmt.Println("error occurred when killing process")
				}

				time.Sleep(time.Second * 1)

				fc = exec.Command("cmd", "/C", "go", "build", "testexample/test.go")

				fc.Dir = "C:/gowork/src/github.com/ttufekci/ginmon"

				fc.Stdin = os.Stdin

				fc.Stdout = os.Stdout

				fc.Stderr = os.Stderr

				err = fc.Run()

				if err != nil {
					goto buildErrorInsideLabel
				}

				fc = exec.Command("test.exe")

				fc.Dir = "C:/gowork/src/github.com/ttufekci/ginmon"

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
					if err := watcher.Add("testexample"); err != nil {
						fmt.Println("ERROR", err)
					}
				}
			}
		}
	}()

	// out of the box fsnotify can watch a single file, or a single directory
	if err := watcher.Add("testexample"); err != nil {
		fmt.Println("error for watcher: ", err)
	}

	<-done
}
