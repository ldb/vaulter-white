package vault

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVault_GetAccessToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "HTTP Method should be correct")
		assert.Equal(t, "/v1/auth/approle/login", r.RequestURI, "URI should be correct")

		login := VaultAppRole{}
		d := json.NewDecoder(r.Body)
		d.Decode(&login)

		if login.SecretId == "" {
			w.WriteHeader(400)
			return
		}

		assert.Equal(t, "roleId", login.RoleId, "RoleID should be transmitted correctly")
		assert.Equal(t, "secretId", login.SecretId, "SecretID should be transmitted correctly")

		w.Write([]byte(`{"auth": {"client_token": "accessToken"}}`))
	}))
	defer ts.Close()

	ar := VaultAppRole{RoleId: "roleId", SecretId: "secretId"}
	v := Vault{Hostname: ts.URL, AppRole: ar}

	v.GetAccessToken()
	assert.Equal(t, "accessToken", v.AccessToken, "Correct Access Token should be returned")

	v.AppRole.SecretId = ""
	e := v.GetAccessToken()
	assert.NotNil(t, e, "should produce error if no SecretId was passed")
}

func TestVault_ListSecrets(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "LIST", r.Method, "HTTP Method should be correct")
		assert.Equal(t, "/v1/secret/service/roleId", r.RequestURI, "URI should be correct")
		assert.Equal(t, "accessToken", r.Header.Get("X-Vault-Token"), "should contain correct Access Token")

		w.Write([]byte(`{"data": {"keys": ["testSecret1", "testSecret2"]}}`))
	}))
	defer ts.Close()

	ar := VaultAppRole{RoleId: "roleId", SecretId: "secretId"}
	v := Vault{Hostname: ts.URL, AppRole: ar, AccessToken: "accessToken"}

	s, err := v.ListSecrets()
	assert.Nil(t, err, "should not produce error")

	expected := []string{
		"testSecret1",
		"testSecret2",
	}

	assert.Equal(t, expected, s, "should return correct secrets")
}

func TestVault_GetSecret(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "HTTP Method should be correct")
		assert.Equal(t, "/v1/secret/service/roleId/testSecret1", r.RequestURI, "URI should be correct")
		assert.Equal(t, "accessToken", r.Header.Get("X-Vault-Token"), "should contain correct Access Token")

		w.Write([]byte(`{"data": {"secretKey1": "secretValue1", "secretKey2": {"secretKey2Sub1":"secret2Sub1Value"}, "secretKey3": ["1",2]}}`))
	}))
	defer ts.Close()

	ar := VaultAppRole{RoleId: "roleId", SecretId: "secretId"}
	v := Vault{Hostname: ts.URL, AppRole: ar, AccessToken: "accessToken"}

	s, err := v.GetSecret("testSecret1")
	assert.Nil(t, err, "should not produce error")

	expected := VaultSecretData{
		"secretKey1": "secretValue1",
		"secretKey2": `{"secretKey2Sub1":"secret2Sub1Value"}`,
		"secretKey3": "[\"1\",2]",
	}

	assert.Equal(t, expected, s, "should unmarshal secret values correctly")
}
