package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

var average_time time.Duration
var numOfSuccess int
var numOfFail int
var numOfnon200 int
var urlStr string

func sendRequest(wg *sync.WaitGroup, req *fasthttp.Request, resp *fasthttp.Response, client *fasthttp.Client) {
	defer wg.Done()
	for {
		time_start := time.Now()
		err := client.Do(req, resp)

		if err != nil {
			numOfFail++
			time.Sleep(1 * time.Second)
			continue
		}
		average_time = average_time + time.Since(time_start)
		statusCode := resp.StatusCode()
		if statusCode == fasthttp.StatusOK {
			numOfSuccess++
		} else {
			numOfnon200++
		}
	}
}

func main() {

	requestFile := flag.String("r", "", "Path of request file")
	numWorker := flag.Int("w", 500, "Number of worker. Default: 500")
	duration := flag.Int("d", 0, "Test duration. Default: infinite")
	protocol := flag.String("p", "https", "Protol to attack. http or https.")
	flag.Parse()
	if *requestFile == "" {
		fmt.Println("Please specify all arguments!")
		flag.PrintDefaults()
		return
	}
	content, err := ioutil.ReadFile(*requestFile)

	if err != nil {
		panic(err)
	}
	httpRequest, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(content)))
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	if *protocol == "http" {
		urlStr = "http://" + httpRequest.Host + httpRequest.RequestURI
	} else {
		urlStr = "https://" + httpRequest.Host + httpRequest.RequestURI
	}

	bodyBytes, _ := ioutil.ReadAll(httpRequest.Body)
	bodyString := string(bodyBytes)
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	req.SetRequestURI(urlStr)
	req.Header.SetMethod(httpRequest.Method)
	if httpRequest.Method == "POST" {
		req.SetBodyString(bodyString)
	}
	for key, element := range httpRequest.Header {
		req.Header.Set(key, strings.Join(element, ","))
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Content-Length", string(rune(httpRequest.ContentLength)))

	client := &fasthttp.Client{
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
	}

	fmt.Println("Starting the action...")
	fmt.Println("The request will be sent:")
	fmt.Println("_________________________________")
	fmt.Println(string(content))
	fmt.Println("_________________________________")
	for i := 0; i < *numWorker; i++ {
		wg.Add(1)
		go sendRequest(&wg, req, resp, client)
	}
	if *duration != 0 {

		time.Sleep(time.Duration(*duration) * time.Second)
		fmt.Println("")
		fmt.Println("---------Results--------------------")
		fmt.Println("")
		fmt.Println("Total Requests Sent:", numOfSuccess+numOfnon200+numOfFail)
		fmt.Println("Number of 200OK Responses:", numOfSuccess)
		fmt.Println("Number of non-200OK Responses:", numOfnon200)
		fmt.Println("Number of Failed Responses:", numOfFail)
		fmt.Println("Avg. Response Time:", average_time/time.Duration(numOfSuccess+numOfnon200+numOfFail))
		os.Exit(1)
	}

	fmt.Println("Infinite mode is active. No output will be shown..")

	wg.Wait()

}
