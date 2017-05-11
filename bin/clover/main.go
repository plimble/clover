package main

import (
	"github.com/plimble/clover/server"
)

func main() {
	config := server.GetConfig()

	s, err := server.New(config, nil)
	if err != nil {
		panic(err)
	}

	s.Run()
}
