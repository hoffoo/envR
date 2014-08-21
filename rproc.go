package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"
)

type R struct {
	wait  chan bool
	stdin *os.File
}

// Send string expression to R via rpc
func (r *R) Pipe(b []byte, result *string) error {

	r.stdin.Write(b)
	if b[len(b)-1] != '\n' {
		r.stdin.Write([]byte{byte('\n')})
	}

	return nil
}

// listen for signal and forward to process
func listenSignal(proc *os.Process) {

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		for {
			<-sig
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
