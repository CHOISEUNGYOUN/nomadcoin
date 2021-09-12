package wallet

import (
	"crypto/x509"
	"encoding/hex"
	"io/fs"
	"reflect"
	"strings"
	"testing"
)

const (
	test_key       = "307702010104208ce0e68c605d224c1a7225bfc077a8272900b1803651098f5c731ec7b737f1eea00a06082a8648ce3d030107a14403420004847ed6c4747dae5bfa84a8c50b0dd720249ef738c41bfa07eb0a2a35a77ad44e7bd4240e33eb15a137fac5363048a213d478de8eafac8516b5188cd0d4ada957"
	test_payload   = "000002cef2f257f6ad891370e62c2552904fab179a00bca08420ad1bd105e655"
	test_signature = "a0a294be7c41a724967b1b48d902ef38306616e44c22cd6218c62faae583426226a4cb506ce77239a6509b992f8b6a1c96f1f110c08256e1efaac292aa3ad3da"
)

type fakeLayer struct {
	fakeHasWalletFile func() bool
}

func (f fakeLayer) hasWalletFile() bool {
	return f.fakeHasWalletFile()
}

func (fakeLayer) writeFile(fileName string, data []byte, perm fs.FileMode) error {
	return nil
}

func (fakeLayer) readFile(name string) ([]byte, error) {
	return x509.MarshalECPrivateKey(makeTestWallet().privateKey)
}

func TestWallet(t *testing.T) {
	t.Run("New wallet is created", func(t *testing.T) {
		files = fakeLayer{
			fakeHasWalletFile: func() bool { return false },
		}
		tw := Wallet()
		if reflect.TypeOf(tw) != reflect.TypeOf(&wallet{}) {
			t.Error("New wallet should return a new wallet instance")
		}
	})
	t.Run("Wallet is restored", func(t *testing.T) {
		files = fakeLayer{
			fakeHasWalletFile: func() bool { return true },
		}
		w = nil
		tw := Wallet()
		if reflect.TypeOf(tw) != reflect.TypeOf(&wallet{}) {
			t.Error("Wallet should return a wallet instance")
		}
	})
}

func makeTestWallet() *wallet {
	w := &wallet{}
	b, _ := hex.DecodeString(test_key)
	key, _ := x509.ParseECPrivateKey(b)
	w.privateKey = key
	w.Address = aFromK(key)
	return w
}

func TestSign(t *testing.T) {
	s := Sign(test_payload, *makeTestWallet())
	_, err := hex.DecodeString(s)
	if err != nil {
		t.Errorf("Sign() should return a hex decoded string, got %s", s)
	}
}

func TestVerify(t *testing.T) {
	type test struct {
		input string
		ok    bool
	}
	tests := []test{
		{test_payload, true},
		{strings.Replace(test_payload, "0", "d", 1), false},
	}
	for _, tc := range tests {
		w := makeTestWallet()
		ok := Verify(test_signature, tc.input, w.Address)
		if ok != tc.ok {
			t.Error("Verify() could not verify test_signature and test_payload")
		}
	}
}

func TestRestoreBigInts(t *testing.T) {
	_, _, err := restoreBigInts("xx")
	if err == nil {
		t.Error("restoreBigInts() should return error when payload is not hex.")
	}
}
