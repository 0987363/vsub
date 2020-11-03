package models

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"time"
)

func LoadShareFromRemote(address string) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second,
	}

	req, err := http.NewRequest("GET", address, nil)
	if err != nil {
		return nil, err
	}

	rsp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer rsp.Body.Close()
	result, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	if rsp.StatusCode != http.StatusOK && rsp.StatusCode != http.StatusNoContent && rsp.StatusCode != http.StatusCreated {
		return nil, Errorf("Request failed. status:%v response:%s", rsp.StatusCode, string(result))
	}

	return result, nil
}
