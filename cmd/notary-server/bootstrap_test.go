package main

import (
	"testing"

	"github.com/docker/notary/tuf/testutils"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

func TestBootstrap(t *testing.T) {
	ctx := context.Background()
	err := bootstrap(ctx)
	require.Error(t, err)
	ctx = context.WithValue(ctx, "metaStore", 1)
	err = bootstrap(ctx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "does not support bootstrapping")
	bs := &testutils.TestBootstrapper{}
	ctx = context.WithValue(ctx, "metaStore", bs)
	err = bootstrap(ctx)
	require.NoError(t, err)
	require.True(t, bs.Booted)
}
