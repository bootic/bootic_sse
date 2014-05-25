package bootic_sse

import (
	"bufio"
	"bytes"
	"fmt"
	data "github.com/bootic/bootic_go_data"
	"log"
	"net/http"
	"net/url"
)

func handler(line []byte) {
	log.Println(string(line))
}

type Client struct {
	conn      *http.Client
	resp      *http.Response
	token     string
	observers []data.EventsChannel
	url       *url.URL
}

func NewClient(urlStr, token string) (client *Client, err error) {

	url, err := url.Parse(urlStr)
	if err != nil {
		return
	}

	client = &Client{
		url:       url,
		token:     token,
		observers: []data.EventsChannel{},
	}

	err = client.connect()
	if err != nil {
		return
	}

	go client.listen()

	return
}

func (c *Client) Subscribe(observer data.EventsChannel) {
	c.observers = append(c.observers, observer)
}

func (c *Client) SubscribeToType(observer data.EventsChannel, topic string) {
	log.Println("SSE client does not filter by topic yet", topic)
	c.Subscribe(observer)
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

	buffer := make(data.EventsChannel, 20)

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
		event, err := data.DecodeJSON(line)
		if err != nil {
			log.Println("Could not decode", string(line))
			continue
		}
		buffer <- event
	}
}
