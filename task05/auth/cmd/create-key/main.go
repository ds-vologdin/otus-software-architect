package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
)

func main() {
	reader := rand.Reader
	bitSize := 2048

	key, err := rsa.GenerateKey(reader, bitSize)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		return
	}

	privateKey := x509.MarshalPKCS1PrivateKey(key)
	privateKeyStr := base64.StdEncoding.EncodeToString(privateKey)
	fmt.Printf("privateKey: %s\n", privateKeyStr)

	publicKey := x509.MarshalPKCS1PublicKey(&key.PublicKey)
	publicKeyStr := base64.StdEncoding.EncodeToString(publicKey)
	fmt.Printf("publicKey: %s\n", publicKeyStr)
}
