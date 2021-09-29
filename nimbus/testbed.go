package nimbus

import (
	"fmt"
)

type Testbed struct {
	client testbedsClient
	name   string
}

func NewTestbed(client testbedsClient, name string) Testbed {
	return Testbed{
		client: client,
		name:   name,
	}
}

func (t Testbed) Delete() error {
	err := t.client.DeleteTestbed(t.name)
	if err != nil {
		return fmt.Errorf("Delete: %s", err)
	}

	return nil
}

func (t Testbed) Name() string {
	return t.name
}

func (t Testbed) Type() string {
	return "Nimbus Testbed"
}
