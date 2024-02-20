package test

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"testing"

	"github.com/vmihailenco/msgpack"
	"github.com/yangkequn/saavuu/config"
	_ "github.com/yangkequn/saavuu/https"
)

func TestHTTPPost(t *testing.T) {
	if !config.Cfg.Http.Enable {
		return
	}
	var (
		body  []byte
		err   error
		param = &Demo1{Text: "TestCallAt 10s later"}
		url   = "http://127.0.0.1:" + strconv.Itoa(int(config.Cfg.Http.Port)) + "/API-!api:demo1"
		resp  *http.Response
	)
	//create a http context
	body, _ = msgpack.Marshal(param)
	reader := bytes.NewReader(body)

	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, reader)
	req.Header.Add("Content-Type", "application/octet-stream")

	if resp, err = client.Do(req); err != nil {
		t.Error(err)
		return
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)

	if err != nil {
		t.Error(err)
	} else if string(body) != "hello world" {
		t.Error("result is not hello world")
	}
}
