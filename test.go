package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"strconv"
	"strings"
	"net/url"
	"bytes"
)


func ft_my_post(r *http.Response, q_id string, data url.Values, ch chan<-string) {

	req, err := http.NewRequest("POST", "http://185.204.3.165/question/" + q_id, bytes.NewBufferString(data.Encode()))
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}
	s_c := r.Header.Get("Set-Cookie")
	req.Header.Set("Set-Cookie", s_c)

	b := req.Header["Set-Cookie"]
	cookie:= b[0][0:36]
	req.Header.Set("Cookie", cookie)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	
	client := &http.Client{Timeout: time.Second * 1}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading body. ", err)
	}
	if string(body) != "404 - not found" {
		ch <- string(body)
		//fmt.Println(string(body))
		fmt.Println(<- ch)
	} else {
		return
	}
}

func ft_my_get(r *http.Response, q_id string, ch chan<-string) {
	var max_r string
	var max_s string
	var name string
	data := url.Values{}
	req, err := http.NewRequest("GET", "http://185.204.3.165/question/" + q_id, nil)
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}
	s_c := r.Header.Get("Set-Cookie")
	req.Header.Set("Set-Cookie", s_c)

	b := req.Header["Set-Cookie"]
	cookie:= b[0][0:36]
	req.Header.Set("Cookie", cookie)
	
	client := &http.Client{Timeout: time.Second * 1}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading body. ", err)
	}
	if string(body) != "404 - not found" {
		ch <- string(body)
		//fmt.Println(string(body))
		fmt.Println(<-ch)
	} else {
		return
	}
	
	contain := strings.Index(string(body), "<h2>Answer</h2>")
	count := 0
	for flag := 0; flag >= 0; flag++ {
		//count := 0

		if string(body[contain + 50 + 1 + count]) == "i" && string(body[contain + 50 + 13 + count]) == "t" {
			//fmt.Println(string(body[contain + 50 + 13 + count]))
			name = string(body[(contain + 50 + 13 + 12 + count) : (contain + 50 + 13 + 12 + 15 + 1 + count)])
			//fmt.Println(name + " name")//name
			data.Add(name, "test")
			//data = ft_my_post(r , q_id, name, "test", data)
		} else if string(body[contain + 50 + 1 + count]) == "i" && string(body[contain + 50 + 13 + count]) == "r" {
			p_end := strings.Index(string(body[contain + 50 + 13 + count :]), "</p>")

			r_t := 0
			for i := contain + 50 + 13 + count; i < p_end + contain + 50 + 13 + count; i++ {
				v_i := strings.Index(string(body[i + r_t : p_end + contain + 50 + 13 + count]), "value")
				if r_t + v_i > r_t {
					r_t = r_t + v_i
					ko := 0
					for j := 0; j >= 0; j++ {
						if string(body[i + r_t + 7 + j]) == ">" {
							ko = j
							j = -2
						}
					} 
					val := string(body[i + r_t + 7 : i + r_t + 7 + ko - 1])
					if len(val) > len(max_r) {
						max_r = val
					}
				} else {
					break
				}
			}
			max_i := strings.Index(string(body[contain + 50 + 13 + count :]), max_r)
			name = string(body[max_i - 25 + contain + 50 + 13 + count : max_i - 9 + contain + 50 + 13 + count])
			data.Add(name, max_r)
			max_r = ""


		} else if string(body[contain + 50 + 1 + count]) == "s" {

			s_end := strings.Index(string(body[contain + 50 + 1 + count :]), "</p>")

			s := 0
			for i := contain + 50 + 1 + count; i < s_end + contain + 50 + 1 + count; i++ {
				sv_i := strings.Index(string(body[i + s: s_end + contain + 50 + 1 + count]), "value")
				if s + sv_i > s {
					s = s + sv_i
					so := 0
					for j := 0; j >= 0; j++ {
						if string(body[i + s + 7 + j]) == ">" {
							so = j
							j = -2
						}
					} 
					val := string(body[i + s + 7 : i + s + 7 + so - 1])
					if len(val) > len(max_s) {
						max_s = val
					}
				} else {
					break
				}
			}

			name = string(body[contain + 50 + 1 + 13 + count : contain + 50 + 1 + 10 + 19 + count])

			data.Add(name, max_s)
			max_s = ""
		}
		line_end := strings.Index(string(body[contain + 50 + count :]), "\n")
		if string(body[line_end + contain + 50 + 10 + 2 + 1 + count + 1]) == "b" {
			flag = -2
				break
		}
		count = count + line_end + 13
	}
	
	ft_my_post(r , q_id, data, ch)
}

func start_request(r *http.Response, n int) {
	ch := make(chan string)
	for i := 1; i <= 10; i++ {
		go ft_my_get(r, strconv.Itoa(i), ch)
	}
}

func main() {
	n := 6
	r, _ := http.Get("http://185.204.3.165")
	start_request(r, n)
}

