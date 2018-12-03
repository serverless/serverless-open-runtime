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
)

func getInvocation(runtime_api string) (string, []byte) {
	resp, _ := http.Get(fmt.Sprintf("%s/invocation/next", runtime_api))
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	invocation_id := resp.Header["Lambda-Runtime-Aws-Request-Id"][0]
	log.Printf("EVENT open-runtime %s %s", invocation_id, body)
	return invocation_id, body
}

func main() {
	runtime_api := fmt.Sprintf("http://%s/2018-06-01/runtime", os.Getenv("AWS_LAMBDA_RUNTIME_API"))

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
		invocation_id, body := getInvocation(runtime_api)
		stdin.Write(body)
		io.WriteString(stdin, "\n")
		for scanner.Scan() {
			response := scanner.Bytes()
			log.Printf("Response! %s", invocation_id)
			log.Printf("%s/invocation/%s/response", runtime_api, invocation_id)
			_, err := http.Post(
				fmt.Sprintf("%s/invocation/%s/response", runtime_api, invocation_id),
				"application/json",
				bytes.NewBuffer(response),
			)
			if nil != err {
				log.Fatalf("Error POSTing response")
			}
			log.Print("done sending response")
			invocation_id, body = getInvocation(runtime_api)
			stdin.Write(body)
			io.WriteString(stdin, "\n")
		}
	}(reader)

	if err := cmd.Start(); nil != err {
		log.Fatalf("Error starting program: %s, %s", cmd.Path, err.Error())
	}

	cmd.Wait()
}
