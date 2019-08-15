package main

import (
	"fmt"
	"os"

	"github.com/csbgl/pwp/pwp"
	"github.com/jessevdk/go-flags"
)

type general struct {
	Version bool `short:"v" long:"version" description:"Prints version information"`
}
type i struct {
	AsUser bool `short:"r" long:"asuser" description:"Initialize as user"`
}
type add struct {
	AsUser bool   `short:"r" long:"asuser" description:"Initialize as user"`
	User   string `short:"u" long:"user" value-name:"USER" description:"Password is restricted to <user>" required:"yes"`
}

func exist(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	}
	return false
}

func (c *i) Execute(args []string) error {
	if c.AsUser {
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
	return nil
}

func (c *add) Execute(args []string) error {
	return nil
}

func main() {
	var General general
	parser := flags.NewNamedParser("pwpc", flags.Default)
	parser.AddCommand("init", "Initiallize PWP storage", "Initialize PWP storage", new(i))
	parser.AddCommand("add", "Add password", "Add password to PWP storage", new(add))
	parser.AddGroup("General", "General options", &General)
	args, err := parser.Parse()
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
	fmt.Println(args)

}
