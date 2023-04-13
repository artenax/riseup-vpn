package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	tokenPath = "/tmp/bitmask-token"
)

type cmd string

const (
	start          cmd = "start"
	stop           cmd = "stop"
	status         cmd = "status"
	gw_get         cmd = "gw/get"
	gw_set         cmd = "gw/set"
	gw_list        cmd = "gw/list"
	transport_get  cmd = "transport/get"
	transport_set  cmd = "transport/set"
	transport_list cmd = "transport/list"
	quit           cmd = "quit"
)

func getLocalAuthToken() ([]byte, error) {
	b, err := ioutil.ReadFile(tokenPath)
	if err != nil {
		return nil, errors.New("no token found")
	}
	return b, nil
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: bitmask-cli [start|stop|status]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()

	c := status
	if len(args) > 2 {
		usage()
	}

	body := io.Reader(nil)
	data := url.Values{}

	switch args[0] {
	case string(stop):
		c = stop
	case string(start):
		c = start
	case string(transport_list):
		c = transport_list
	case string(transport_get):
		c = transport_get
	case string(transport_set):
		c = transport_set
		data.Set("transport", args[1])
	case string(gw_list):
		c = gw_list
	case string(gw_get):
		c = gw_get
	case string(gw_set):
		c = gw_set
	}

	url := "http://127.0.0.1:8000/vpn/" + string(c)

	token, err := getLocalAuthToken()
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	method := getMethodForEndpoint(string(c))
	if method == "POST" {
		body = bytes.NewBuffer([]byte(data.Encode()))
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("X-Auth-Token", string(token))
	if method == "POST" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s: %s\n", string(c), string(b))
}

func getMethodForEndpoint(e string) string {
	if strings.HasSuffix(e, "/set") {
		return "POST"
	} else {
		return "GET"
	}
}
