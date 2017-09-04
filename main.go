package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/go-fsnotify/fsnotify"
)

func printError(err error) {
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
	}
}

func printOutput(outs []byte) {
	if len(outs) > 0 {
		fmt.Printf("==> Output: %s\n", string(outs))
	}
}

func kill(cmd *exec.Cmd) error {
	kill := exec.Command("TASKKILL", "/T", "/F", "/PID", strconv.Itoa(cmd.Process.Pid))
	kill.Stderr = os.Stderr
	kill.Stdout = os.Stdout
	return kill.Run()
}

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
		goto errorLabel
	}

	fc = exec.Command("test.exe")

	fc.Dir = "C:/gowork/src/github.com/ttufekci/ginmon"

	fc.Stdin = os.Stdin

	fc.Stdout = os.Stdout

	fc.Stderr = os.Stderr

	fc.Start()

errorLabel:

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
					// goto errorInsideLabel
				}

				time.Sleep(time.Second * 1)

				fc = exec.Command("cmd", "/C", "go", "build", "testexample/test.go")

				fc.Dir = "C:/gowork/src/github.com/ttufekci/ginmon"

				fmt.Println("testing")

				fc.Stdin = os.Stdin

				fc.Stdout = os.Stdout

				fc.Stderr = os.Stderr

				err = fc.Run()

				if err != nil {
					fmt.Println("error occurred 2")
					goto errorInsideLabel
				}

				fc = exec.Command("test.exe")

				fc.Dir = "C:/gowork/src/github.com/ttufekci/ginmon"

				fmt.Println("testing")

				fc.Stdin = os.Stdin

				fc.Stdout = os.Stdout

				fc.Stderr = os.Stderr

				fc.Start()

			errorInsideLabel:

				fmt.Println("deneme3")

				time.Sleep(time.Second * 1)

				fmt.Println("before restart true")

				restart <- true

				fmt.Println("restart is starting")

			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()

	go func() {
		for {
			select {
			case restarted := <-restart:
				fmt.Println("restarted another func", restarted)
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
		fmt.Println("ERROR_2", err)
	}

	<-done
}
