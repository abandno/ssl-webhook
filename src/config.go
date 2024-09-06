package src

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	ContextPath string `yaml:"CONTEXT_PATH"`
	// 回调http端口，缺省 10010
	Port              int    `yaml:"PORT"`
	CallbackToken     string `yaml:"CALLBACK_TOKEN"`
	NginxCertBasePath string `yaml:"NGINX_CERT_BASE_PATH"`
}

var _config *Config

func GetConfig() *Config {
	if _config != nil {
		return _config
	}
	config := Config{
		ContextPath:       "/sslwebhook",
		Port:              10010,
		CallbackToken:     os.Getenv("CALLBACK_TOKEN"),
		NginxCertBasePath: "/etc/nginx/cert",
	}
	configFile, err := os.Open("config.yaml")
	if os.IsNotExist(err) {
		configFile, err = os.Open("config.yml")
	}
	if err != nil {
		//log.Fatal(err)
		log.Println("=> config.yaml not found, use default config")
	} else {
		log.Printf("=> `%s` is found\n", configFile.Name())
		bytes, _ := ioutil.ReadAll(configFile)
		yaml.Unmarshal(bytes, &config)
	}
	defer configFile.Close()

	log.Println(config.ContextPath)
	//log.Println(config.CallbackToken)
	log.Println(config.NginxCertBasePath)
	_config = &config
	return _config
}
