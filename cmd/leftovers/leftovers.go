package main

import (
	"errors"
	"fmt"

	"github.com/notrepo05/leftovers/app"
	"github.com/notrepo05/leftovers/aws"
	"github.com/notrepo05/leftovers/azure"
	"github.com/notrepo05/leftovers/gcp"
	"github.com/notrepo05/leftovers/nimbus"
	"github.com/notrepo05/leftovers/nsxt"
	"github.com/notrepo05/leftovers/openstack"
	"github.com/notrepo05/leftovers/vsphere"
)

type leftovers interface {
	Delete(filter string, regex bool) error
	DeleteByType(filter, rType string, regex bool) error
	List(filter string, regex bool)
	ListByType(filter, rType string, regex bool)
	Types()
}

func GetLeftovers(logger *app.Logger, o app.Options) (leftovers, error) {
	var (
		l   leftovers
		err error
	)

	switch o.IAAS {
	case app.AWS:
		l, err = aws.NewLeftovers(logger, o.AWSAccessKeyID, o.AWSSecretAccessKey, o.AWSSessionToken, o.AWSRegion)
	case app.Azure:
		l, err = azure.NewLeftovers(logger, o.AzureClientID, o.AzureClientSecret, o.AzureSubscriptionID, o.AzureTenantID)
	case app.GCP:
		fmt.Println("in gcp")
		l, err = gcp.NewLeftovers(logger, o.GCPServiceAccountKey)
	case app.NSXT:
		l, err = nsxt.NewLeftovers(logger, o.NSXTManagerHost, o.NSXTUser, o.NSXTPassword)
	case app.VSphere:
		if o.NoConfirm {
			err = errors.New("--no-confirm is not supported for vSphere.")
		} else {
			l, err = vsphere.NewLeftovers(logger, o.VSphereIP, o.VSphereUser, o.VSpherePassword, o.VSphereDC)
		}
	case app.Nimbus:
		l, err = nimbus.NewLeftovers(logger, o.NimbusUser)
	case app.Openstack:
		l, err = openstack.NewLeftovers(logger, o.OpenstackAuthUrl, o.OpenstackUsername, o.OpenstackPassword, o.OpenstackDomain, o.OpenstackTenant, o.OpenstackRegion)
	default:
		err = errors.New("Missing or unsupported BBL_IAAS.")
	}

	if err != nil {
		return nil, err
	}
	return l, nil
}
