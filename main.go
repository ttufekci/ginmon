package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/go-fsnotify/fsnotify"
)

func printCommand(cmd *exec.Cmd) {
	fmt.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
}

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

	fc := exec.Command("cmd", "/C", "go", "run", "testexample/test.go")

	fmt.Println("testing")

	fc.Stdin = os.Stdin

	fc.Stdout = os.Stdout

	fc.Stderr = os.Stderr

	fc.Start()

	fmt.Println("deneme")

	// fmt.Printf("testing", out.String())

	//
	done := make(chan bool)

	//
	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				fmt.Printf("EVENT! %#v\n", event)

				kill := exec.Command("TASKKILL", "/T", "/F", "/PID", strconv.Itoa(fc.Process.Pid))

				kill.Stdin = os.Stdin

				kill.Stdout = os.Stdout

				kill.Stderr = os.Stderr

				kill.Run()

				time.Sleep(time.Second * 3)

				fc := exec.Command("cmd", "/C", "go", "run", "testexample/test.go")

				fmt.Println("testing")

				fc.Stdin = os.Stdin

				fc.Stdout = os.Stdout

				fc.Stderr = os.Stderr

				fc.Start()

				fmt.Println("deneme3")

				// c := exec.Command("cmd", "/C", "go", "run", event.Name)

				// c.Stdin = os.Stdin
				// c.Stdout = os.Stdout
				// c.Stderr = os.Stderr

				// c.Run()
			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()

	// out of the box fsnotify can watch a single file, or a single directory
	if err := watcher.Add("testexample"); err != nil {
		fmt.Println("ERROR", err)
	}

	<-done
}
