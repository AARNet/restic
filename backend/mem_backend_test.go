// DO NOT EDIT, AUTOMATICALLY GENERATED
package backend_test

import (
	"testing"

	"github.com/restic/restic/backend/test"
)

var SkipMessage string

func TestMemBackendCreate(t *testing.T) {
	if SkipMessage != "" {
		t.Skip(SkipMessage)
	}
	test.Create(t)
}

func TestMemBackendOpen(t *testing.T) {
	if SkipMessage != "" {
		t.Skip(SkipMessage)
	}
	test.Open(t)
}

func TestMemBackendCreateWithConfig(t *testing.T) {
	if SkipMessage != "" {
		t.Skip(SkipMessage)
	}
	test.CreateWithConfig(t)
}

func TestMemBackendLocation(t *testing.T) {
	if SkipMessage != "" {
		t.Skip(SkipMessage)
	}
	test.Location(t)
}

func TestMemBackendConfig(t *testing.T) {
	if SkipMessage != "" {
		t.Skip(SkipMessage)
	}
	test.Config(t)
}

func TestMemBackendGetReader(t *testing.T) {
	if SkipMessage != "" {
		t.Skip(SkipMessage)
	}
	test.GetReader(t)
}

func TestMemBackendLoad(t *testing.T) {
	if SkipMessage != "" {
		t.Skip(SkipMessage)
	}
	test.Load(t)
}

func TestMemBackendWrite(t *testing.T) {
	if SkipMessage != "" {
		t.Skip(SkipMessage)
	}
	test.Write(t)
}

func TestMemBackendGeneric(t *testing.T) {
	if SkipMessage != "" {
		t.Skip(SkipMessage)
	}
	test.Generic(t)
}

func TestMemBackendDelete(t *testing.T) {
	if SkipMessage != "" {
		t.Skip(SkipMessage)
	}
	test.Delete(t)
}

func TestMemBackendCleanup(t *testing.T) {
	if SkipMessage != "" {
		t.Skip(SkipMessage)
	}
	test.Cleanup(t)
}
