package nimbus

import (
	"errors"
	"fmt"

	"github.com/fatih/color"
	"github.com/genevieve/leftovers/app"
	"github.com/genevieve/leftovers/common"
)

type resource interface {
	List(filter string, regex bool) ([]common.Deletable, error)
	Type() string
}

//go:generate faux --interface logger --output fakes/logger.go
type logger interface {
	Printf(message string, a ...interface{})
	Println(message string)
	Debugln(message string)
	PromptWithDetails(resourceType, resourceName string) bool
	NoConfirm()
}

type Leftovers struct {
	logger       logger
	asyncDeleter app.AsyncDeleter
	resources    []resource
}

// NewLeftovers returns a new Leftovers for Nimbus that can be used to list resources,
// list types, or delete resources for the provided account. It returns an error
// if the credentials provided are invalid or a client cannot be created.
func NewLeftovers(logger logger, nimbusUser string) (Leftovers, error) {
	if nimbusUser == "" {
		return Leftovers{}, errors.New("Missing nimbus user.")
	}

	client := NewClient(nimbusUser)
	asyncDeleter := app.NewAsyncDeleter(logger)

	return Leftovers{
		logger:       logger,
		asyncDeleter: asyncDeleter,
		resources: []resource{
			NewTestbeds(client, logger),
		},
	}, nil
}

// List will print all of the resources that match the provided filter.
func (l Leftovers) List(filter string, regex bool) {
	l.logger.NoConfirm()

	var deletables []common.Deletable

	for _, r := range l.resources {
		list, err := r.List(filter, regex)
		if err != nil {
			l.logger.Println(color.YellowString(err.Error()))
		}

		deletables = append(deletables, list...)
	}

	for _, d := range deletables {
		l.logger.Println(fmt.Sprintf("[%s: %s]", d.Type(), d.Name()))
	}
}

// ListByType defaults to List.
func (l Leftovers) ListByType(filter, rType string, regex bool) {
	l.List(filter, regex)
}

// Types will print all the resource types that can
// be deleted on this IaaS.
func (l Leftovers) Types() {
	l.logger.NoConfirm()

	for _, r := range l.resources {
		l.logger.Println(r.Type())
	}
}

// Delete will collect all resources that contain
// the provided filter in the resource's identifier, prompt
// you to confirm deletion (if enabled), and delete those
// that are selected.
func (l Leftovers) Delete(filter string, regex bool) error {
	deletables := [][]common.Deletable{}

	for _, r := range l.resources {
		list, err := r.List(filter, regex)
		if err != nil {
			l.logger.Println(color.YellowString(err.Error()))
		}

		deletables = append(deletables, list)
	}

	return l.asyncDeleter.Run(deletables)
}

// DeleteByType will collect all resources of the provied type that contain
// the provided filter in the resource's identifier, prompt
// you to confirm deletion, and delete those
// that are selected.
func (l Leftovers) DeleteByType(filter, rType string, regex bool) error {
	deletables := [][]common.Deletable{}

	for _, r := range l.resources {
		if r.Type() == rType {
			list, err := r.List(filter, regex)
			if err != nil {
				l.logger.Println(color.YellowString(err.Error()))
			}

			deletables = append(deletables, list)
		}
	}

	return l.asyncDeleter.Run(deletables)
}
