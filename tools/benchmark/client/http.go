package client

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type HttpClient struct {
	*http.Client
	serverAddr string
}

func (cli *HttpClient) Do(operation *Operation) error {
	switch operation.Name {
	case OperationTypeGet:
		value, err := cli.get(operation.Key)
		if err != nil {
			return err
		}
		operation.Value = value
	case OperationTypeSet:
		if err := cli.set(operation.Key, operation.Value); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown operation %s", operation.Name)
	}
	return nil
}

func (cli *HttpClient) PipelinedDo(operations []*Operation) error {
	return fmt.Errorf("not implemented")
}

func newHttpClient(serverAddr string) *HttpClient {
	return &HttpClient{
		Client:     &http.Client{Transport: &http.Transport{MaxIdleConnsPerHost: 1}},
		serverAddr: "http://" + serverAddr,
	}
}

func (cli *HttpClient) get(key string) (string, error) {
	resp, err := cli.Get(cli.serverAddr + "/cache/" + key)
	if err != nil {
		return "", err
	}
	if resp.StatusCode == http.StatusNotFound {
		return "", nil
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status code %d", resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (cli *HttpClient) set(key, value string) error {
	req, err := http.NewRequest(http.MethodPut, cli.serverAddr+"/cache/"+key, strings.NewReader(value))
	if err != nil {
		return err
	}
	resp, err := cli.Client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code %d", resp.StatusCode)
	}
	return nil
}
