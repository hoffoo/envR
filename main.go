// A mini R environment daemon
package main

import (
	"bufio"
	"fmt"
	"net/rpc"
	"os"
)

// If R is not running start it up. Otherwise send stdin to running R.
func main() {

	rpcR := checkRunning()
	if rpcR == nil {
		fmt.Println("Starting R")
		r := startR(os.Args[1:]...)
		<-r.wait
	} else {
		sendToR(rpcR)
	}
}

// check if rpc running and return client if so
func checkRunning() (rpcR *rpc.Client) {

	env := GetEnv()

	rpcR, err := rpc.DialHTTP("tcp", env.String())
	if err != nil {
		return nil
	}

	return rpcR
}

// start an R process and expose an rpc interface
func startR(args ...string) *R {

	// connect R's stdin
	rStdin, rWrite, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	attrs := &os.ProcAttr{
		Dir: "./",
		Files: []*os.File{
			rStdin,
			os.Stdout,
			os.Stderr,
		},
	}

	// make args
	rArgs := make([]string, len(args)+1)
	// XXX find R from $PATH
	rArgs[0] = "/usr/bin/R"
	for i, arg := range args {
		rArgs[i+1] = arg
	}

	// XXX find R from $PATH
	proc, err := os.StartProcess("/usr/bin/R", rArgs, attrs)
	if err != nil {
		panic(err)
	}

	wait := make(chan bool, 1)
	r := R{wait, rWrite}

	go listenSignal(proc)
	go listenRPC(&r)

	go func() {
		proc.Wait()
		rStdin.Close()
		rWrite.Close()
		wait <- false
	}()

	return &r
}

// reads all stdin and sends to R
func sendToR(client *rpc.Client) {

	reader := bufio.NewReader(os.Stdin)
	var resp string

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		if line[len(line)-1] == '\n' {
			err = client.Call("R.Pipe", line, &resp)
		} else {
			err = client.Call("R.Pipe", line, &resp)
		}

		if err != nil {
			fmt.Println("ERROR: ", err)
		}
	}
}
