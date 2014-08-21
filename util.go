package main

import "os"

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
