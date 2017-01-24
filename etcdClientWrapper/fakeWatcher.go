package etcdClientWrapper

import (
	"errors"
	"fmt"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
	"time"
)

type FakeWatcher struct {
	Nodes []*client.Node
}

func (f *FakeWatcher) ExpectResponse(node *client.Node) {
	f.Nodes = append(f.Nodes, node)
}

func (f *FakeWatcher) ExpectationsWereFulfilled() error {
	if len(f.Nodes) != 0 {
		return errors.New(fmt.Sprintf("Unmet expectations in Watcher: %+v", f.Nodes))
	}
	return nil
}

func (f *FakeWatcher) Next(ctx context.Context) (*client.Response, error) {
	for len(f.Nodes) < 1 {
		<-time.After(time.Millisecond)
	}
	resp := client.Response{
		Node: f.Nodes[0],
	}
	f.Nodes = f.Nodes[1:]
	return &resp, nil
}
