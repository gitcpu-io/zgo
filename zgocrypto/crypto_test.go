package zgocrypto

import (
	"encoding/hex"
	"fmt"
	"git.zhugefang.com/gocore/zgo/zgoutils"
	"testing"
)

/*
@Time : 2019-03-15 10:16
@Author : rubinus.chu
@File : crypto_test
@project: zgo
*/

var u = New()

//var str = "http://www.baidu.com"
var str = "This is a test!"
var key = "ZhugeAd"

var token_key = "zhugefang2016"

func TestCrypto_AesDecrypt(t *testing.T) {
	key := "9871267812345mn812345xyz"
	encrypt := u.AesEncrypt(str, key)
	t.Log(encrypt)

	decryptCode := u.AesDecrypt(encrypt, key)
	fmt.Println("解密结果：", decryptCode)

}

func TestCrypto_RsaEncrypt(t *testing.T) {
	//	var publicKey = []byte(`
	//-----BEGIN PUBLIC KEY-----
	//MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDfw1/P15GQzGGYvNwVmXIGGxea
	//8Pb2wJcF7ZW7tmFdLSjOItn9kvUsbQgS5yxx+f2sAv1ocxbPTsFdRc6yUTJdeQol
	//DOkEzNP0B8XKm+Lxy4giwwR5LJQTANkqe4w/d9u129bRhTu/SUzSUIr65zZ/s6TU
	//GQD6QzKY1Y8xS+FoQQIDAQAB
	//-----END PUBLIC KEY-----
	//`)
	//	var privateKey = []byte(`
	//-----BEGIN RSA PRIVATE KEY-----
	//MIICXQIBAAKBgQDfw1/P15GQzGGYvNwVmXIGGxea8Pb2wJcF7ZW7tmFdLSjOItn9
	//kvUsbQgS5yxx+f2sAv1ocxbPTsFdRc6yUTJdeQolDOkEzNP0B8XKm+Lxy4giwwR5
	//LJQTANkqe4w/d9u129bRhTu/SUzSUIr65zZ/s6TUGQD6QzKY1Y8xS+FoQQIDAQAB
	//AoGAbSNg7wHomORm0dWDzvEpwTqjl8nh2tZyksyf1I+PC6BEH8613k04UfPYFUg1
	//0F2rUaOfr7s6q+BwxaqPtz+NPUotMjeVrEmmYM4rrYkrnd0lRiAxmkQUBlLrCBiF
	//u+bluDkHXF7+TUfJm4AZAvbtR2wO5DUAOZ244FfJueYyZHECQQD+V5/WrgKkBlYy
	//XhioQBXff7TLCrmMlUziJcQ295kIn8n1GaKzunJkhreoMbiRe0hpIIgPYb9E57tT
	///mP/MoYtAkEA4Ti6XiOXgxzV5gcB+fhJyb8PJCVkgP2wg0OQp2DKPp+5xsmRuUXv
	//720oExv92jv6X65x631VGjDmfJNb99wq5QJBAMSHUKrBqqizfMdOjh7z5fLc6wY5
	//M0a91rqoFAWlLErNrXAGbwIRf3LN5fvA76z6ZelViczY6sKDjOxKFVqL38ECQG0S
	//pxdOT2M9BM45GJjxyPJ+qBuOTGU391Mq1pRpCKlZe4QtPHioyTGAAMd4Z/FX2MKb
	//3in48c0UX5t3VjPsmY0CQQCc1jmEoB83JmTHYByvDpc8kzsD8+GmiPVrausrjj4p
	//y2DQpGmUic2zqCxl6qXMpBGtFEhrUbKhOiVOJbRNGvWW
	//-----END RSA PRIVATE KEY-----
	//`)

	var publicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIGeMA0GCSqGSIb3DQEBAQUAA4GMADCBiAKBgGRjsDmDjmo+aMJHt4OCoXNz3Cc3
rwP52yCGwgLKx62bZl9bZ9N9qYxz1P1u4ol1JbMayhekub3g5VEAh0QOQD4ptXCl
bf9KKvOl4u+peSyaI2uJBZwxH4kM1PmCAb/bekoWdM34/g6t0u4b0uW+6faR08/2
G5T6lQzfonYg9W09AgMBAAE=
-----END PUBLIC KEY-----
`)

	var privateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgGRjsDmDjmo+aMJHt4OCoXNz3Cc3rwP52yCGwgLKx62bZl9bZ9N9
qYxz1P1u4ol1JbMayhekub3g5VEAh0QOQD4ptXClbf9KKvOl4u+peSyaI2uJBZwx
H4kM1PmCAb/bekoWdM34/g6t0u4b0uW+6faR08/2G5T6lQzfonYg9W09AgMBAAEC
gYAGstQuckXNdw6ccBahstS9FG0LPN7zv/dMK1HEx+0c+3KQSFQcO5VVpuD+bCi3
jVaPpej4POib7+a9SfC8rMrChZywYoIJC0SOdtnliyKdlg+YbJNAOqPRfx5lO7N0
GwlwQBAruM8Mj5+lK3WyfntLnbAoL+ehP4AbDk8h2kxzgQJBALDV6BzVhG6dAlfv
geh3fjfsjnyC+MaowLPGtdj99kV5ql8H5r/oa5njfkP4ss+qyJPbBV0qwSxyuN4h
brqmjB0CQQCRVLuR4QnbDQ0HckOaPjqbv+Ex/t+4NCeoCK6q6AzhIyV0FBvpP3jj
kX2YG+J9GqTyzcViNy9MF/4LF/8tvVuhAkEAsGSfhWujD6eEt+zMe/P+LpDETQVk
XN3mYy3bJVmnOAQ141eBYHraIoZf5DOHQklMjNhUEGxehxTkvYdB6ePKLQJANFfE
RsI6fn6w1qG832zSRw5Q2WTmxr1iGw8lecCiAuRZtPkZJGUL8UqRZZuFc3GJPzoF
gV5G1q09flliOuMcgQJBAKij+dhYhtBcqj2430OBHDN/QTk0rAnVVMqqo2VAZ+SO
sQU9l8bA5uU6ASJ25JLqxSLj8IpcVcccqsB+0/eHZ9Y=
-----END RSA PRIVATE KEY-----
`)

	//publicKey, privateKey = u.GenerateRSAKey(1024)
	fmt.Println(string(publicKey))
	fmt.Println("md5公匙", u.Md5(string(publicKey)))
	fmt.Println(string(privateKey))
	data, _ := u.RsaEncrypt([]byte(str), publicKey)
	fmt.Println("加密后：", zgoutils.Utils.ByteToBase64(data))
	b64 := "C4BcuNlEFe3mh/jihPJAekIsS8Dl0bUM7ds65cUT48ZMfIKJg5mgljG7wnFirejBXHCZfiGE4XhU9NaLpj5AFfY/ddRj84IA6PnHwWaycqKfG07vO+dYyUqpNXi2OJG/u0e4iB2D/vpZqDVD6OoHVHm6b1O9obD/AKWSHRPTQfY="
	bytes, _ := zgoutils.Utils.Base64ToByte(b64)
	origData, _ := u.RsaDecrypt(bytes, privateKey)
	fmt.Println("解密后：", string(origData))

	fmt.Println("sha256：", u.SHA256String(str))
	fmt.Println("sha1：", u.SHA1(str))
}

func TestUtils_Encrypt(t *testing.T) {
	d, err := u.AESCFBEncrypt([]byte(str), []byte(key))
	if err != nil {
		panic(err)
	}
	r, _ := u.AESCFBDecrypt(d, []byte(key))
	fmt.Println(hex.EncodeToString(d))
	fmt.Println(string(r))
}

func TestHmacSha256AndSha1(t *testing.T) {
	KEY := []byte(key)
	s := []byte(str)
	r, err := u.HmacSha256(s, KEY)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hex.EncodeToString(r), len(r))

	ok, err := u.HmacSha256Check(s, r, KEY)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("HmacSha256 Error")
	}

	r, err = u.HmacSha1(s, KEY)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hex.EncodeToString(r), len(r))

	ok, err = u.HmacSha1Check(s, r, KEY)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("HmacSha1 Error")
	}
}

func TestHkdfSha256AndSha1(t *testing.T) {
	info := []byte{0x62, 0x72, 0x6f, 0x6f, 0x6b}
	r, s, err := u.HkdfSha256RandomSalt([]byte("hello"), info, 12)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hex.EncodeToString(r), hex.EncodeToString(s), len(r))

	r, err = u.HkdfSha256WithSalt([]byte("hello"), info, info)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hex.EncodeToString(r), len(r))

	r, s, err = u.HkdfSha1RandomSalt([]byte("hello"), info, 12)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hex.EncodeToString(r), hex.EncodeToString(s), len(r))

	r, err = u.HkdfSha1WithSalt([]byte("hello"), info, info)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hex.EncodeToString(r), len(r))

}

func TestTpTokenEncode(t *testing.T) {
	data := "13501242231|73721|1564130996"
	encodeData := u.TokenEncode(data, token_key)
	t.Log("encodeData...", encodeData)
}

func TestTpTokenDecode(t *testing.T) {
	data := "vWRPOvW2wM_GHApp8NgGeA6tJpfPaO_L3mMrC3CBwypScROI4CMPPLCeFK_WgvLC5Mum78G-g8shHF0pW7Utkg=="
	decodeData := u.TokenDecode(data, token_key)
	t.Log("decodeData...", decodeData)
}
