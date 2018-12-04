package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func main() {

	syscall.Mkfifo("/tmp/runtime-output", 0600)
	syscall.Mkfifo("/tmp/runtime-input", 0600)
	cmd := exec.Command(os.Args[1], os.Args[2], os.Args[3])
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Start(); nil != err {
		log.Fatalf("Error starting program: %s, %s", cmd.Path, err.Error())
	}
	input, err := os.OpenFile("/tmp/runtime-input", os.O_RDWR, 0600)
	if nil != err {
		log.Fatalf("Error obtaining input: %s", err.Error())
	}
	output, err := os.OpenFile("/tmp/runtime-output", os.O_RDONLY, 0600)
	if nil != err {
		log.Fatalf("Error obtaining output: %s", err.Error())
	}
	io.WriteString(input, "{\"asdf\": 53423}\n")
	input.Close()
	reader := bufio.NewReader(output)
	go func(reader io.Reader) {
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			log.Printf("result from sidecar: %s", scanner.Text())
		}
	}(reader)
	cmd.Wait()

}
