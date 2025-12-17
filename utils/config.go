package utils

import (
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Template Template
	Keys     []Key     `json:"keys" toml:"keys"`
	Lint     Lint      `json:"lint" toml:"lint"`
	LLM      LLMConfig `json:"llm" toml:"llm"`
}

type Lint struct {
	MaxSubjectLength int `json:"max_subject_length" toml:"max_subject_length"`
}

type LLMConfig struct {
	Backend string `json:"backend" toml:"backend"`
	Model   string `json:"model" toml:"model"`
}

type Template struct {
	Region string `json:"region" toml:"region"`
	Normal string `json:"normal" toml:"normal"`
}

type Key struct {
	Prefix      string `json:"prefix" toml:"prefix"`
	Description string `json:"description" toml:"description"`
}

func (k Key) FilterValue() string {
	return k.Prefix
}

func GenerateConfigFromFile(path string) (Config, error) {
	var c Config
	_, err := toml.DecodeFile(path, &c)
	return c, err
}

func GenerateConfig(data string) (Config, error) {
	var c Config
	_, err := toml.Decode(data, &c)
	return c, err
}

func LoadConfig(embeddedConfig string) (Config, error) {
	cfg, err := GenerateConfig(embeddedConfig)
	if err != nil {
		return Config{}, err
	}

	setDefaults(&cfg)

	repoConfigPath := os.ExpandEnv("$PWD/.overcommit.toml")
	if _, err := os.Stat(repoConfigPath); os.IsNotExist(err) {
		return cfg, nil
	}

	repoCfg, err := GenerateConfigFromFile(repoConfigPath)
	if err != nil {
		return cfg, nil
	}

	return mergeConfigs(cfg, repoCfg), nil
}

func setDefaults(cfg *Config) {
	if cfg.Lint.MaxSubjectLength == 0 {
		cfg.Lint.MaxSubjectLength = 50
	}
	if cfg.LLM.Backend == "" {
		cfg.LLM.Backend = "ollama"
	}
	if cfg.LLM.Model == "" {
		cfg.LLM.Model = "tinyllama"
	}
}

func mergeConfigs(base, repo Config) Config {
	if len(repo.Keys) > 0 {
		base.Keys = repo.Keys
	}
	if repo.Template.Region != "" {
		base.Template.Region = repo.Template.Region
	}
	if repo.Template.Normal != "" {
		base.Template.Normal = repo.Template.Normal
	}
	if repo.Lint.MaxSubjectLength > 0 {
		base.Lint.MaxSubjectLength = repo.Lint.MaxSubjectLength
	}
	if repo.LLM.Backend != "" {
		base.LLM.Backend = repo.LLM.Backend
	}
	if repo.LLM.Model != "" {
		base.LLM.Model = repo.LLM.Model
	}
	return base
}
