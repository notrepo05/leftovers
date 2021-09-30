package fakes

import (
	"sync"

	"github.com/genevieve/leftovers/nimbus"
)

type TestbedsClient struct {
	DeleteTestbedCall struct {
		sync.Mutex
		CallCount int
		Receives  struct {
			Name string
		}
		Returns struct {
			Error error
		}
		Stub func(string) error
	}
	ListTestbedsCall struct {
		sync.Mutex
		CallCount int
		Returns   struct {
			TestbedSlice []nimbus.Testbed
			Error        error
		}
		Stub func() ([]nimbus.Testbed, error)
	}
}

func (f *TestbedsClient) DeleteTestbed(param1 string) error {
	f.DeleteTestbedCall.Lock()
	defer f.DeleteTestbedCall.Unlock()
	f.DeleteTestbedCall.CallCount++
	f.DeleteTestbedCall.Receives.Name = param1
	if f.DeleteTestbedCall.Stub != nil {
		return f.DeleteTestbedCall.Stub(param1)
	}
	return f.DeleteTestbedCall.Returns.Error
}
func (f *TestbedsClient) ListTestbeds() ([]nimbus.Testbed, error) {
	f.ListTestbedsCall.Lock()
	defer f.ListTestbedsCall.Unlock()
	f.ListTestbedsCall.CallCount++
	if f.ListTestbedsCall.Stub != nil {
		return f.ListTestbedsCall.Stub()
	}
	return f.ListTestbedsCall.Returns.TestbedSlice, f.ListTestbedsCall.Returns.Error
}
