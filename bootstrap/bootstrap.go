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
)

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
	return result
}

func main() {
	runtimeAPI := fmt.Sprintf("http://%s/2018-06-01/runtime", os.Getenv("AWS_LAMBDA_RUNTIME_API"))
	middlewares := strings.Split(os.Getenv("SLSMIDDLEWARES"), ",")

	cmd := exec.Command("/opt/language-runtime", os.Getenv("_HANDLER"))
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

	go func(reader io.Reader) {
		scanner := bufio.NewScanner(reader)
		invocationID, body := getInvocation(runtimeAPI)
		body = runMiddlewares(middlewares, "before", body)
		stdin.Write(body)
		io.WriteString(stdin, "\n")
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
			stdin.Write(body)
			io.WriteString(stdin, "\n")
		}
	}(reader)

	if err := cmd.Start(); nil != err {
		log.Fatalf("Error starting program: %s, %s", cmd.Path, err.Error())
	}

	cmd.Wait()
}
