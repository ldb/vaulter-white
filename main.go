package main

import (
	"flag"
	"fmt"
	"golang.org/x/sys/unix"
	"log"
	"os"
	"os/exec"
	"strings"
)

func init() {

}

func main() {
	var c = flag.String("c", "vaulter-white.yaml", "Configuration file")
	flag.Parse()

	if flag.NFlag() < 1 {
		flag.Usage()
	}

	file, err := os.Open(*c)
	if err != nil {
		log.Fatal("Could not open config: ", err)
	}
	defer file.Close()

	config, err := LoadConfig(file)
	if err != nil {
		log.Fatal("Could not load config: ", err)
	}

	command := config.Command
	if len(os.Args) > 3 {
		command = os.Args[3:]
	}

	if command == nil {
		log.Fatal("No Command provided. Please specify in config or provide as argument!")
	}

	vr := VaultAppRole{RoleId: config.RoleID, SecretId: config.SecretId}
	v := Vault{Hostname: config.Host, AccessToken: config.Token, AppRole: vr}

	err = v.GetAccessToken()
	if err != nil {
		log.Fatal("Authentication Error: ", err)
	}

	list, err := v.ListSecrets()
	if err != nil {
		log.Fatal("Error listing secrets: ", err)
	}

	var secrets = make(map[string]VaultSecretData)
	for _, s := range list {
		secret, err := v.GetSecret(s)
		if err != nil {
			log.Fatal("Error getting secret: ", err)
		}
		secrets[s] = secret
	}

	environment := PrepareEnvironment(secrets, config)
	binary, err := exec.LookPath(command[0])
	if err != nil {
		log.Fatal("Command not found: ", err)
	}
	unix.Exec(binary, command, environment)
}

func PrepareEnvironment(secrets map[string]VaultSecretData, config Config) []string {
	environment := os.Environ()
	for name, secret := range secrets {
		for sk, sv := range secret {
			if cv := config.SecretPaths[name][sk]; cv != "" {
				e := fmt.Sprintf("%s=%s", cv, sv)
				environment = append(environment, e)
				continue
			}
			e := fmt.Sprintf("%s=%s", strings.ToUpper(fmt.Sprintf("%s_%s", name, sk)), sv)
			environment = append(environment, e)
		}
	}
	return environment
}
