package main

import (
	"log"

	"github.com/jrmanes/go-toggl/cmd/api/bootstrap"
)

func main()  {
	if err := bootstrap.Run(); err != nil {
		log.Fatal(err)
	}
}