package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

var averageTime time.Duration
var numOfSuccess int
var numOfFail int
var numOfnon200 int
var mtx sync.Mutex
var urlStr string

func sendRequest(client *fasthttp.Client, req *fasthttp.Request, wg *sync.WaitGroup, timeout int) {
	defer wg.Done()
	for {
		resp := fasthttp.AcquireResponse()
		timeStart := time.Now()
		err := client.DoTimeout(req, resp, time.Duration(timeout)*time.Second)
		if err != nil {
			go safeInc(&numOfFail)
			continue
		}

		go setAvarageTime(averageTime + time.Since(timeStart))
		if resp.StatusCode() == fasthttp.StatusOK {
			go safeInc(&numOfSuccess)
		} else {
			go safeInc(&numOfnon200)
		}
	}
}

// safeInc safely incerements the given integer for aoiding race conditions
func safeInc(num *int) {
	mtx.Lock()
	*num++
	mtx.Unlock()
}

func setAvarageTime(newTime time.Duration) {
	mtx.Lock()
	averageTime = newTime
	mtx.Unlock()
}

func main() {

	requestFile := flag.String("r", "", "Path of request file")
	numWorker := flag.Int("w", 500, "Number of worker. Default: 500")
	duration := flag.Int("d", 0, "Test duration. Default: infinite")
	timeout := flag.Int("t", 5, "HTTP request timeout (sec.)")
	https := flag.Bool("https", true, "Enable https protocol")
	flag.Parse()
	if *requestFile == "" {
		fmt.Println("Please specify all arguments!")
		flag.PrintDefaults()
		os.Exit(1) // Exit condition requires non-zero exit code
	}
	content, err := ioutil.ReadFile(*requestFile)
	if err != nil {
		panic(err)
	}

	httpRequest, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(content)))
	if err != nil && err != io.ErrUnexpectedEOF {
		panic(err)
	}

	var wg sync.WaitGroup
	urlStr = "http://" + httpRequest.Host + httpRequest.RequestURI
	if *https {
		urlStr = "https://" + httpRequest.Host + httpRequest.RequestURI
	}

	bodyBytes, err := ioutil.ReadAll(httpRequest.Body)
	if err != nil {
		panic(err)
	}
	req := fasthttp.AcquireRequest()

	req.SetRequestURI(urlStr)
	req.Header.SetMethod(httpRequest.Method)
	if httpRequest.Method == "POST" {
		req.SetBody(bodyBytes)
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
		go sendRequest(client, req, &wg, *timeout)
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
		fmt.Println("Avg. Response Time:", averageTime/time.Duration(numOfSuccess+numOfnon200+numOfFail))
		os.Exit(0) // Exit without error
	}
	fmt.Println("Infinite mode is active. No output will be shown..")
	wg.Wait()
}
