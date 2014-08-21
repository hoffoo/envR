package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
)

// start an R process and expose an rpc interface
func start(args ...string) *R {

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

	// XXX push env args into R process
	// XXX find R from $PATH
	proc, err := os.StartProcess("/usr/bin/R", []string{"/usr/bin/R", "--no-save"}, attrs)
	if err != nil {
		panic(err)
	}

	wait := make(chan bool, 1)
	r := R{wait, rWrite}

	go listenSignal(proc)
	go listenRPC(&r)

	go func() {
		proc.Wait()
		wait <- false
	}()

	return &r
}

// listen for signals and forward them to R
func listenSignal(proc *os.Process) {

	sig := make(chan os.Signal)
	signal.Notify(sig)

	go func() {
		for {
			// XXX be smarter about which signals kill
			<-sig
			fmt.Println("\n\n\nShutting down R")
			proc.Kill()
		}
	}()
}

// listen rpc
func listenRPC(r *R) {

	rpc.Register(r)
	rpc.HandleHTTP()

	env := GetEnv()
	l, e := net.Listen("tcp", env.String())
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}

// check if rpc running and return client if so
func checkRunning() (running bool, client *rpc.Client) {

	env := GetEnv()

	client, err := rpc.DialHTTP("tcp", env.String())
	if err != nil {
		return false, nil
	}

	return true, client
}
