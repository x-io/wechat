package util

import (
	"crypto/tls"
	"encoding/pem"

	"golang.org/x/crypto/pkcs12"
)

//P12ToPem 将Pkcs12转成Pem
func P12ToPem(p12 []byte, password string) (cert tls.Certificate, err error) {

	var blocks []*pem.Block
	blocks, err = pkcs12.ToPEM(p12, password)

	if err != nil {
		return
	}

	var pemData []byte
	for _, b := range blocks {
		pemData = append(pemData, pem.EncodeToMemory(b)...)
	}

	cert, err = tls.X509KeyPair(pemData, pemData)

	return
}