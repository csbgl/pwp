package pwp

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"runtime"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
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
	osx   = "darwin"
	win   = "windows"
	linux = "linux"
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
	case linux:
		opsys.OSName = linux
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

func getkey(AsUser bool) ([]byte, error) {
	StaticPart := make([]byte, 32)
	Key := make([]byte, 32)
	opsys := getOS()
	mID, err := getMachineID(opsys.OSName)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	if AsUser == true {
		fp, err := os.Open(opsys.LibUserDir + "key.pem")
		if err != nil {
			return nil, errors.New(err.Error())
		}
		fp.Read(StaticPart)
	} else {
		fp, err := os.Open(opsys.LibDir + "key.pem")
		if err != nil {
			return nil, errors.New(err.Error())
		}
		fp.Read(StaticPart)

	}
	for i := 0; i < 32; i++ {
		Key[i] = StaticPart[i] ^ mID[i]
	}
	return Key, nil
}

//IsInitialized : Checks that PWP is initialized or not
//opsys : OpSys struct
func IsInitialized(opsys OpSys) bool {
	if !exist(opsys.LibDir+"key.pem") && !exist(opsys.LibUserDir+"key.pem") {
		return false
	}
	return true
}

func getMachineID(_os string) ([]byte, error) {
	var result string
	var byteresult [32]byte
	switch _os {
	case osx:
		re := regexp.MustCompile("UUID.*")
		out, err := exec.Command("ioreg", "-rd1", "-c", "IOPlatformExpertDevice").Output()
		if err != nil {
			return nil, err
		}
		result = re.FindString(string(out))
		byteresult = sha256.Sum256([]byte(result))

	case linux:
		b, err := ioutil.ReadFile("/var/lib/dbus/machine-id")
		if err != nil && os.IsNotExist(err) {
			b, err = ioutil.ReadFile("/etc/machine-id")
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
		byteresult = sha256.Sum256(b)
	}
	return byteresult[:], nil
}

func encrypt(data []byte, key []byte) (string, error) {
	nonce := make([]byte, 12)
	io.ReadFull(rand.Reader, nonce)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	encdata := gcm.Seal(nil, nonce, data, nil)
	return hex.EncodeToString(append(nonce, encdata...)), nil

}

func decrypt(bucket []byte, key []byte) (string, error) {
	if len(bucket) < 28 {
		return "", errors.New("decrypt - Length of data is insufficient")
	}
	nonce := bucket[0:12]
	data := bucket[12:]
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", errors.New("NewCipher: " + err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.New("NewGCM: " + err.Error())
	}
	decdata, err := gcm.Open(nil, nonce, data, nil)
	if err != nil {
		return "", errors.New("Open: " + err.Error())
	}
	return hex.EncodeToString(decdata), nil
}

func objectExist(ObjectName string, FileName string) (bool, error) {
	_, err := os.Stat(FileName)
	if os.IsNotExist(err) {
		return false, nil
	}
	b, err := ioutil.ReadFile(FileName)
	if err != nil {
		return false, errors.New("objectExist - " + err.Error())
	}
	s := string(b)
	if strings.Contains(s, ObjectName) {
		return true, nil
	}
	return false, nil
}

func getObject(ObjectName string, FileName string) (string, error) {
	fp, err := os.OpenFile(FileName, os.O_RDONLY, 0)
	if err != nil {
		return "", err
	}
	defer fp.Close()
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Split(line, " ")[0] == ObjectName {
			return line, nil
		}
	}
	return "", errors.New("Object not found")
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

// AddPW - Get password from STDIN, encrypt and write to file
func AddPW(AsUser bool, FileName string, ObjectName string) error {
	opsys := getOS()
	if FileName == "" {
		if AsUser {
			FileName = opsys.LibUserDir + "password"
		} else {
			FileName = opsys.LibDir + "password"
		}
	}
	exist, err := objectExist(ObjectName, FileName)
	if err != nil {
		return err
	}
	if exist {
		return errors.New("Object " + ObjectName + " is already exist")
	}
	usr, _ := user.Current()
	Key, err := getkey(AsUser)
	if err != nil {
		return err
	}

	fmt.Print("Enter 1st part: ")
	bytePassword1, _ := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println("")
	fmt.Print("Enter 2nd part: ")
	bytePassword2, _ := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println("")
	bytePassword := append(bytePassword1, bytePassword2...)
	strPass, err := encrypt(bytePassword, Key)
	if err != nil {
		return err
	}
	stw := ObjectName + " " + usr.Username + " " + strPass
	h := sha256.Sum256([]byte(stw))
	signature, err := encrypt(h[:], Key)
	if err != nil {
		return err
	}
	fp, err := os.OpenFile(FileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	defer fp.Close()
	fp.WriteString(stw + " " + signature + "\n")
	fp.Close()
	return nil
}

// GetPW - Decrypt password and returns.
func GetPW(AsUser bool, FileName string, ObjectName string) (string, error) {
	opsys := getOS()
	if FileName == "" {
		if AsUser {
			FileName = opsys.LibUserDir + "password"
		} else {
			FileName = opsys.LibDir + "password"
		}
	}
	exist, err := objectExist(ObjectName, FileName)
	if err != nil {
		return "", err
	}
	if !exist {
		return "", errors.New("Object " + ObjectName + " does not exist")
	}
	usr, _ := user.Current()
	Key, err := getkey(AsUser)
	if err != nil {
		return "", err
	}
	line, err := getObject(ObjectName, FileName)
	if err != nil {
		return "", err
	}
	parts := strings.Split(line, " ")
	bSignature, _ := hex.DecodeString(parts[3])
	strHash, err := decrypt(bSignature, Key)
	if err != nil {
		return "", errors.New("Signature verification: " + err.Error())
	}
	stc := strings.Join(parts[0:3], " ")
	cHash := sha256.Sum256([]byte(stc))
	if strHash != hex.EncodeToString(cHash[:]) {
		return "", errors.New("Signature verification failed - data currupted")
	}
	if parts[1] != usr.Username {
		return "", errors.New("GetPW - User: " + usr.Username + " is not authorized to read " + ObjectName)
	}
	bytePassword, _ := hex.DecodeString(parts[2])
	strPass, err := decrypt(bytePassword, Key)
	if err != nil {
		return "", errors.New("Password decrypt failed: " + err.Error())
	}
	buff, err := hex.DecodeString(strPass)
	if err != nil {
		return "", errors.New("Password decode failed: " + err.Error())
	}
	return string(buff), nil

}
