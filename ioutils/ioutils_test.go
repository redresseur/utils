package ioutils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestCopy(t *testing.T) {
	err := WriteFile("test.data", []byte("hello world"), 0600)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	defer os.Remove("test.data")

	err = Copy("test.data", "test.data.bak")
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	defer os.Remove("test.data.bak")

	data, err := ReadFile("test.data.bak")
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	assert.Equal(t, []byte("hello world"), data)
}
