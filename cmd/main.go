package main

import (
	"log"
	"github.com/qnib/go-secbench"
)

func main() {
	sb, err := secbench.NewSecBenc()
	if err != nil {
		log.Panicf(err.Error())
	}
	sb.RunBench()
}
