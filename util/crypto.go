package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/tls"
	"encoding/pem"

	"golang.org/x/crypto/pkcs12"
)

// func GetTransport2(certData []byte, password string) (*http.Transport, error) {
// 	cert, err := P12ToPem(certData, password)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &http.Transport{
// 		TLSClientConfig: &tls.Config{
// 			Certificates: []tls.Certificate{cert}, // 将pkcs12证书转成pem
// 		},
// 		DisableCompression: true,
// 	}, nil
// }

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

//AESDecrypt AESDecrypt
func AESDecrypt(key, iv, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockMode := cipher.NewCBCDecrypter(block, iv)
	blockMode.CryptBlocks(data, data)
	data = unpadding(data)
	return data, nil
}

func unpadding(src []byte) []byte {
	n := len(src)
	unPadNum := int(src[n-1])
	return src[:n-unPadNum]
}
