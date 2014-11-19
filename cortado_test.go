package main

import (
	assert "github.com/pilu/miniassert"
	"testing"
)

func TestValidUrl(t *testing.T) {
	valid := []string{
		"http://example.com",
		"http://example.com/foo",
		"http://example.com/?foo=bar",
		"http://example.com#foo",
	}

	invalid := []string{
		"foo",
		"http://",
		"http://example",
		"example.com",
		"http://exam ple.com",
	}

	for _, url := range valid {
		assert.True(t, validUrl(url))
	}

	for _, url := range invalid {
		assert.False(t, validUrl(url))
	}
}
