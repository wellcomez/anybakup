package util

import (
	"fmt"
	"os"
	"path/filepath"

	"go.yaml.in/yaml/v3"
)

type Profile struct {
	RepoDir RepoRoot `yaml:"repodir"`
}
type Config struct {
	RepoDir RepoRoot           `yaml:"repodir"`
	Profile map[string]Profile `yaml:"profile"`
	Default Profile            `yaml:"default"`
}

func (r RepoRoot) String() string {
	return string(r)
}
func (r RepoRoot) With(s string) string {
	ret := fmt.Sprintf("%s/%s", r.String(), s)
	return ret
}
func NewConfig() *Config {
	ret := Config{}
	ret.Load()
	return &ret
}
func (c *Config) SetProfile(name string, p Profile)error {
	if name == "" {
		c.Default = p
	} else {
		c.Profile[name] = p
	}
	return c.Save()
}
func (c *Config) GetProfile(name string) *Config {
	if p, ok := c.Profile[name]; ok {
		return &Config{RepoDir: p.RepoDir}
	}
	return nil
}
func (c *Config) Load() error {
	configFilePath, err := c.configfile()
	if err != nil {
		return fmt.Errorf("error getting config file path: %v", err)
	}
	b, err := os.ReadFile(configFilePath)
	if err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}
	return yaml.Unmarshal(b, c)
}
func (c *Config) configfile() (string, error) {
	configDir, err := c.Configdir()
	if err != nil {
		return "", err
	}
	// c.RepoDir = dir
	configFilePath := filepath.Join(configDir, "config.yaml")
	return configFilePath, nil
}

func (Config) Configdir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting home directory: %v", err)
	}

	configDir := filepath.Join(home, ".config", "anybakup")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("error creating config directory: %v", err)
	}
	return configDir, nil
}
func (c *Config) Save() error {
	c.Default.RepoDir = c.RepoDir
	configFilePath, err := c.configfile()
	if err != nil {
		return fmt.Errorf("error getting config file path: %v", err)
	}
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("error marshaling config: %v", err)
	}
	if err := os.WriteFile(configFilePath, data, 0644); err != nil {
		return fmt.Errorf("error writing config file: %v", err)
	}
	return nil
}
