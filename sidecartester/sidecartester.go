package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"os/exec"
)

func main() {

	cmd := exec.Command(os.Args[1], os.Args[2], os.Args[3])
	cmd.Stderr = os.Stderr
	stdin, err := cmd.StdinPipe()
	if nil != err {
		log.Fatalf("Error obtaining stdin: %s", err.Error())
	}
	stdout, err := cmd.StdoutPipe()
	if nil != err {
		log.Fatalf("Error obtaining stdout: %s", err.Error())
	}
	reader := bufio.NewReader(stdout)
	io.WriteString(stdin, "{\"asdf\": 53423}\n")
	stdin.Close()
	go func(reader io.Reader) {
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			log.Printf("result from sidecar: %s", scanner.Text())
		}
	}(reader)
	if err := cmd.Start(); nil != err {
		log.Fatalf("Error starting program: %s, %s", cmd.Path, err.Error())
	}
	cmd.Wait()

}
