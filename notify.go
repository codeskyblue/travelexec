package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func postForm(url string, params map[string]string, files map[string][]byte) (err error) {
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	for name, data := range files {
		part, er := w.CreateFormFile(name, "rep.html")
		if er != nil {
			err = er
			return er
		}
		_, er = part.Write(data)
		if er != nil {
			err = er
			return
		}
	}
	req, err := http.NewRequest("POST", "http://www.baidu.com", body)
	if err != nil {
		return
	}
	if err = w.Close(); err != nil {
		return
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	io.Copy(os.Stderr, res.Body)
	return
}

func sendNotify(msg string, users ...string) (err error) {
	return nil
}
