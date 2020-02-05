package ioutils

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCreateDirIfMissing(t *testing.T) {
	missed, err := CreateDirIfMissing("./test/child")
	if !assert.NoError(t, err) {
		t.Skipped()
	}

	if !assert.Equal(t, missed, true) {
		t.Skipped()
	}

	missed, err = CreateDirIfMissing("./test")
	if assert.NoError(t, err) {
		assert.Equal(t, missed, false)
	}

	os.RemoveAll("./test")
}
