package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWithOneFromDefaultValueAndOneFromEnvVar(t *testing.T) {
	os.Setenv("PORT", "8089")
	conf, err := New()
	if assert.NoError(t, err) {
		assert.Equal(t, 8089, conf.Port)
		assert.Equal(t, "mongodb://localhost:27017", conf.DBURI)
	}
}

func TestNewWithInvalidEnvVar(t *testing.T) {
	os.Setenv("PORT", "abc")
	_, err := New()
	assert.Error(t, err)
}
