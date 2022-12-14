package main

import (
	"github.com/theamniel/scheduler/database"
	"github.com/theamniel/scheduler/ipc"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		panic(err)
	}
	process := ipc.New()
	process.SetDatabase(db)
	process.Start()
}
