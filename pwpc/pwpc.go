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
	AsUser   bool   `short:"r" long:"asuser" description:"Initialize as user"`
	Name     string `short:"n" long:"name" value-name:"NAME" description:"NAME identifying the object" required:"yes"`
	FileName string `short:"f" long:"file" value-name:"FN" description:"File name where the passwords are stored"`
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
	err := pwp.AddPW(c.AsUser, c.FileName, c.Name)
	if err != nil {
		fmt.Println("Error adding password: ", err)
	}
	return nil
}

/*func (c *get) Execute(args []string) error {
	str, err := pwp.GetPW(c.AsUser, c.FileName, c.Name)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println("Password: " + str)
	return nil
}
*/

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
