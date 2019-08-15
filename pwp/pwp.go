package pwp

import (
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"os/user"
	"runtime"
)

//OpSys is a struct that contains os dependent settings for PWP
//OSName : Operating system name
//LibDir : PWP library dir
//LibUserDir : PWP user lib dir
//PrivUser : Username of the root/admin/whatever
type OpSys struct {
	OSName      string
	LibDir      string
	LibUserDir  string
	PrivUser    string
	CurrentUser string
}

const (
	osx = "darwin"
	win = "windows"
)

func getOS() OpSys {
	var opsys OpSys
	usr, _ := user.Current()
	switch runtime.GOOS {
	case osx:
		opsys.OSName = osx
		opsys.LibDir = "/usr/local/pwp/"
		opsys.LibUserDir = usr.HomeDir + "/.pwp/"
		opsys.PrivUser = "root"
		opsys.CurrentUser = usr.Username
	}
	return opsys
}
func exist(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	}
	return false
}

//IsInitialized : Checks that PWP is initialized or not
//opsys : OpSys struct
func IsInitialized(opsys OpSys) bool {
	if !exist(opsys.LibDir+"key.pem") && !exist(opsys.LibUserDir+"key.pem") {
		return false
	}
	return true
}

//Init - Initializes PWP before the first use
//asUser : indicates that init should be done in usermode (no root)
func Init(asUser bool) error {
	sp := make([]byte, 32)
	rand.Read(sp)
	opsys := getOS()
	fmt.Printf("Start initialization. OS: %s, User: %s\n", opsys.OSName, opsys.CurrentUser)
	if !asUser && opsys.CurrentUser != opsys.PrivUser {
		return errors.New("You need to be root/admin to init PWP system-wide")
	}
	if IsInitialized(opsys) {
		return errors.New("PWP has already been initialized")
	}

	if asUser {
		if _, err := os.Stat(opsys.LibUserDir); err != nil {
			err := os.MkdirAll(opsys.LibUserDir, 0700)
			if err != nil {
				return errors.New("Error creating lib directory: " + opsys.LibUserDir)
			}
		}
		fp, err := os.OpenFile(opsys.LibUserDir+"key.pem", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)
		if err != nil {
			return errors.New("Error creating key file: " + err.Error())
		}
		_, err = fp.Write(sp)
		fp.Close()
		if err != nil {
			os.Remove(opsys.LibUserDir + "key.pem")
			return errors.New("Error writing key to file: " + err.Error())
		}
	}

	return nil
}
