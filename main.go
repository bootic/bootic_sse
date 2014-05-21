package main

import (
	"log"
	"net/url"
	"net/http"
	"bufio"
	"bytes"
	"fmt"
)

func handler(line []byte) {
	log.Println(string(line))
}

type eventsChan chan []byte

type Client struct {
	conn *http.Client
	resp *http.Response
	token string
	observers []eventsChan
	url *url.URL
}

func NewClient (urlStr, token string) (client *Client, err error) {

	url, err := url.Parse(urlStr)
	if err != nil {
		return
	}

	client = &Client{
		url: url,
		token: token,
		observers: []eventsChan{},
	}

	err = client.connect()
	if err != nil {
		return
	}

	go client.listen()

	return
}

func (c *Client) Subscribe(observer eventsChan) {
	c.observers = append(c.observers, observer)
}

func (c *Client) connect() (err error) {
	httpConn := &http.Client{}

	var req http.Request
	req.URL = c.url
	req.Method = "GET"
	req.Header = http.Header{}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	resp, err := httpConn.Do(&req)
	if err != nil {
		return
	}

	c.conn = httpConn

	if resp.StatusCode != 200 {
		log.Fatal("HTTP error", resp.StatusCode)
	}
	c.resp = resp

	return
}

func (c *Client) listen() {
	reader := bufio.NewReader(c.resp.Body)

	buffer := make(eventsChan, 20)

	go func() {
		for {
			evt := <-buffer
			for i := range c.observers {
				c.observers[i] <- evt
			}
		}
	}()

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			log.Fatal("Could not read line. Broken connection? Reconnect?", err)
		}
		line = bytes.TrimSpace(line)

		if len(line) == 0 {
			continue
		}
		buffer <- line
	}
}

func main() {
	client, _ := NewClient("https://tracker.bootic.net/stream?raw=1", "b00t1csse")

	events := make(eventsChan)
	client.Subscribe(events)
	
	for {
		log.Println(string(<-events))
	}
	
}