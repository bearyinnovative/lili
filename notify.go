package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type NotifierType interface {
	Notify(text string)
}

type BCIncommingNotifier struct {
	Domain string
	Token  string
}

func (bc *BCIncommingNotifier) Notify(text string) {
	path := fmt.Sprintf("https://hook.bearychat.com/%s/incoming/%s", bc.Domain, bc.Token)

	dic := map[string]string{
		// "user": "rocry", // TODO: remove this
		"text": text,
	}
	jsonValue, err := json.Marshal(dic)
	FatalIfErr(err)

	body := bytes.NewBuffer(jsonValue)

	client := &http.Client{}
	req, err := http.NewRequest("POST", path, body)
	FatalIfErr(err)

	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	_, err = client.Do(req)
	LogIfErr(err)
}
