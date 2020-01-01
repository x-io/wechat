package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

//AesCBCEncrypt AES-128-CBC Encrypt
func AesCBCEncrypt(data []byte, key []byte) ([]byte, error) {
	// 分组秘钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()                              // 获取秘钥块的长度
	data = pkcs5Padding(data, blockSize)                        // 补全码
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize]) // 加密模式
	encrypted := make([]byte, len(data))                        // 创建数组
	blockMode.CryptBlocks(encrypted, data)                      // 加密
	return encrypted, nil
}

//AesCBCDecrypt AES-128-CBC Decrypt
func AesCBCDecrypt(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key) // 分组秘钥
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()                              // 获取秘钥块的长度
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize]) // 加密模式
	decrypted := make([]byte, len(data))                        // 创建数组
	blockMode.CryptBlocks(decrypted, data)                      // 解密
	decrypted = pkcs5UnPadding(decrypted)                       // 去除补全码
	return decrypted, nil
}

//AesCBCDecrypt2 AES-128-CBC Decrypt
func AesCBCDecrypt2(data, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key) // 分组秘钥
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv) // 加密模式
	decrypted := make([]byte, len(data))           // 创建数组
	blockMode.CryptBlocks(decrypted, data)         // 解密
	decrypted = pkcs5UnPadding(decrypted)          // 去除补全码
	return decrypted, nil
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//AesECBEncrypt AES-ECB Encrypt
//bit: 128 192  256位的其中一个 长度 对应分别是 16 24  32字节长度
func AesECBEncrypt(bit int, data, key []byte) ([]byte, error) {
	cipher, err := aes.NewCipher(generateKey(key, bit))
	if err != nil {
		return nil, err
	}
	length := (len(data) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, data)
	pad := byte(len(plain) - len(data))
	for i := len(data); i < len(plain); i++ {
		plain[i] = pad
	}
	encrypted := make([]byte, len(plain))
	// 分组分块加密
	for bs, be := 0, cipher.BlockSize(); bs <= len(data); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Encrypt(encrypted[bs:be], plain[bs:be])
	}

	return encrypted, nil
}

//AesECBDecrypt AES-ECB Decrypt
func AesECBDecrypt(bit int, data, key []byte) ([]byte, error) {
	cipher, err := aes.NewCipher(generateKey(key, bit))
	if err != nil {
		return nil, err
	}
	decrypted := make([]byte, len(data))
	//
	for bs, be := 0, cipher.BlockSize(); bs < len(data); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Decrypt(decrypted[bs:be], data[bs:be])
	}

	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}

	return decrypted[:trim], nil
}
func generateKey(key []byte, bit int) (genKey []byte) {
	genKey = make([]byte, bit)
	copy(genKey, key)
	for i := bit; i < len(key); {
		for j := 0; j < bit && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}

//AesCFBEncrypt AES-CFB Encrypt
func AesCFBEncrypt(origData []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	encrypted := make([]byte, aes.BlockSize+len(origData))
	iv := encrypted[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(encrypted[aes.BlockSize:], origData)
	return encrypted, nil
}

//AesCFBDecrypt AES-CFB Decrypt
func AesCFBDecrypt(encrypted []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(encrypted) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	iv := encrypted[:aes.BlockSize]
	encrypted = encrypted[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(encrypted, encrypted)
	return encrypted, nil
}
