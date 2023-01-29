package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
	"os"
)

const (
	testUrl			= "http://185.204.3.165"
	testEndpoint	= "http://185.204.3.165/question/"
)

func sendRequest(r *http.Response, questionID string, data url.Values) {

	req, err := http.NewRequest(http.MethodPost, testEndpoint+questionID, bytes.NewBufferString(data.Encode()))
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}
	setCookie := r.Header.Get("Set-Cookie")
	req.Header.Set("Set-Cookie", setCookie)

	b := req.Header["Set-Cookie"]
	cookie := b[0][0:36]
	req.Header.Set("Cookie", cookie)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: time.Second * 1}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error reading response: ", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading body: ", err)
	}

	if string(body) != "404 - not found" {
		fmt.Println(string(body))
	}
}

func parseTest(r *http.Response, questionID string) {
	var (
		maxRadio  string
		maxSelect string
		name      string
	)

	data := url.Values{}
	req, err := http.NewRequest("GET", "http://185.204.3.165/question/"+questionID, nil)
	if err != nil {
		log.Fatal("Error reading request: ", err)
	}
	setCookie := r.Header.Get("Set-Cookie")
	req.Header.Set("Set-Cookie", setCookie)

	b := req.Header["Set-Cookie"]
	cookie := b[0][0:36]
	req.Header.Set("Cookie", cookie)

	client := &http.Client{Timeout: time.Second * 1}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error reading response: ", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading body: ", err)
	}

	if string(body) == "404 - not found" {
		return
	}
	fmt.Println(string(body))

	contain := strings.Index(string(body), "<h2>Answer</h2>")
	count := 0
	for {
		if string(body[contain+50+1+count]) == "i" && string(body[contain+50+13+count]) == "t" {
			//fmt.Println(string(body[contain + 50 + 13 + count]))
			name = string(body[(contain + 50 + 13 + 12 + count):(contain + 50 + 13 + 12 + 15 + 1 + count)])
			//fmt.Println(name + " name")//name
			data.Add(name, "test")
			//data = sendRequest(r , questionID, name, "test", data)
		} else if string(body[contain+50+1+count]) == "i" && string(body[contain+50+13+count]) == "r" {
			pEnd := strings.Index(string(body[contain+50+13+count:]), "</p>")

			sumValueIndex := 0
			for i := contain + 50 + 13 + count; i < pEnd+contain+50+13+count; i++ {
				valueIndex := strings.Index(string(body[i+sumValueIndex:pEnd+contain+50+13+count]), "value")
				if sumValueIndex+valueIndex > sumValueIndex {
					sumValueIndex += valueIndex
					savePosition := 0
					for j := 0; j >= 0; j++ {
						if string(body[i+sumValueIndex+7+j]) == ">" {
							savePosition = j
							j = -2
						}
					}
					val := string(body[i+sumValueIndex+7 : i+sumValueIndex+7+savePosition-1])
					if len(val) > len(maxRadio) {
						maxRadio = val
					}
				} else {
					break
				}
			}
			maxI := strings.Index(string(body[contain+50+13+count:]), maxRadio)
			name = string(body[maxI-25+contain+50+13+count : maxI-9+contain+50+13+count])
			data.Add(name, maxRadio)
			maxRadio = ""

		} else if string(body[contain+50+1+count]) == "s" {

			sEnd := strings.Index(string(body[contain+50+1+count:]), "</p>")

			s := 0
			for i := contain + 50 + 1 + count; i < sEnd+contain+50+1+count; i++ {
				svI := strings.Index(string(body[i+s:sEnd+contain+50+1+count]), "value")
				if s+svI > s {
					s = s + svI
					position := 0
					for j := 0; j >= 0; j++ {
						if string(body[i+s+7+j]) == ">" {
							position = j
							break
						}
					}
					res := string(body[i+s+7 : i+s+7+position-1])
					if len(res) > len(maxSelect) {
						maxSelect = res
					}
					continue
				}
				break
			}

			name = string(body[contain+50+1+13+count : contain+50+1+10+19+count])

			data.Add(name, maxSelect)
			maxSelect = ""
		}

		lineEnd := strings.Index(string(body[contain+50+count:]), "\n")
		if string(body[lineEnd+contain+50+10+2+1+count+1]) == "b" {
			break
		}
		count += lineEnd + 13
	}

	sendRequest(r, questionID, data)
}

func startTest(r *http.Response, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 1; i <= 10; i++ {
		parseTest(r, strconv.Itoa(i))
	}
}

func main() {
	req, _ := http.Get(testUrl)

	var wg sync.WaitGroup

	workerPool := 1

	if len(os.Args) == 2 {
		workerPool, _ = strconv.Atoi(os.Args[1])
	} else if len(os.Args) > 2 {
		log.Fatal("Error too much args: ")
	}

	for i := 0; i <= workerPool - 1; i++ {
		wg.Add(1)
		go startTest(req, &wg)
	}
	wg.Wait()
}
