package main

import (
	"bytes"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
)

func postForm(uri string, params map[string]string, files map[string][]byte) (err error) {
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	for field, value := range params {
		err = w.WriteField(field, value)
		if err != nil {
			return
		}
	}
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
	if err = w.Close(); err != nil {
		return
	}
	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()
	//io.Copy(os.Stderr, res.Body)
	return
}

func sendNotify(msg string, users ...string) (err error) {
	params := map[string]string{
		"username": strings.Join(users, ","),
		"message":  msg,
	}
	err = postForm("http://localhost:8080", params, nil)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
