package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

// pkcs7Padding 填充模式
func pkcs7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	//Repeat()函数的功能是把切片[]byte{byte(padding)}复制padding个，然后合并成新的字节切片返回
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

// pkcs7UnPadding 填充的反向操作，删除填充字符串
func pkcs7UnPadding(origData []byte) ([]byte, error) {
	//获取数据长度
	length := len(origData)
	if length == 0 {
		return nil, errors.New("加密字符串错误！")
	} else {
		//获取填充字符串长度
		unPadding := int(origData[length-1])
		//截取切片，删除填充字节，并且返回明文
		return origData[:(length - unPadding)], nil
	}
}

// AesEncrypt 实现加密
func AesEncrypt(origData []byte, key []byte) ([]byte, error) {
	//创建加密算法实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create cipher: %w", err)
	}
	//获取块的大小
	blockSize := block.BlockSize()
	//对数据进行填充，让数据长度满足需求
	origData = pkcs7Padding(origData, blockSize)

	// 生成随机IV
	iv := make([]byte, blockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, fmt.Errorf("generate IV: %w", err)
	}

	//采用AES加密方法中CBC加密模式
	blocMode := cipher.NewCBCEncrypter(block, iv)
	encryptedData := make([]byte, len(origData))
	//执行加密
	blocMode.CryptBlocks(encryptedData, origData)

	// 返回 IV + 密文
	return append(iv, encryptedData...), nil
}

// AesDecrypt 实现解密
func AesDecrypt(encryptedData []byte, key []byte) (string, error) {
	//创建加密算法实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("create cipher: %w", err)
	}
	//获取块大小
	blockSize := block.BlockSize()
	if len(encryptedData) < blockSize {
		return "", errors.New("ciphertext too short")
	}

	// 提取IV和密文
	iv := encryptedData[:blockSize]
	ciphertext := encryptedData[blockSize:]

	if len(ciphertext)%blockSize != 0 {
		return "", errors.New("ciphertext is not a multiple of the block size")
	}

	//创建加密客户端实例
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(ciphertext))
	//这个函数也可以用来解密
	blockMode.CryptBlocks(origData, ciphertext)
	//去除填充字符串
	origData, err = pkcs7UnPadding(origData)
	if err != nil {
		return "", err
	}
	return string(origData), nil
}
