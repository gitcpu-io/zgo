package zgocrypto

import (
	"fmt"
	"testing"
)

/*
@Time : 2019-03-15 10:16
@Author : rubinus.chu
@File : crypto_test
@project: zgo
*/

var u = New()

func TestAESCFBEncrypt(t *testing.T) {
	bytes, err := u.DecryptDataForWeixinUniond("TKrsd6pprjJQ68SVhhPqkVzmjW9j4389ZbI6ehhqHGo", "AES-128-ECB", "ZhugeAd")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
}

func TestRsaEncrypt(t *testing.T) {
	var publicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDZsfv1qscqYdy4vY+P4e3cAtmv
ppXQcRvrF1cB4drkv0haU24Y7m5qYtT52Kr539RdbKKdLAM6s20lWy7+5C0Dgacd
wYWd/7PeCELyEipZJL07Vro7Ate8Bfjya+wltGK9+XNUIHiumUKULW4KDx21+1NL
AUeJ6PeW+DAkmJWF6QIDAQAB
-----END PUBLIC KEY-----
`)
	r, err := u.RsaDecrypt([]byte("git@github.com/mrkt"), publicKey)
	if err != nil {
		panic(err)
	}
	var privateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDZsfv1qscqYdy4vY+P4e3cAtmvppXQcRvrF1cB4drkv0haU24Y
7m5qYtT52Kr539RdbKKdLAM6s20lWy7+5C0DgacdwYWd/7PeCELyEipZJL07Vro7
Ate8Bfjya+wltGK9+XNUIHiumUKULW4KDx21+1NLAUeJ6PeW+DAkmJWF6QIDAQAB
AoGBAJlNxenTQj6OfCl9FMR2jlMJjtMrtQT9InQEE7m3m7bLHeC+MCJOhmNVBjaM
ZpthDORdxIZ6oCuOf6Z2+Dl35lntGFh5J7S34UP2BWzF1IyyQfySCNexGNHKT1G1
XKQtHmtc2gWWthEg+S6ciIyw2IGrrP2Rke81vYHExPrexf0hAkEA9Izb0MiYsMCB
/jemLJB0Lb3Y/B8xjGjQFFBQT7bmwBVjvZWZVpnMnXi9sWGdgUpxsCuAIROXjZ40
IRZ2C9EouwJBAOPjPvV8Sgw4vaseOqlJvSq/C/pIFx6RVznDGlc8bRg7SgTPpjHG
4G+M3mVgpCX1a/EU1mB+fhiJ2LAZ/pTtY6sCQGaW9NwIWu3DRIVGCSMm0mYh/3X9
DAcwLSJoctiODQ1Fq9rreDE5QfpJnaJdJfsIJNtX1F+L3YceeBXtW0Ynz2MCQBI8
9KP274Is5FkWkUFNKnuKUK4WKOuEXEO+LpR+vIhs7k6WQ8nGDd4/mujoJBr5mkrw
DPwqA3N5TMNDQVGv8gMCQQCaKGJgWYgvo3/milFfImbp+m7/Y3vCptarldXrYQWO
AQjxwc71ZGBFDITYvdgJM1MTqc8xQek1FXn1vfpy2c6O
-----END RSA PRIVATE KEY-----
`)
	b, err := u.RsaDecrypt(r, privateKey)
	if err != nil {
		panic(err)
	}
	fmt.Println(r)
	fmt.Println(b)
}

func TestUtils_Encrypt(t *testing.T) {
	d, err := u.AESCFBEncrypt([]byte("http://www.baidu.com"), []byte("ZhugeAd"))
	if err != nil {
		panic(err)
	}
	r, _ := u.AESCFBDecrypt(d, []byte("ZhugeAd"))
	fmt.Println(string(d))
	fmt.Println(string(r))
}
