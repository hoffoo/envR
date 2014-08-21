// A mini R environment daemon
package main

import (
	"bufio"
	"fmt"
	"io"
	"net/rpc"
	"os"
)

// If R is not running start it up. Otherwise send stdin to running R.
func main() {

	if rpcR := checkRunning(); rpcR != nil {
		// R is running, send stdin to it
		sendToR(rpcR)
	} else {
		// R is not running, start
		fmt.Println("Starting R")
		if r := startR(os.Args[1:]...); r != nil {
			<-r.wait
		}
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

	rBin := findR()
	if rBin == "" {
		fmt.Println("Couldnt find R, check your $PATH")
		return nil
	}

	// make args
	rArgs := append([]string{rBin}, args...)

	proc, err := os.StartProcess(rBin, rArgs, attrs)
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
		close(r.wait) // done, close and exit
	}()

	// copy stdin to R
	// XXX make this not halt R when it fails
	go func() {
		reader := bufio.NewReader(os.Stdin)

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}

			io.WriteString(r.stdin, line)
			if err != nil {
				fmt.Println("ERROR: ", err)
			}
		}
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
