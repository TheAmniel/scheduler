package main

import "github.com/theamniel/scheduler/ipc"

func main() {
	process := ipc.New()
	process.Start()
}
