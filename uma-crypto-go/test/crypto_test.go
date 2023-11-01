package test

import (
	"github.com/uma-universal-money-address/uma-crypto-uniffi/uma-crypto-go"
	"testing"
)

func TestEncryptAndDecrypt(t *testing.T) {
	pubKey, privKey, err := umacrypto.GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}
	secretMessage := "Hello World!"

	encryptedMessage, err := umacrypto.EncryptEcies([]byte(secretMessage), pubKey)
	if err != nil {
		t.Fatal(err)
	}

	decryptedMessage, err := umacrypto.DecryptEcies(encryptedMessage, privKey)
	if err != nil {
		t.Fatal(err)
	}

	if string(decryptedMessage) != secretMessage {
		t.Fatal("Decrypted message does not match original message")
	}
}

func TestSignAndVerify(t *testing.T) {
	pubKey, privKey, err := umacrypto.GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}
	message := "Hello World!"

	signature, err := umacrypto.SignEcdsa([]byte(message), privKey)
	if err != nil {
		t.Fatal(err)
	}

	valid, err := umacrypto.VerifyEcdsa([]byte(message), signature, pubKey)
	if err != nil {
		t.Fatal(err)
	}

	if !valid {
		t.Fatal("Signature is not valid")
	}
}
