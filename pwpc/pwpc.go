package main

import (
	"fmt"
	"os"

	"github.com/csbgl/pwp/pwp"
	"github.com/jessevdk/go-flags"
)

type opts struct {
	Init   bool `short:"i" long:"init" description:"Initialize PWP"`
	AsUser bool `short:"u" long:"user" description:"Initialize as user"`
}

func exist(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	}
	return false
}

func main() {
	var op opts
	args, err := flags.Parse(&op)
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
	if op.Init {
		if op.AsUser {
			fmt.Println("Initializing PWP in User mode")
			err := pwp.Init(true)
			if err != nil {
				fmt.Println("Error initializing PWP: ", err)
			}
		} else {
			fmt.Println("Initializing PWP in system-wide mode")
			err := pwp.Init(false)
			if err != nil {
				fmt.Println("Error initializing PWP: ", err)
			}
		}
	}
	fmt.Println(args)

}
