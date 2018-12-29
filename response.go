package photon

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type Response struct {
	*http.Response
	ctx *Context
}

func (resp *Response) Text() (text string, err error) {
	bodyBytes, err := resp.Bytes()
	return string(bodyBytes), err
}

func (resp *Response) Bytes() (bodyBytes []byte, err error) {
	bodyBytes, err = ioutil.ReadAll(resp.Response.Body)
	defer resp.Response.Body.Close()
	return
}

func (resp *Response) JSON(i interface{}) (err error) {
	bodyBytes, err := resp.Bytes()
	if err != nil {
		return
	}
	return json.Unmarshal(bodyBytes, &i)
}

func (resp *Response) XML(i interface{}) (err error) {
	bodyBytes, err := resp.Bytes()
	if err != nil {
		return
	}
	return xml.Unmarshal(bodyBytes, &i)
}

func (resp *Response) Document() (*goquery.Document, error) {
	return goquery.NewDocumentFromResponse(resp.Response)
}
