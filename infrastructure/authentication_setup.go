package infrastructure

import (
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"

	"github.com/go-chi/jwtauth"
)

const Algorithm = "RS256"

func loadAuthToken() error {
	// Load private key
	privateReader, err := ioutil.ReadFile(rsaPrivatePath)
	if err != nil {
		log.Println("No RSA private pem file.")
		return err
	}

	privatePem, _ := pem.Decode(privateReader)
	privateKey, err = x509.ParsePKCS1PrivateKey(privatePem.Bytes)
	if err != nil {
		log.Println("problem creating authentication:", err)
		return err
	}
	// Read public key
	publicReader, err := ioutil.ReadFile(rsaPublicPath)
	if err != nil {
		log.Println("No RSA public pem file.")
		return err
	}
	publicPem, _ := pem.Decode(publicReader)
	publicKey, _ = x509.ParsePKIXPublicKey(publicPem.Bytes)

	// signKey = privateKey
	// verifyKey = publicKey
	encodeAuth = jwtauth.New(Algorithm, privateKey, publicKey)
	//decodeAuth = jwtauth.New(Algorithm, nil, publicKey)

	return nil
}
