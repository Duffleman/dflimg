package storageproviders

import (
	"testing"
)

func TestAWSConformsToInterface(t *testing.T) {
	var sp StorageProvider = &AWS{}

	if sp == nil {
		t.Error("not possible")
	}
}

func TestLFSConformsToInterface(t *testing.T) {
	var sp StorageProvider = &LocalFileSystem{}

	if sp == nil {
		t.Error("not possible")
	}
}
