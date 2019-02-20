package main

import (
	"io"
	"io/ioutil"
	"net/http"
)

func makeGETRequest(ip string, url string, sesInfo string) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://" + ip + "/" + url, nil)
	if err != nil {
		return make([]byte, 0), err
	}

	if len(sesInfo) > 0 {
		req.AddCookie(&http.Cookie{Name: "SessionID", Value: sesInfo})
	}

	resp, err := client.Do(req)
	if err != nil {
		return make([]byte, 0), err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return make([]byte, 0), err
	}

	return body, nil
}

func makePOSTRequest(ip string, url string, data io.Reader, sesInfo string, tokInfo string) ([]byte, string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", "http://" + ip + "/" + url, data)
	if err != nil {
		return make([]byte, 0), tokInfo, err
	}

	req.AddCookie(&http.Cookie{Name: "SessionID", Value: sesInfo})
	req.Header.Add("__RequestVerificationToken", tokInfo)

	resp, err := client.Do(req)
	if err != nil {
		return make([]byte, 0), tokInfo, err
	}

	verificationToken := resp.Header.Get("__RequestVerificationToken")

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return make([]byte, 0), verificationToken, err
	}

	return body, verificationToken, nil
} 
