package repository

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/pkg/errors"
)

type repository struct {
	client  *datastore.Client
	timeout time.Duration
}

// newDatastoreClient returns a new datastore client.
func newDatastoreClient(projectID string, timeout time.Duration) (*datastore.Client, error) {

	// Create context.
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(timeout)*time.Second,
	)
	defer cancel()

	// Create client.
	return datastore.NewClient(ctx, projectID)
}

// NewRepository returns a new repository.
func NewRepository(projectID string, timeout time.Duration) (*repository, error) {

	const operation = "repository.NewDatastoreRepository"

	// Create repository.
	repo := &repository{
		timeout: timeout,
	}

	// Create datastore client.
	client, err := newDatastoreClient(projectID, repo.timeout)
	if err != nil {
		return nil, errors.Wrap(err, operation)
	}

	repo.client = client
	return repo, nil
}
