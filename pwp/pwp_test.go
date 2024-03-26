package pwp

import (
	"encoding/hex"
	"fmt"
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
