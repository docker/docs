package main

import (
	"fmt"

	"github.com/docker/notary/storage"
	"golang.org/x/net/context"
)

func bootstrap(ctx context.Context) error {
	s := ctx.Value("metaStore")
	store, ok := s.(storage.Bootstrapper)
	if !ok {
		return fmt.Errorf("Store does not support bootstrapping.")
	}
	return store.Bootstrap()
}
