package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func mkPipes() {
	syscall.Mkfifo("/tmp/runtime-output", 0600)
	syscall.Mkfifo("/tmp/runtime-input", 0600)
}

func getInvocation(runtimeAPI string) (string, []byte) {
	resp, err := http.Get(fmt.Sprintf("%s/invocation/next", runtimeAPI))
	if nil != err {
		log.Fatalf("Error getting next invocation: %s", err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		log.Fatalf("Error reading invocation body: %s", err.Error())
	}
	invocationID := resp.Header["Lambda-Runtime-Aws-Request-Id"][0]
	log.Printf("EVENT open-runtime %s %s", invocationID, body)
	return invocationID, body
}

func runMiddlewares(middlewares []string, hook string, body []byte) []byte {
	for _, middleware := range middlewares {
		body = runMiddleware(middleware, hook, body)
	}
	return body
}

func runMiddleware(name string, hook string, body []byte) []byte {
	cmd := exec.Command(fmt.Sprintf("/opt/middlewares/%s", name), hook)
	cmd.Stderr = os.Stderr
	stdin, err := cmd.StdinPipe()
	if nil != err {
		log.Fatalf("Error obtaining stdin: %s", err.Error())
	}
	stdin.Write(body)
	stdin.Close()
	result, err := cmd.Output()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	result = bytes.TrimSuffix(result, []byte("\n")) // middlewares can easilly add newlines by accident
	return result
}

func main() {
	runtimeAPI := fmt.Sprintf("http://%s/2018-06-01/runtime", os.Getenv("AWS_LAMBDA_RUNTIME_API"))
	middlewaresString := os.Getenv("SLSMIDDLEWARES")
	var middlewares []string
	if len(middlewaresString) > 0 {
		middlewares = strings.Split(os.Getenv("SLSMIDDLEWARES"), ",")
	}

	mkPipes()

	// start language runtime
	cmd := exec.Command("/opt/language-runtime", os.Getenv("_HANDLER"))
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Start(); nil != err {
		log.Fatalf("Error starting program: %s, %s", cmd.Path, err.Error())
	}

	// open pipes
	input, err := os.OpenFile("/tmp/runtime-input", os.O_RDWR, 0600)
	if nil != err {
		log.Fatalf("Error obtaining input: %s", err.Error())
	}
	output, err := os.OpenFile("/tmp/runtime-output", os.O_RDONLY, 0600)
	if nil != err {
		log.Fatalf("Error obtaining output: %s", err.Error())
	}

	reader := bufio.NewReader(output)

	go func(reader io.Reader) {
		scanner := bufio.NewScanner(reader)
		invocationID, body := getInvocation(runtimeAPI)
		body = runMiddlewares(middlewares, "before", body)
		input.Write(body)
		io.WriteString(input, "\n")
		for scanner.Scan() {
			response := scanner.Bytes()
			response = runMiddlewares(middlewares, "after", response)
			log.Printf("Response! %s", invocationID)
			log.Printf("%s/invocation/%s/response", runtimeAPI, invocationID)
			_, err := http.Post(
				fmt.Sprintf("%s/invocation/%s/response", runtimeAPI, invocationID),
				"application/json",
				bytes.NewBuffer(response),
			)
			if nil != err {
				log.Fatalf("Error POSTing response")
			}
			log.Print("done sending response")
			invocationID, body = getInvocation(runtimeAPI)
			body = runMiddlewares(middlewares, "before", body)
			input.Write(body)
			io.WriteString(input, "\n")
		}
	}(reader)

	cmd.Wait()
}
