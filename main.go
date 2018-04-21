package main

import (
	"flag"
	"fmt"
	"golang.org/x/sys/unix"
	"log"
	"os"
	"os/exec"
	"strings"
	"github.com/cosmonawt/vaulter-white/vault"
	"github.com/cosmonawt/vaulter-white/conf"
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
		log.Fatal("could not open config: ", err)
	}
	defer file.Close()

	config, err := conf.LoadConfig(file)
	if err != nil {
		log.Fatal("could not load config: ", err)
	}

	command := config.Command
	if len(os.Args) > 3 {
		command = os.Args[3:]
	}

	if command == nil {
		log.Fatal("No Command provided. Please specify in config or provide as argument!")
	}

	vr := vault.AppRole{RoleId: config.RoleID, SecretId: config.SecretId}
	v := vault.Vault{Hostname: config.Host, AccessToken: config.Token, AppRole: vr}

	err = v.GetAccessToken()
	if err != nil {
		log.Fatal("authentication Error: ", err)
	}

	list, err := v.ListSecrets()
	if err != nil {
		log.Fatal("error listing secrets: ", err)
	}

	var secrets = make(map[string]vault.SecretData)
	for _, s := range list {
		secret, err := v.GetSecret(s)
		if err != nil {
			log.Fatal("error getting secret: ", err)
		}
		secrets[s] = secret
	}

	environment := PrepareEnvironment(secrets, config)
	binary, err := exec.LookPath(command[0])
	if err != nil {
		log.Fatal("command not found: ", err)
	}
	unix.Exec(binary, command, environment)
}

func PrepareEnvironment(secrets map[string]vault.SecretData, config conf.Config) []string {
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
