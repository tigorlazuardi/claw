package otel

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
)

var Resource *resource.Resource

func init() {
	res, err := resource.New(context.Background(),
		resource.WithFromEnv(),
		resource.WithHost(),
		resource.WithOS(),
		resource.WithContainer(),
	)
	if err != nil {
		Resource = resource.Empty()
		otel.Handle(err)
	}
	Resource, err = resource.Merge(resource.Default(), res)
	if err != nil {
		Resource = resource.Default()
		otel.Handle(err)
	}
}

// SetResourceAttr sets a resource attribute if it is not already set.
func SetResourceAttr(attr attribute.KeyValue) {
	var found bool
	iter := Resource.Iter()
	for iter.Next() {
		if iter.Attribute().Key == attr.Key {
			found = true
			break
		}
	}
	if !found {
		newRes, _ := resource.New(context.Background(), resource.WithAttributes(attr))
		Resource, _ = resource.Merge(Resource, newRes)
	}
}
