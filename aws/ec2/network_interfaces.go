package ec2

import (
	"fmt"
	awsec2 "github.com/aws/aws-sdk-go/service/ec2"

	"github.com/genevieve/leftovers/common"
)

//go:generate faux --interface networkInterfacesClient --output fakes/network_interfaces_client.go
type networkInterfacesClient interface {
	DescribeNetworkInterfaces(*awsec2.DescribeNetworkInterfacesInput) (*awsec2.DescribeNetworkInterfacesOutput, error)
	DeleteNetworkInterface(*awsec2.DeleteNetworkInterfaceInput) (*awsec2.DeleteNetworkInterfaceOutput, error)
}

type NetworkInterfaces struct {
	client networkInterfacesClient
	logger logger
}

func NewNetworkInterfaces(client networkInterfacesClient, logger logger) NetworkInterfaces {
	return NetworkInterfaces{
		client: client,
		logger: logger,
	}
}

func (e NetworkInterfaces) List(filter string, regex bool) ([]common.Deletable, error) {
	networkInterfaces, err := e.client.DescribeNetworkInterfaces(&awsec2.DescribeNetworkInterfacesInput{})
	if err != nil {
		return nil, fmt.Errorf("Describing EC2 Network Interfaces: %s", err)
	}

	var resources []common.Deletable
	for _, i := range networkInterfaces.NetworkInterfaces {
		r := NewNetworkInterface(e.client, i.NetworkInterfaceId, i.TagSet)

		if !common.ResourceMatches(r.Name(),  filter, regex) {
			continue
		}

		proceed := e.logger.PromptWithDetails(r.Type(), r.Name())
		if !proceed {
			continue
		}

		resources = append(resources, r)
	}

	return resources, nil
}

func (e NetworkInterfaces) Type() string {
	return "ec2-network-interface"
}
