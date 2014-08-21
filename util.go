package main

import "os"
import "path/filepath"
import "strings"

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

// search $PATH for R
func findR() string {

	pathEnv := os.Getenv("PATH")

	if pathEnv == "" {
		return ""
	}

	pathList := strings.SplitAfter(pathEnv, string(os.PathListSeparator))

	for _, path := range pathList {

		if path[len(path)-1] == os.PathListSeparator {
			path = path[:len(path)-1]
		}

		pathR := filepath.Join(path, "R")
		if _, err := os.Stat(pathR); !os.IsNotExist(err) {
			return pathR
		}
	}

	return ""
}
