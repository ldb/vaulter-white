package vault

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Vault struct {
	Hostname    string
	AccessToken string
	AppRole     AppRole
}

type AppRole struct {
	RoleId   string `json:"role_id"`
	SecretId string `json:"secret_id"`
}

type Secret struct {
	RequestID     string                     `json:"request_id"`
	LeaseID       string                     `json:"lease_id"`
	Renewable     bool                       `json:"renewable"`
	LeaseDuration int                        `json:"lease_duration"`
	Auth          Auth                       `json:"auth"`
	Data          map[string]json.RawMessage `json:"data"`
}

type Auth struct {
	ClientToken string   `json:"client_token"`
	Accessor    string   `json:"accessor"`
	Policies    []string `json:"policies"`
}

type SecretData map[string]string

func (v *Vault) GetAccessToken() error {
	p, err := json.Marshal(v.AppRole)
	if err != nil {
		return err
	}

	r, err := v.makeRequest("POST", "/v1/auth/approle/login", string(p))
	if err != nil {
		return err
	}

	v.AccessToken = r.Auth.ClientToken
	return nil
}

func (v Vault) GetSecret(secretName string) (secret SecretData, err error) {
	p := fmt.Sprintf("%s%s/%s", "/v1/secret/service/", v.AppRole.RoleId, secretName)

	r, err := v.makeRequest("GET", p, "")
	if err != nil {
		return nil, err
	}

	secret = make(map[string]string)
	for k, v := range r.Data {
		var s string
		err = json.Unmarshal(v, &s)
		if err != nil {
			secret[k] = string(v)
			continue
		}
		secret[k] = s
	}

	return secret, nil
}

func (v Vault) ListSecrets() (secrets []string, err error) {
	p := fmt.Sprintf("%s%s", "/v1/secret/service/", v.AppRole.RoleId)

	r, err := v.makeRequest("LIST", p, "")
	if err != nil {
		return secrets, err
	}

	var secretList []string
	err = json.Unmarshal(r.Data["keys"], &secretList)
	if err != nil {
		return secrets, fmt.Errorf("error listing secrets")
	}

	return secretList, nil
}

func (v Vault) makeRequest(requestType string, path string, params string) (response Secret, err error) {
	url := fmt.Sprintf("%s%s", v.Hostname, path)

	req, err := http.NewRequest(requestType, url, bytes.NewBufferString(params))
	if err != nil {
		return Secret{}, err
	}

	if v.AccessToken != "" {
		req.Header.Set("X-Vault-Token", v.AccessToken)
	}

	client := http.Client{}

	r, err := client.Do(req)
	if err != nil {
		return Secret{}, err
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		return Secret{}, fmt.Errorf("bad response code %d %s", r.StatusCode, r.Body)
	}

	vaultResponse := Secret{}

	d := json.NewDecoder(r.Body)
	d.Decode(&vaultResponse)
	if err != nil {
		return Secret{}, err
	}

	return vaultResponse, nil
}
