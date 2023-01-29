package http

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// Get connect to other
// url: url
// token: token
// appName: use of scheduler tracking information, optional
func Get(url, appName string) ([]byte, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("App-Name", appName)

	// request do
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	// Judge the return status
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s", resp.Status)
	}

	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func PostFromGateway(url string) ([]byte, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	request, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	// request do
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	// Judge the return status
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s", resp.Status)
	}

	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return result, nil
}
