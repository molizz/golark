package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
)

func AesDecrypt(key string, encrypt string) (string, error) {
	kbs := SHA256(key)
	decode, err := base64.StdEncoding.DecodeString(encrypt)
	if err != nil {
		return "", err
	}
	if len(decode) < aes.BlockSize {
		return "", errors.New("encrypt short")
	}
	iv := decode[:aes.BlockSize]
	block, err := aes.NewCipher(kbs)
	if err != nil {
		return "", err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	plantText := make([]byte, len(decode))
	blockMode.CryptBlocks(plantText, decode)
	plantText = PKCS7UnPadding(plantText)
	plantText = plantText[aes.BlockSize:]
	return string(plantText), nil
}

func PKCS7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	unpadding := int(plantText[length-1])
	return plantText[:(length - unpadding)]
}

func SHA256(source string) []byte {
	mac := sha256.New()
	mac.Write([]byte(source))
	return mac.Sum(nil)
}

func SHA1(source string) string {
	sha1 := sha1.New()
	sha1.Write([]byte(source))
	return hex.EncodeToString(sha1.Sum([]byte(nil)))
}
