package consulclient

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

var bodyReader = ioutil.ReadAll

type HTTPKV struct {
	ConsulAddress string
}

func NewHTTPKV(consulAddress string) HTTPKV {
	return HTTPKV{
		ConsulAddress: fmt.Sprintf("%s/consul", consulAddress),
	}
}

func (kv HTTPKV) Address() string {
	return kv.ConsulAddress
}

func (kv HTTPKV) Set(key, value string) error {
	request, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/kv/%s", kv.ConsulAddress, key), strings.NewReader(value))
	if err != nil {
		return fmt.Errorf("new request: %s", err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   1 * time.Second,
				KeepAlive: 1 * time.Second,
			}).Dial,
		},
	}

	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("do request: %s", err)
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %d %s", response.StatusCode, http.StatusText(response.StatusCode))
	}

	defer response.Body.Close()
	body, err := bodyReader(response.Body)
	if err != nil {
		return fmt.Errorf("read body: %s", err)
	}

	if string(body) != "true" {
		return errors.New(fmt.Sprintf("unexpected response body: %s", string(body)))
	}

	return nil
}

func (kv HTTPKV) Get(key string) (string, error) {
	response, err := http.Get(fmt.Sprintf("%s/v1/kv/%s?raw", kv.ConsulAddress, key))
	if err != nil {
		return "", err
	}

	if response.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("key %q not found", key)
	}

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("consul http error: %d %s", response.StatusCode, http.StatusText(response.StatusCode))
	}

	body, err := bodyReader(response.Body)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	return string(body), nil
}
