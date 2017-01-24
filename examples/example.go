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
