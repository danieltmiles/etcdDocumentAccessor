package etcdDocumentAccessor

import (
	"context"
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

func TestAllFulfilled(t *testing.T) {
	RegisterTestingT(t)
	fakeKeysApi := etcdClientWrapper.FakeKeysAPI{}
	fakeKeysApi.ExpectGet("/testkey", "expected result")
	fakeKeysApi.Get(context.Background(), "/testkey", nil)
	err := fakeKeysApi.ExpectationsFulfilled()
	Expect(err).NotTo(HaveOccurred())
}

func TestAllNotFulfilled(t *testing.T) {
	RegisterTestingT(t)
	fakeKeysApi := etcdClientWrapper.FakeKeysAPI{}
	fakeKeysApi.ExpectGet("/testkey", "expected result")
	err := fakeKeysApi.ExpectationsFulfilled()
	Expect(err).To(HaveOccurred())
}
