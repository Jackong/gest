package gest

import (
	"net/http"
	"menteslibres.net/gosexy/rest"
	"net/http/cookiejar"
	"net/url"
	"crypto/hmac"
	"crypto/sha256"
	"sort"
	"strings"
	"encoding/base64"
	"errors"
	"encoding/json"
	"time"
	"fmt"
)

const (
	SIGN_HEADER = "Content-HMAC"
)

type Gest struct {
	timeDelta int64
	secretKey string
	*rest.Client
}

func New(serverTime int64, secretKey string, header http.Header) (g *Gest) {
	if header == nil {
		header = http.Header{}
	}
	g = new(Gest)
	g.Client = new(rest.Client)
	g.timeDelta = serverTime - time.Now().Unix()
	g.secretKey = secretKey
	g.Header = header
	g.CookieJar, _ = cookiejar.New(nil)
	return
}

func (g *Gest) Get(res interface {}, path string, params url.Values) error {
	if params == nil {
		params = url.Values{}
	}
	params.Set("time", fmt.Sprint("", time.Now().Unix() + g.timeDelta))
	reqSign := g.reqSign("GET", path, params)
	g.Header.Set(SIGN_HEADER, reqSign)
	var response rest.Response
	if err := g.Client.Get(&response, path, params); err != nil {
		return err
	}
	return g.decode(reqSign, response, res)
}

func (g *Gest) Post(res interface {}, path string, params url.Values) error {
	if params == nil {
		params = url.Values{}
	}
	params.Set("time", fmt.Sprint("", time.Now().Unix() + g.timeDelta))
	reqSign := g.reqSign("POST", path, params)
	g.Header.Set(SIGN_HEADER, reqSign)
	var response rest.Response
	if err := g.Client.Post(&response, path, params); err != nil {
		return err
	}
	return g.decode(reqSign, response, res)
}

func (g *Gest) Put(res interface {}, path string, params url.Values) error {
	if params == nil {
		params = url.Values{}
	}
	params.Set("time", fmt.Sprint("", time.Now().Unix() + g.timeDelta))
	reqSign := g.reqSign("PUT", path, params)
	g.Header.Set(SIGN_HEADER, reqSign)
	var response rest.Response
	if err := g.Client.Put(&response, path, params); err != nil {
		return err
	}
	return g.decode(reqSign, response, res)
}

func (g *Gest) Delete(res interface {}, path string, params url.Values) error {
	if params == nil {
		params = url.Values{}
	}
	params.Set("time", string(time.Now().Unix() + g.timeDelta))
	reqSign := g.reqSign("DELETE", path, params)
	g.Header.Set(SIGN_HEADER, reqSign)
	var response rest.Response
	if err := g.Client.Delete(&response, path, params); err != nil {
		return err
	}
	return g.decode(reqSign, response, res)
}

func (g *Gest) reqSign(method, uri string, params url.Values) (string) {
	mac := hmac.New(sha256.New, []byte(g.secretKey))
	mac.Write([]byte(method))
	u, _ := url.Parse(uri)
	mac.Write([]byte(u.Path))
	var keys []string
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		mac.Write([]byte(key + "=" + strings.Join(params[key], ",")))
	}
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func (g *Gest) resSign(reqSign string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(g.secretKey))
	mac.Write([]byte(reqSign))
	mac.Write(body)
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func (g *Gest) decode(reqSign string, response rest.Response, res interface {}) error {
	actual := response.Header.Get(SIGN_HEADER)
	expect := g.resSign(reqSign, response.Body)
	if actual != expect {
		return errors.New("invalid sign")
	}
	return json.Unmarshal(response.Body, res)
}
