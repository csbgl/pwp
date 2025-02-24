package pwp

import (
	"encoding/hex"
	"fmt"
	"os"

	//"os"
	"reflect"
	"testing"
)

func TestCrypt(t *testing.T) {
	Password := "TestPassword"
	keyenc, err := getkey(true)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("[TEST] Encryption Key: %s\n", hex.EncodeToString(keyenc))
	enc, err := encrypt([]byte(Password), keyenc)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("[TEST] Encrypted: %s\n", enc)
	keydec, err := getkey(true)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("[TEST] Decryption Key: %s\n", hex.EncodeToString(keydec))
	benc, _ := hex.DecodeString(enc)
	dec, err := decrypt(benc, keydec)
	fmt.Printf("[TEST] Decrypted: %s\n", dec)
	if err != nil {
		t.Error(err)
	}
	bdec, _ := hex.DecodeString(dec)
	if !reflect.DeepEqual(bdec, []byte(Password)) {
		t.Errorf("Test failed %s != %s", Password, dec)
	}
}

func TestMain(t *testing.T) {
	fmt.Printf("[TEST] Starting tests\n")
	if !IsInitialized(getOS()) {
		err := Init(true)
		if err != nil {
			t.Error(err)
		}
		fmt.Printf("[TEST] Init completed\n")
	}
	err := AddPW(true, "./test.pwp", "TestPassword", "/usr/local/go/bin/go test -timeout 30s -run ^TestMain$ github.com/csbgl/pwp/pwp", "Password123")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("[TEST] Password added\n")
	err = ListPW(true, "./test.pwp")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("[TEST] Password listed\n")
	pw, err := GetPW(true, "./test.pwp", "TestPassword", "/usr/local/go/bin/go test -timeout 30s -run ^TestMain$ github.com/csbgl/pwp/pwp")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("[TEST] Password retrieved: %s\n", pw)
	err = DeletePW(true, "./test.pwp", "TestPassword")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("[TEST] Password deleted\n")
	err = ListPW(true, "./test.pwp")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("[TEST] Password listed\n")
	os.Remove("./test.pwp")
	fmt.Printf("[TEST] Tests completed\n")
}
