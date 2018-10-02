# Simple wrapper for common etcd operations


## Constructing your accessor:
An accessor interacts with exactly one key in etcd and must be supplied with a github.com/coreos/etcd/client.KeysAPI
instance when you construct it.

```go
import "github.com/monsooncommerce/etcdDocumentAccessor"

func main(){
        accessor := etcdDocumentAccessor.EtcdDocumentAccessor{
                DocumentKey: "/testkey",
                KeysAPI:     keysApi,
        }
}
```

## Finding a value:
Once you have an instance of an accessor, you can get the value of its key as a string.
```go
value, err := accessor.Get()
```

## Setting a value:
Once you have an instance of accessor, you can set the value of its key as a string.
```go
accessor.Set("this is a whole new value")
```

## Watching:
Once you have an instance of accessor, you may ask it to watch its key and write any new values
into a chan string.
```go
newValueChan := accessor.Watch()
newValue := <-newValueChan // blocks until some other thread writes to your etcd key
```

## Putting it all together (this code is available in examples/example.go):
```go
package main

import (
	"fmt"
	"time"

	"github.com/coreos/etcd/client"
	"github.com/monsooncommerce/etcdDocumentAccessor"
)

func getKeysApi() (client.KeysAPI, error) {
	cfg := client.Config{
		Endpoints:               []string{"http://your-etcd-cluster.domain:2379"},
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}
	conn, err := client.New(cfg)
	if err != nil {
		return nil, err
	}
	return client.NewKeysAPI(conn), nil

}

func main() {
	keysApi, err := getKeysApi()
	if err != nil {
		panic(err.Error())
	}

	// make the accessor
	accessor := etcdDocumentAccessor.EtcdDocumentAccessor{
		DocumentKey: "/testkey",
		KeysAPI:     keysApi,
	}

	// find a value
	value, _ := accessor.Get()
	fmt.Printf("value: %v\n", value)

	// set a value
	accessor.Set("this is a whole new value")
	value, _ = accessor.Get()
	fmt.Printf("I just set this value: \"%v\"\n", value)

	// make a watcher
	resultChan := accessor.Watch()

	// write something for your watcher to find
	go func(accessor *etcdDocumentAccessor.EtcdDocumentAccessor) {
		<-time.After(2 * time.Second)
		accessor.Set("value from inside go func")
	}(&accessor)

	
	// check the watcher
	gotValue := false
	for !gotValue {
		select {
		case value := <-resultChan:
			fmt.Printf("watched value: %v\n", value)
			gotValue = true
		case <-time.After(time.Second):
			fmt.Printf("waited a second, still nothing from the watcher\n")
		}
	}
}
```

## Testing and Mocking
This package provides wrapper interfaces and implementations for a limited set of etcd client library features. These
can be used for mocking. See clientAccessor_test.go for examples:
```go
package etcdDocumentAccessor

import (
	"testing"
	"time"

	"github.com/monsooncommerce/etcdDocumentAccessor/etcdClientWrapper"
	. "github.com/onsi/gomega"
)

func TestGet(t *testing.T) {
	RegisterTestingT(t)
	fakeKeysApi := etcdClientWrapper.FakeKeysAPI{}
	accessor := EtcdDocumentAccessor{
		DocumentKey: "/testkey",
		KeysAPI:     &fakeKeysApi,
	}
	fakeKeysApi.ExpectGet("/testkey", "fake value")
	value, err := accessor.Get()
	Expect(err).NotTo(HaveOccurred())
	Expect(value).To(Equal("fake value"))
}

func TestSet(t *testing.T) {
	RegisterTestingT(t)
	fakeKeysApi := etcdClientWrapper.FakeKeysAPI{}
	accessor := EtcdDocumentAccessor{
		DocumentKey: "/testkey",
		KeysAPI:     &fakeKeysApi,
	}
	err := accessor.Set("test value")
	Expect(err).NotTo(HaveOccurred())
	Expect(fakeKeysApi.Nodes).To(HaveLen(1))
	Expect(fakeKeysApi.Nodes[0].Key).To(Equal("/testkey"))
	Expect(fakeKeysApi.Nodes[0].Value).To(Equal("test value"))
}

func TestWatch(t *testing.T) {
	RegisterTestingT(t)
	fakeKeysApi := etcdClientWrapper.FakeKeysAPI{}
	accessor := EtcdDocumentAccessor{
		DocumentKey: "/testkey",
		KeysAPI:     &fakeKeysApi,
	}
	resultChan := accessor.Watch()
	<-time.After(time.Millisecond)
	fakeKeysApi.ExpectGet("/testkey", "expected result")
	select {
	case result := <-resultChan:
		Expect(result).To(Equal("expected result"))
	case <-time.After(time.Second):
		Expect(1).To(Equal(2), "watcher failed to put anything on its channel")
	}
}
```
