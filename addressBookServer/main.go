package main

import (
	"httpserver/controller/stdhttp"
	"httpserver/gates/psg"
	"os"
)

func main() {
	psgr := psg.NewPsg("localhost", "postgres", os.Getenv("DB_PASSWORD"))
	defer psgr.Close()
	serv := stdhttp.NewController(":8080", psgr)
	serv.Start()
}
