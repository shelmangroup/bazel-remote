package config

import (
	"fmt"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// Config provides the configuration
type Config struct {
	Host               string `yaml:"host"`
	Port               int    `yaml:"port"`
	Dir                string `yaml:"dir"`
	MaxSize            int    `yaml:"max_size"`
	HtpasswdFile       string `yaml:"htpasswd_file"`
	TLSCertFile        string `yaml:"tls_cert_file"`
	TLSKeyFile         string `yaml:"tls_key_file"`
	GoogleCloudStorage *struct {
		Bucket                string `yaml:"bucket"`
		UseDefaultCredentials bool   `yaml:"use_default_credentials"`
		JSONCredentialsFile   string `yaml:"json_credentials_file"`
	} `yaml:"gcs_proxy"`
	HTTPBackend *struct {
		BaseURL string `yaml:"base_url"`
	} `yaml:"http_proxy"`
}

// NewConfig ...
func NewConfig(dir string, maxSize int, host string, port int, htpasswdFile string,
	tlsCertFile string, tlsKeyFile string) *Config {
	return &Config{
		Host:               host,
		Port:               port,
		Dir:                dir,
		MaxSize:            maxSize,
		HtpasswdFile:       htpasswdFile,
		TLSCertFile:        tlsCertFile,
		TLSKeyFile:         tlsKeyFile,
		GoogleCloudStorage: nil,
		HTTPBackend:        nil,
	}
}

// NewConfigFromYamlFile ...
func NewConfigFromYamlFile(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to open config file '%s': %v", path, err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("Failed to read config file '%s': %v", path, err)
	}

	return newConfigFromYaml(data)
}

func newConfigFromYaml(data []byte) (*Config, error) {
	c := Config{}
	err := yaml.Unmarshal(data, &c)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse YAML config: %v", err)
	}

	if c.Dir == "" {
		return nil, fmt.Errorf("The 'dir' key is required in the YAML config %v", c)
	}

	if c.MaxSize == 0 {
		return nil, fmt.Errorf("The 'max_size' key is required in the YAML config")
	}

	if (c.TLSCertFile != "" && c.TLSKeyFile == "") || (c.TLSCertFile == "" && c.TLSKeyFile != "") {
		return nil, fmt.Errorf("When enabling TLS, one must specify both keys " +
			"'tls_key_file' and 'tls_cert_file' in the YAML config")
	}

	return &c, nil
}
