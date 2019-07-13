package zgocrypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"golang.org/x/crypto/hkdf"
	"hash"
	"io"
)

/*
@Time : 2019-03-15 10:14
@Author : rubinus.chu
@File : crypto
@project: zgo
*/

// AES256KeyLength is the length of key for AES 256 crypt
const AES256KeyLength = 32

var Crypto Cryptoer

type crypto struct {
}

func New() Cryptoer {
	return &crypto{}
}

type Cryptoer interface {
	//md5 对字符串md5
	Md5(s string) string

	//对字符串进行sha1算法
	SHA1(s string) string

	//对字符串进行sha256算法
	SHA256String(s string) string

	//对[]byte进行sha256算法
	SHA256(s []byte) ([]byte, error)

	//aes 256
	AESMake256Key(k []byte) []byte

	//aes-cfb
	AESCFBEncrypt(s, k []byte) ([]byte, error)
	AESCFBDecrypt(c, k []byte) ([]byte, error)

	//aes-cbc
	AESCBCEncrypt(s, k []byte) ([]byte, error)
	AESCBCDecrypt(c, k []byte) ([]byte, error)

	//aes-gcm
	AESGCMEncrypt(s, k, n []byte) ([]byte, error)
	AESGCMDecrypt(c, k, n []byte) ([]byte, error)

	//aes
	AesEncrypt(orig string, key string) string
	AesDecrypt(cryted string, key string) string

	//rsa
	RsaEncrypt(origData []byte, publicKey []byte) ([]byte, error)
	RsaDecrypt(ciphertext []byte, privateKey []byte) ([]byte, error)

	DecryptDataForWXBizData(encryptedData, sessionKey, iv string) ([]byte, error)

	//pkcs5
	PKCS5Padding(c []byte, blockSize int) []byte
	PKCS5UnPadding(s []byte) ([]byte, error)

	//pkcs7
	PKCS7Padding(ciphertext []byte, blocksize int) []byte
	PKCS7UnPadding(origData []byte) []byte

	//HmacSha1
	HmacSha1(message, key []byte) ([]byte, error)
	HmacSha1Check(message, messageMAC, key []byte) (bool, error)

	//HmacSha256
	HmacSha256(message, key []byte) ([]byte, error)
	HmacSha256Check(message, messageMAC, key []byte) (bool, error)

	//hkdf-sha256
	HkdfSha256RandomSalt(secret, info []byte, sl int) (key []byte, salt []byte, err error)
	HkdfSha256WithSalt(secret, salt, info []byte) (key []byte, err error)

	//hkdf-sha1
	HkdfSha1RandomSalt(secret, info []byte, sl int) (key []byte, salt []byte, err error)
	HkdfSha1WithSalt(secret, salt, info []byte) (key []byte, err error)
}

// Md5
func (cp *crypto) Md5(s string) string {
	md5 := md5.New()
	md5.Write([]byte(s))
	return hex.EncodeToString(md5.Sum(nil))
}

// SHA1 encrypt s according to sha1 algorithm
func (cp *crypto) SHA1(s string) string {
	var h hash.Hash
	h = sha1.New()
	io.WriteString(h, s)
	return hex.EncodeToString(h.Sum(nil))
}

// SHA256 encrypt s according to sha256 algorithm
func (cp *crypto) SHA256String(s string) string {
	var h hash.Hash
	h = sha256.New()
	io.WriteString(h, s)
	return hex.EncodeToString(h.Sum(nil))
}

// SHA256 encrypt s according to sha256 algorithm
func (cp *crypto) SHA256(s []byte) ([]byte, error) {
	var h hash.Hash
	h = sha256.New()
	n, err := h.Write(s)
	if err != nil {
		return nil, err
	}
	if n != len(s) {
		return nil, errors.New("Write length error")
	}
	r := h.Sum(nil)
	return r, nil
}

// AESMake256Key cut or append empty data on the key
// and make sure the key lenth equal 32
func (cp *crypto) AESMake256Key(k []byte) []byte {
	if len(k) < AES256KeyLength {
		a := make([]byte, AES256KeyLength-len(k))
		return append(k, a...)
	}
	if len(k) > AES256KeyLength {
		return k[:AES256KeyLength]
	}
	return k
}

// AESCFBEncrypt encrypt s with given k.
// k should be 128/256 bits, otherwise it will append empty data or cut until 256 bits.
// First 16 bytes of cipher data is the IV.
func (cp *crypto) AESCFBEncrypt(s, k []byte) ([]byte, error) {
	if len(k) != 16 && len(k) != 32 {
		k = cp.AESMake256Key(k)
	}
	block, err := aes.NewCipher(k)
	if err != nil {
		return nil, err
	}

	cb := make([]byte, aes.BlockSize+len(s))
	iv := cb[:aes.BlockSize]
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return nil, err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(cb[aes.BlockSize:], s)
	return cb, nil
}

// AESDecrypt decrypt c with given k
// k should be 128/256 bits, otherwise it will append empty data or cut until 256 bits
// First 16 bytes of cipher data is the IV.
func (cp *crypto) AESCFBDecrypt(c, k []byte) ([]byte, error) {
	if len(k) != 16 && len(k) != 32 {
		k = cp.AESMake256Key(k)
	}
	block, err := aes.NewCipher(k)
	if err != nil {
		return nil, err
	}
	if len(c) < aes.BlockSize {
		err := errors.New("crypt data is too short")
		return nil, err
	}

	iv := c[:aes.BlockSize]
	cb := c[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(cb, cb)
	return cb, nil
}

// AESCBCEncrypt encrypt s with given k
// k should be 128/256 bits, otherwise it will append empty data or cut until 256 bits
// First 16 bytes of cipher data is the IV.
func (cp *crypto) AESCBCEncrypt(s, k []byte) ([]byte, error) {
	if len(k) != 16 && len(k) != 32 {
		k = cp.AESMake256Key(k)
	}
	if len(s)%aes.BlockSize != 0 {
		return nil, errors.New("invalid length of s")
	}
	block, err := aes.NewCipher(k)
	if err != nil {
		return nil, err
	}
	cb := make([]byte, aes.BlockSize+len(s))
	iv := cb[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cb[aes.BlockSize:], s)
	return cb, nil
}

// AESCBCDecrypt decrypt c with given k
// k should be 128/256 bits, otherwise it will append empty data or cut until 256 bits
// First 16 bytes of cipher data is the IV.
func (cp *crypto) AESCBCDecrypt(c, k []byte) ([]byte, error) {
	if len(k) != 16 && len(k) != 32 {
		k = cp.AESMake256Key(k)
	}
	if len(c) < aes.BlockSize {
		return nil, errors.New("c too short")
	}
	block, err := aes.NewCipher(k)
	if err != nil {
		return nil, err
	}

	iv := c[:aes.BlockSize]
	cb := c[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cb, cb)
	return cb, nil
}

// PKCS5Padding can append data of PKCS5
// Common blockSize is aes.BlockSize
func (cp *crypto) PKCS5Padding(c []byte, blockSize int) []byte {
	pl := blockSize - len(c)%blockSize
	p := bytes.Repeat([]byte{byte(pl)}, pl)
	return append(c, p...)
}

// PKCS5UnPadding can unappend data of PKCS5
func (cp *crypto) PKCS5UnPadding(s []byte) ([]byte, error) {
	l := len(s)
	if l == 0 {
		return nil, errors.New("s too short")
	}
	pl := int(s[l-1])
	if l < pl {
		return nil, errors.New("s too short")
	}
	return s[:(l - pl)], nil
}

// AESGCMEncrypt encrypt s use k and nonce
func (cp *crypto) AESGCMEncrypt(s, k, n []byte) ([]byte, error) {
	if len(k) != 16 && len(k) != 32 {
		k = cp.AESMake256Key(k)
	}
	block, err := aes.NewCipher(k)
	if err != nil {
		return nil, err
	}
	g, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	c := g.Seal(nil, n, s, nil)
	return c, nil
}

// AESGCMDecrypt decrypt s use k and nonce
func (cp *crypto) AESGCMDecrypt(c, k, n []byte) ([]byte, error) {
	if len(k) != 16 && len(k) != 32 {
		k = cp.AESMake256Key(k)
	}
	block, err := aes.NewCipher(k)
	if err != nil {
		return nil, err
	}
	g, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	s, err := g.Open(nil, n, c, nil)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// AesEncrypt
func (cp *crypto) AesEncrypt(orig string, key string) string {
	// 转成字节数组
	origData := []byte(orig)
	k := []byte(key)
	// 分组秘钥
	block, _ := aes.NewCipher(k)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 补全码
	origData = cp.PKCS7Padding(origData, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	// 创建数组
	cryted := make([]byte, len(origData))
	// 加密
	blockMode.CryptBlocks(cryted, origData)
	return base64.StdEncoding.EncodeToString(cryted)
}

// AesDecrypt
func (cp *crypto) AesDecrypt(cryted string, key string) string {
	// 转成字节数组
	crytedByte, _ := base64.StdEncoding.DecodeString(cryted)
	k := []byte(key)
	// 分组秘钥
	block, _ := aes.NewCipher(k)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	// 创建数组
	orig := make([]byte, len(crytedByte))
	// 解密
	blockMode.CryptBlocks(orig, crytedByte)
	// 去补全码
	orig = cp.PKCS7UnPadding(orig)
	return string(orig)
}

// PKCS7Padding补码
func (cp *crypto) PKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS7UnPadding去码
func (cp *crypto) PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// RsaEncrypt
func (cp *crypto) RsaEncrypt(origData []byte, publicKey []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

// RsaDecrypt
func (cp *crypto) RsaDecrypt(ciphertext []byte, privateKey []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}

func (cp *crypto) AesCBCDecrypt(encryptData, sessionKey, iv []byte) ([]byte, error) {
	var aesBlockDecrypter cipher.Block
	aesBlockDecrypter, err := aes.NewCipher(sessionKey)
	if err != nil {
		return nil, err
	}
	decrypted := make([]byte, len(encryptData))
	aesDecrypter := cipher.NewCBCDecrypter(aesBlockDecrypter, iv)
	aesDecrypter.CryptBlocks(decrypted, encryptData)

	decrypted = cp.PKCS7UnPadding(decrypted)

	return decrypted, nil
}

// DecryptDataForWXBizData
func (cp *crypto) DecryptDataForWXBizData(encryptedData, sessionKey, iv string) ([]byte, error) {

	if len(sessionKey) != 24 {
		return nil, errors.New("error sessionKey len must be 24")
	}

	if len(iv) != 24 {
		return nil, errors.New("error iv len must be 24")
	}

	data, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}
	iKey, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		return nil, err
	}
	iIv, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return nil, err
	}
	dnData, err := cp.AesCBCDecrypt(data, iKey, iIv)
	if err != nil {
		return nil, err
	}
	return dnData, nil
}

// HmacSha256
func (cp *crypto) HmacSha256(message, key []byte) ([]byte, error) {
	mac := hmac.New(sha256.New, key)
	if _, err := mac.Write(message); err != nil {
		return nil, err
	}
	return mac.Sum(nil), nil
}

// HmacSha256Check
func (cp *crypto) HmacSha256Check(message, messageMAC, key []byte) (bool, error) {
	mac := hmac.New(sha256.New, key)
	if _, err := mac.Write(message); err != nil {
		return false, err
	}
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC), nil
}

// HmacSha1
func (cp *crypto) HmacSha1(message, key []byte) ([]byte, error) {
	mac := hmac.New(sha1.New, key)
	if _, err := mac.Write(message); err != nil {
		return nil, err
	}
	return mac.Sum(nil), nil
}

// HmacSha1Check
func (cp *crypto) HmacSha1Check(message, messageMAC, key []byte) (bool, error) {
	mac := hmac.New(sha1.New, key)
	if _, err := mac.Write(message); err != nil {
		return false, err
	}
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC), nil
}

// HkdfSha256RandomSalt
func (cp *crypto) HkdfSha256RandomSalt(secret, info []byte, sl int) (key []byte, salt []byte, err error) {
	hash := sha256.New

	salt = make([]byte, sl)
	if _, err = io.ReadFull(rand.Reader, salt); err != nil {
		return
	}

	hkdf := hkdf.New(hash, secret, salt, info)

	key = make([]byte, hash().Size())
	if _, err = io.ReadFull(hkdf, key); err != nil {
		return
	}
	return
}

// HkdfSha256WithSalt
func (cp *crypto) HkdfSha256WithSalt(secret, salt, info []byte) (key []byte, err error) {
	hash := sha256.New

	hkdf := hkdf.New(hash, secret, salt, info)

	key = make([]byte, hash().Size())
	if _, err = io.ReadFull(hkdf, key); err != nil {
		return
	}
	return
}

// HkdfSha1RandomSalt
func (cp *crypto) HkdfSha1RandomSalt(secret, info []byte, sl int) (key []byte, salt []byte, err error) {
	hash := sha1.New

	salt = make([]byte, sl)
	if _, err = io.ReadFull(rand.Reader, salt); err != nil {
		return
	}

	hkdf := hkdf.New(hash, secret, salt, info)

	key = make([]byte, hash().Size())
	if _, err = io.ReadFull(hkdf, key); err != nil {
		return
	}
	return
}

// HkdfSha1WithSalt
func (cp *crypto) HkdfSha1WithSalt(secret, salt, info []byte) (key []byte, err error) {
	hash := sha1.New

	hkdf := hkdf.New(hash, secret, salt, info)

	key = make([]byte, hash().Size())
	if _, err = io.ReadFull(hkdf, key); err != nil {
		return
	}
	return
}
