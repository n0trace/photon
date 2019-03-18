package photon

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

//Response wraps an http.ResponseWriter and implements its interface to be used
type Response struct {
	*http.Response
	ctx          Context
	bodyBytes    []byte
	readBodyOnce sync.Once
}

// Text get text from response.Body
func (resp *Response) Text() (text string, err error) {
	bs, err := resp.Bytes()
	return string(bs), err
}

// Bytes get []byte from response.Body
func (resp *Response) Bytes() (bodyBytes []byte, err error) {
	var buf = bytes.NewBuffer(nil)
	var newReader io.Reader
	resp.readBodyOnce.Do(func() {
		newReader = io.TeeReader(resp.Body, buf)
		resp.Body = ioutil.NopCloser(buf)
	})
	return ioutil.ReadAll(newReader)
}

//JSON parses the resp.Body data and stores the result
func (resp *Response) JSON(i interface{}) (err error) {
	bodyBytes, err := resp.Bytes()
	if err != nil {
		return
	}
	return json.Unmarshal(bodyBytes, &i)
}

//XML parses the resp.Body data and stores the result
func (resp *Response) XML(i interface{}) (err error) {
	bodyBytes, err := resp.Bytes()
	if err != nil {
		return
	}
	return xml.Unmarshal(bodyBytes, &i)
}

//Document get query element from resp.Body
func (resp *Response) Document() (document *goquery.Document, err error) {
	bodyBytes, err := resp.Bytes()
	if err != nil {
		return
	}
	return goquery.NewDocumentFromReader(bytes.NewReader(bodyBytes))
}
