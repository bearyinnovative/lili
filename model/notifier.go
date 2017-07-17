package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	. "../util"
)

type NotifierType interface {
	Notify(text string)
}

type BCIncommingNotifier struct {
	Domain    string
	Token     string
	ToUser    string
	ToChannel string
}

func DefaultChannelNotifier(to string) NotifierType {
	return &BCIncommingNotifier{
		Domain:    "=bw52O",
		Token:     "08c0d225efc37cb33d31d089b91233d1",
		ToChannel: to,
	}
}

func DefaultUserNotifier(to string) NotifierType {
	return &BCIncommingNotifier{
		Domain: "=bw52O",
		Token:  "08c0d225efc37cb33d31d089b91233d1",
		ToUser: to,
	}
}

func (bc *BCIncommingNotifier) Notify(text string) {
	path := fmt.Sprintf("https://hook.bearychat.com/%s/incoming/%s", bc.Domain, bc.Token)

	dic := map[string]string{
		"text": text,
	}
	if bc.ToUser != "" {
		dic["user"] = bc.ToUser
	}
	if bc.ToChannel != "" {
		dic["channel"] = bc.ToChannel
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
