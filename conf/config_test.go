package conf

import (
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	c := `
host: "testHost"
roleId: "testRole"
secretId: "testSecret"`

	config, err := LoadConfig(strings.NewReader(c))
	assert.Nil(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "testHost", config.Host)
	assert.Equal(t, "testRole", config.RoleID)
	assert.Equal(t, "testSecret", config.SecretId)

	c = `
host: "testHost"
some: "nonsensfield"
secretId: "testSecret"`

	config, err = LoadConfig(strings.NewReader(c))
	assert.Nil(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "testHost", config.Host)
	assert.Equal(t, "", config.RoleID)
	assert.Equal(t, "testSecret", config.SecretId)

	c = `
host: "testHost"
some: "nonsensfield"
secretIdEnv: "SECRET"`
	os.Setenv("SECRET", "testSecret")

	config, err = LoadConfig(strings.NewReader(c))
	assert.Nil(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "testSecret", config.SecretId, "should read secret from 'secretIdEnv'")

	c = `
host: "testHost"
some: "nonsensfield"`
	os.Setenv("VAULT_SECRET_ID", "testSecret")

	config, err = LoadConfig(strings.NewReader(c))
	assert.Nil(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "testSecret", config.SecretId, "should read secret from 'VAULT_SECRET_ID'")

	c = ""
	config, err = LoadConfig(strings.NewReader(c))
	assert.NotNil(t, err)
}

func TestSafeLookupEnv(t *testing.T) {
	os.Setenv("VAULT_HOST", "testHost")
	os.Setenv("VAULT_HOST_FALLBACK", "testHostFallback")

	env := safeLookupEnv("VAULT_HOST", "VAULT_HOST_FALLBACK")
	assert.Equal(t, "testHost", env, "should read primary value successfully")

	os.Unsetenv("VAULT_HOST")
	env = safeLookupEnv("VAULT_HOST", "VAULT_HOST_FALLBACK")
	assert.Equal(t, "testHostFallback", env, "should read fallback value successfully")
}
