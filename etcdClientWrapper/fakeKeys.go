package etcdClientWrapper

import (
	"errors"
	"fmt"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
	"time"
)

type KeysAPIWrapper interface {
	Watcher(key string, opts *client.WatcherOptions) client.Watcher
	Get(ctx context.Context, key string, opts *client.GetOptions) (*client.Response, error)
	Set(ctx context.Context, key, val string, opts *client.SetOptions) (*client.Response, error)
}

type FakeKeysAPI struct {
	Nodes []*client.Node
}

func (f *FakeKeysAPI) Get(ctx context.Context, key string, opts *client.GetOptions) (resp *client.Response, err error) {
	resp = &client.Response{}
	if len(f.Nodes) == 0 {
		return nil, errors.New(fmt.Sprintf("Unexpected key Get for %v", key))
	}
	if f.Nodes[0].Key != key {
		return nil, errors.New(fmt.Sprintf("100: Key not found (%v) [39881395]", key))
	}
	resp.Node = f.Nodes[0]
	f.Nodes = f.Nodes[1:]
	return resp, nil
}

func (f *FakeKeysAPI) Set(ctx context.Context, key, val string, opts *client.SetOptions) (*client.Response, error) {
	f.ExpectGet(key, val)
	resp := client.Response{
		Action: "set",
		Node: &client.Node{
			Key:   key,
			Value: val,
		},
	}
	return &resp, nil
}

func (f *FakeKeysAPI) Delete(ctx context.Context, key string, opts *client.DeleteOptions) (*client.Response, error) {
	return nil, nil
}

func (f *FakeKeysAPI) Create(ctx context.Context, key, value string) (*client.Response, error) {
	return nil, nil
}

func (f *FakeKeysAPI) CreateInOrder(ctx context.Context, dir, value string, opts *client.CreateInOrderOptions) (*client.Response, error) {
	return nil, nil
}

func (f *FakeKeysAPI) Update(ctx context.Context, key, value string) (*client.Response, error) {
	return nil, nil
}

func (f *FakeKeysAPI) Watcher(key string, opts *client.WatcherOptions) client.Watcher {
	watcher := &FakeWatcher{}
	go func(watcher *FakeWatcher) {
		for {
			<-time.After(500 * time.Microsecond)
			watcher.Nodes = append(watcher.Nodes, f.Nodes...)
		}
	}(watcher)
	return watcher
}

func (f *FakeKeysAPI) ExpectationsFulfilled() error {
	if len(f.Nodes) != 0 {
		return errors.New("unmet expectations in FakeKeysAPI")
	}
	return nil
}

func (f *FakeKeysAPI) ExpectGet(key, value string) {
	f.Nodes = append(f.Nodes, &client.Node{Key: key, Value: value})
}
