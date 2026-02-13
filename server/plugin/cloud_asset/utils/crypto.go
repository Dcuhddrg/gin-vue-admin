package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
)

var encryptionKey []byte

// InitCrypto 初始化加密密钥
func InitCrypto(key string) {
	var k []byte
	if key == "" {
		// 使用默认密钥（生产环境应从配置文件读取）
		// 之前 "default-32-byte-encryption-key!" 长度为 31，导致 aes.NewCipher 报错
		// 这里修正为 32 字节长度
		k = []byte("default-32-byte-encryption-key!!")
	} else {
		k = []byte(key)
	}

	// 确保密钥长度符合 AES 要求 (16, 24, 32)
	// 如果不符合，使用 SHA-256 哈希生成 32 字节密钥
	if len(k) != 16 && len(k) != 24 && len(k) != 32 {
		hash := sha256.Sum256(k)
		k = hash[:]
	}
	encryptionKey = k
}

// Encrypt 加密字符串
func Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt 解密字符串
func Decrypt(ciphertext string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, cipherBytes := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, cipherBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
