package photon_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/n0trace/photon"
)

var (
	jsonBody = `
	{
		"username":"n0trace",
		"page":"github.com/n0trace"
	}
	`

	xmlBody = `
	<?xml version="1.0" encoding="UTF-8"?><root>
  <username>n0trace</username>
  <page>github.com/n0trace</page>
</root>
	`

	htmlBody = `
	<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <meta http-equiv="X-UA-Compatible" content="ie=edge">
  <title>Document</title>
</head>
<body>
	<h1>hello</h1>
</body>
</html>
	`
)

func TestResponse_Text(t *testing.T) {
	ctx := &photon.Context{}
	resp := new(http.Response)
	resp.Body = ioutil.NopCloser(bytes.NewBufferString(htmlBody))
	ctx.Response = new(photon.Response)
	ctx.Response.Response = resp

	var wantText = htmlBody
	gotText, err := ctx.Text()
	if err != nil {
		t.Errorf("Response.Text() error = %v, wantErr %v", err, nil)
		return
	}
	if gotText != htmlBody {
		t.Errorf("Response.Text() = %v, want %v", gotText, wantText)
	}
}

func TestResponse_Bytes(t *testing.T) {
	ctx := &photon.Context{}
	resp := new(http.Response)
	resp.Body = ioutil.NopCloser(bytes.NewBufferString(htmlBody))
	ctx.Response = new(photon.Response)
	ctx.Response.Response = resp

	var wantBytes = []byte(htmlBody)
	gotBytes, err := ctx.Bytes()
	if err != nil {
		t.Errorf("Response.Bytes() error = %v, wantErr %v", err, nil)
		return
	}

	if !reflect.DeepEqual(gotBytes, wantBytes) {
		t.Errorf("Response.Bytes() = %v, want %v", gotBytes, wantBytes)
	}
}

func TestResponse_Document(t *testing.T) {
	ctx := &photon.Context{}
	resp := new(http.Response)
	resp.Body = ioutil.NopCloser(bytes.NewBufferString(htmlBody))
	ctx.Response = new(photon.Response)
	ctx.Response.Response = resp
	ctx.Response.Response.Request = new(http.Request)

	gotDocument, err := ctx.Document()
	if err != nil {
		t.Errorf("Response.Document() error = %v, wantErr %v", err, nil)
		return
	}
	var wantH1 = "hello"
	if gotDocument.Find("h1").Text() != wantH1 {
		t.Errorf("Response.Document() h1 = %v, wantH1 %v", gotDocument.Find("h1").Text(), wantH1)
		return
	}
}

func TestResponse_JSON(t *testing.T) {
	ctx := &photon.Context{}
	resp := new(http.Response)
	resp.Body = ioutil.NopCloser(bytes.NewBufferString(jsonBody))
	ctx.Response = new(photon.Response)
	ctx.Response.Response = resp

	type User struct {
		Name string `json:"username"`
		Page string `json:"page"`
	}
	var user = new(User)

	if ctx.JSON(user) != nil {
		t.Errorf("Response.JSON() error = %v, wantErr %v", ctx.JSON(user), nil)
	}
	var wantName = "n0trace"
	if user.Name != wantName {
		t.Errorf("Response.JSON() name = %v, wantName %v", user.Name, wantName)
		return
	}

	var wantPage = "github.com/n0trace"
	if user.Page != wantPage {
		t.Errorf("Response.JSON() page = %v, wantPage %v", user.Page, wantPage)
		return
	}
}

func TestResponse_XML(t *testing.T) {
	ctx := &photon.Context{}
	resp := new(http.Response)
	resp.Body = ioutil.NopCloser(bytes.NewBufferString(xmlBody))
	ctx.Response = new(photon.Response)
	ctx.Response.Response = resp

	type User struct {
		Name string `xml:"username"`
		Page string `xml:"page"`
	}
	var user = new(User)

	if ctx.XML(user) != nil {
		t.Errorf("Response.XML() error = %v, wantErr %v", ctx.XML(user), nil)
	}
	var wantName = "n0trace"
	if user.Name != wantName {
		t.Errorf("Response.XML() name = %v, wantName %v", user.Name, wantName)
		return
	}

	var wantPage = "github.com/n0trace"
	if user.Page != wantPage {
		t.Errorf("Response.XML() page = %v, wantPage %v", user.Page, wantPage)
		return
	}

}
