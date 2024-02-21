package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"testing"

	"github.com/vmihailenco/msgpack/v5"
	"github.com/yangkequn/saavuu/config"
	_ "github.com/yangkequn/saavuu/https"
)

func TestHTTPPostMsgPack(t *testing.T) {
	if !config.Cfg.Http.Enable {
		return
	}
	var (
		body  []byte
		err   error
		param = &Demo1{Text: "TestCallAt 10s later"}
		url   = "http://127.0.0.1:" + strconv.Itoa(int(config.Cfg.Http.Port)) + "/API-!demo1"
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

func TestHTTPPostJson(t *testing.T) {
	if !config.Cfg.Http.Enable {
		return
	}
	var (
		body  []byte
		err   error
		param = &Demo1{Text: "TestCallAt 10s later"}
		url   = "http://127.0.0.1:" + strconv.Itoa(int(config.Cfg.Http.Port)) + "/API-!demo1"
		resp  *http.Response
	)
	//create a http context
	body, _ = json.Marshal(param)
	reader := bytes.NewReader(body)

	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, reader)
	req.Header.Add("Content-Type", "application/json")

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

func TestHTTPGetJson(t *testing.T) {
	if !config.Cfg.Http.Enable {
		return
	}
	var (
		body []byte
		err  error
		url  = "http://127.0.0.1:" + strconv.Itoa(int(config.Cfg.Http.Port)) + "/API-!demo1?Text=TestCallAt"
		resp *http.Response
	)
	postBody := []byte(`{"Attach":{"Text":"TestCallAtBody"}}`)
	reader := bytes.NewReader(postBody)

	req, _ := http.NewRequest("POST", url, reader)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
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
