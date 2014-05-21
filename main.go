// package main
// 
// import (
// 	"bufio"
// 	"bytes"
// 	"fmt"
// 	"crypto/tls"
// 	"net"
// 	"net/http"
// 	"net/http/httputil"
// 	"log"
// )
// 
// func fun(line []byte) {
// 	log.Println(line)
// }
// 
// func main () {
// 	tcpConn, err := net.Dial("tcp", "localhost:5556")
// 	if err != nil {
// 		log.Fatal("Error opening tcp connection")
// 	}
// 	cf := &tls.Config{InsecureSkipVerify: true}
// 	ssl := tls.Client(tcpConn, cf)
// 	reader := bufio.NewReader(ssl)
// 	clientConn := httputil.NewClientConn(ssl, reader)
// 	
// 	req, err := http.NewRequest("GET", "stream", nil)
// 	if err != nil {
// 		log.Fatal("Error GETting path")
// 	}
// 
// 	token := "b00t1csse"
// 
// 	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
// 	req.Header.Set("Connection", "Keep-Alive")
// 
// 	_, err = clientConn.Do(req)
// 	if err != nil {
// 		log.Fatal("Error executing request ", err)
// 	}
// 
// 	// reader := bufio.NewReader(resp.Body)
// 	for {
// 	  // if c.stale {
// // 	    c.clientConn.Close()
// // 	    break
// // 	  }
// 
// 	  line, err := reader.ReadBytes('\r')
// 		if err != nil {
// 			log.Fatal("Could not read line", err)
// 		}
// 	  line = bytes.TrimSpace(line)
// 
// 	  fun(line)
// 	}
// }
package main

import (
	"log"
	"net/url"
	"net/http"
	"bufio"
	"bytes"
)

func handler(line []byte) {
	log.Println(string(line))
}
 
func main() {

	url, _ := url.Parse("https://tracker.bootic.net/stream?raw=1")

	var resp *http.Response
	client := &http.Client{}

	var req http.Request
	req.URL = url
	req.Method = "GET"
	req.Header = http.Header{}
	req.Header.Set("Authorization", "Bearer xxx")

	resp, err := client.Do(&req)

	if err != nil {
		log.Fatal("Could not execute request", err)
	}
	if resp.StatusCode != 200 {
		log.Fatal("HTTP error", resp.StatusCode)
	}

	reader := bufio.NewReader(resp.Body)

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			log.Fatal("Could not read line. Broken connection? Reconnect?", err)
		}
		line = bytes.TrimSpace(line)

		if len(line) == 0 {
			continue
		}
		handler(line)
	}
}