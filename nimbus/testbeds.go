package nimbus

import (
	"fmt"

	"github.com/notrepo05/leftovers/common"
)

//go:generate faux --interface testbedsClient --output fakes/testbeds_client.go
type testbedsClient interface {
	ListTestbeds() ([]Testbed, error)
	DeleteTestbed(name string) error
}

type Testbeds struct {
	client testbedsClient
	logger logger
}

func NewTestbeds(client testbedsClient, logger logger) Testbeds {
	return Testbeds{
		client: client,
		logger: logger,
	}
}

func (t Testbeds) List(filter string, regex bool) ([]common.Deletable, error) {
	t.logger.Debugln("Listing Nimbus Testbeds...")
	testbeds, err := t.client.ListTestbeds()
	if err != nil {
		return nil, fmt.Errorf("List Nimbus Testbeds: %s", err)
	}

	var resources []common.Deletable
	for _, testbed := range testbeds {
		resource := NewTestbed(t.client, testbed.Name())

		if !common.ResourceMatches(resource.Name(), filter, regex) {
			continue
		}

		proceed := t.logger.PromptWithDetails(resource.Type(), resource.Name())
		if !proceed {
			continue
		}

		resources = append(resources, resource)
	}

	return resources, nil
}

func (t Testbeds) Type() string {
	return "Testbed"
}
