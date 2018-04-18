package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"github.com/cosmonawt/vaulter-white/conf"
	"github.com/cosmonawt/vaulter-white/vault"
)

func TestPrepareEnvironment(t *testing.T) {
	config := conf.Config{
		SecretPaths: map[string]map[string]string{
			"testSecret1": {
				"testKey1": "TEST_VAL1",
			},
		},
	}

	secrets := map[string]vault.VaultSecretData{
		"testSecret1": {
			"testKey1": "TestValue1",
			"testKey2": "TestValue2",
		},
	}

	os.Setenv("TESTENV", "TESTVAL")

	testEnv := PrepareEnvironment(secrets, config)

	assert.Contains(t, testEnv, "TEST_VAL1=TestValue1", "Secrets should be saved according to config")
	assert.Contains(t, testEnv, "TESTSECRET1_TESTKEY2=TestValue2", "Secrets should be saved if config is absent")
	assert.Contains(t, testEnv, "TESTENV=TESTVAL", "Existing environment variables should be included")
}
