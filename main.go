// A mini R environment daemon
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {

	var r *R

	running, client := checkRunning()
	if !running {
		fmt.Println("Starting R")
		r = start()
		<-r.wait
	} else {
		reader := bufio.NewReader(os.Stdin)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			var res string
			err = client.Call("R.Pipe", line, &res)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

type R struct {
	wait  chan bool
	stdin *os.File
}

func (r *R) Pipe(s string, result *string) error {
	io.WriteString(r.stdin, s)
	return nil
}

type Env struct {
	addr string
	port string
}

func GetEnv() Env {

	addr := os.Getenv("R_HOST")
	port := os.Getenv("R_PORT")

	if addr == "" {
		addr = "127.0.0.1"
	}

	if port == "" {
		port = "64001"
	}

	return Env{
		addr: addr,
		port: port,
	}
}

func (env *Env) String() string {
	return env.addr + ":" + env.port
}
