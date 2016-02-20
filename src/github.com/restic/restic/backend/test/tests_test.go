package test_test

import (
	"errors"

	"github.com/restic/restic/backend"
	"github.com/restic/restic/backend/mem"
	"github.com/restic/restic/backend/test"
)

var be backend.Backend

//go:generate go run ../test/generate_backend_tests.go

func init() {
	test.CreateFn = func() (backend.Backend, error) {
		if be != nil {
			return nil, errors.New("temporary memory backend dir already exists")
		}

		be = mem.New()

		return be, nil
	}

	test.OpenFn = func() (backend.Backend, error) {
		if be == nil {
			return nil, errors.New("repository not initialized")
		}

		return be, nil
	}

	test.CleanupFn = func() error {
		be = nil
		return nil
	}
}
