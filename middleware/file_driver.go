package middleware

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

type FileDriver struct {
	cacheDir string
}

func NewFileDriver(cacheDir string) *FileDriver {
	if cacheDir == "" {
		cacheDir = ".cache"
	}
	return &FileDriver{
		cacheDir: cacheDir,
	}
}

func (d *FileDriver) Put(req *http.Request, resp *http.Response) (err error) {
	idxDir := d.idxDir(req)
	err = initDir(idxDir)
	if err != nil {
		return
	}

	respFile, err := os.Create(idxDir + "/info")
	if err != nil {
		err = errors.Wrap(err, "create cache info file error")
		return
	}
	defer respFile.Close()
	defer respFile.Sync()
	bodyFile, err := os.Create(idxDir + "/body")
	if err != nil {
		err = errors.Wrap(err, "create cache body file error")
		return
	}
	defer bodyFile.Close()
	defer bodyFile.Sync()

	var buf1, buf2 = bytes.NewBuffer(nil), bytes.NewBuffer(nil)
	io.Copy(io.MultiWriter(buf1, buf2), resp.Body)
	io.Copy(bodyFile, buf1)
	newRequest := resp.Request
	resp.Request = nil

	err = encodeStdResponse(resp, respFile)
	resp.Body = ioutil.NopCloser(buf2)
	resp.Request = newRequest
	return
}

func (d *FileDriver) Get(req *http.Request) (resp *http.Response, has bool) {
	idx := d.idxDir(req)
	if !dirExist(idx) {
		return
	}

	respFile, err := os.Open(idx + "/info")
	if err != nil {
		err = errors.Wrap(err, "open cache info file error")
		return
	}
	defer respFile.Close()
	bodyFile, err := os.Open(idx + "/body")
	if err != nil {
		err = errors.Wrap(err, "open cache body file error")
		return
	}
	defer bodyFile.Close()

	resp, err = decodeStdResponse(respFile)
	if err != nil {
		err = errors.Wrap(err, "decode response info error")
		return
	}
	buf := bytes.NewBuffer([]byte{})
	_, err = io.Copy(buf, bodyFile)
	if err != nil {
		err = errors.Wrap(err, "decode response body error")
		return
	}
	resp.Body = ioutil.NopCloser(buf)
	resp.Request = req
	resp.TLS = req.TLS
	has = true
	return
}

func (d *FileDriver) idxDir(req *http.Request) string {
	urlStr := req.URL.String()
	idx := fmt.Sprintf("%x", md5.Sum([]byte(urlStr)))
	return fmt.Sprintf("%s/%s/%s", d.cacheDir, idx[:2], idx)
}

func dirExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func initDir(path string) (err error) {
	return os.MkdirAll(path, os.ModePerm)
}

func encodeStdResponse(resp *http.Response, writer io.Writer) (err error) {
	jsonEncoder := json.NewEncoder(writer)
	err = jsonEncoder.Encode(resp)
	return
}

func decodeStdResponse(reader io.Reader) (resp *http.Response, err error) {
	jsonDecoder := json.NewDecoder(reader)
	err = jsonDecoder.Decode(resp)
	return
}
