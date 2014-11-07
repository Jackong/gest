package gest_test

import (
	. "github.com/onsi/ginkgo"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"log"
)

var _ = Describe("Gest", func() {

		var sign = func (secretKey string, data string) string {
			mac := hmac.New(sha256.New, []byte(secretKey))
			mac.Write([]byte(data))
			return base64.StdEncoding.EncodeToString(mac.Sum(nil))
		}

		Describe("sign data", func () {
		    Context("With Chinese chars", func () {
				It("will be printed", func () {
					log.Println(sign("test", "string contains 中文"))
				})
		    })
			Context("Without Chinese chars", func () {
			    It("will be printed", func () {
			        log.Println(sign("test", "string without Chinese"))
			    })
			})
		})
})
