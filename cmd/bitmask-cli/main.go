package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	tokenPath = "/tmp/bitmask-token"
)

type cmd string

const (
	start  cmd = "start"
	stop   cmd = "stop"
	status cmd = "status"
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
	if len(args) > 1 {
		usage()
	}
	if len(args) == 1 {
		switch args[0] {
		case string(stop):
			c = stop
		case string(start):
			c = start
		}
	}

	url := "http://127.0.0.1:8000/vpn/" + string(c)

	token, err := getLocalAuthToken()
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Auth-Token", string(token))
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
