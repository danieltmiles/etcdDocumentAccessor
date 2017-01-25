package etcdDocumentAccessor

import (
	"errors"

	"github.com/monsooncommerce/etcdDocumentAccessor/etcdClientWrapper"
	"golang.org/x/net/context"
)

type EtcdDocumentAccessor struct {
	DocumentKey string
	KeysAPI     etcdClientWrapper.KeysAPIWrapper
}

func (e *EtcdDocumentAccessor) Get() (document string, err error) {
	resp, err := e.KeysAPI.Get(context.Background(), e.DocumentKey, nil)
	if err != nil {
		return "", err
	}

	if resp != nil && resp.Node != nil {
		return resp.Node.Value, nil
	}

	return "", errors.New("etcd response was nil or contained node that was nil")
}

func (e *EtcdDocumentAccessor) Set(value string) error {
	resp, err := e.KeysAPI.Set(context.Background(), e.DocumentKey, value, nil)
	if err != nil {
		return err
	}
	if resp == nil || resp.Node == nil {
		return errors.New("etcd response was nil or contained node that was nil")
	}
	return nil
}

func (e *EtcdDocumentAccessor) Watch() <-chan string {
	responseChan := make(chan string, 1)
	go func(responseChan chan string) {
		watcher := e.KeysAPI.Watcher(e.DocumentKey, nil)
		for {
			resp, err := watcher.Next(context.Background())
			if err != nil {
				responseChan <- ""
			} else if resp == nil || resp.Node == nil {
				responseChan <- ""
			} else {
				responseChan <- resp.Node.Value
			}
		}
	}(responseChan)
	return responseChan
}
