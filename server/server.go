package main

import (
	"log"
)

var args struct {
	DbPath   string `arg:"env:MAILINGLIST_DB"`
	BindJson string `arg:"env:MAILINGLIST_BIND_JSON"`
}

func main() {
	arg.mustParse(&args)

	if args.DbPath == "" {
		log.Fatal("MAILINGLIST_DB not set")
	}
}
