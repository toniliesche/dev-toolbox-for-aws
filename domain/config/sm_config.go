// MIT License
// Copyright (c) 2025 Toni Liesche
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

package config

import (
	"dev-toolbox-for-aws/domain/shared"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type SecretsManagerConfig struct {
	Account string                  `yaml:"account"`
	Region  string                  `yaml:"region"`
	Secrets []*SecretsManagerSecret `yaml:"secrets"`
}

type SecretsManagerSecret struct {
	Name     string   `yaml:"name"`
	Versions []string `yaml:"versions"`
}

func (c *SecretsManagerConfig) Validate() error {
	if c.Account == "" {
		return fmt.Errorf("config value `account` must exist")
	}

	if c.Region == "" {
		return fmt.Errorf("config value `region` must exist")
	}

	if c.Secrets == nil {
		return fmt.Errorf("config value `secrets` must exist")
	}

	if len(c.Secrets) == 0 {
		return fmt.Errorf("config value `secrets` must contain at least one secret")
	}

	for i, s := range c.Secrets {
		if err := s.Validate(fmt.Sprintf("secrets[%d]", i)); err != nil {
			return err
		}
	}

	return nil
}

func (s *SecretsManagerSecret) Validate(path string) error {
	if s.Name == "" {
		return fmt.Errorf("config value `%s.name` must exist", path)
	}

	if s.Versions == nil {
		return fmt.Errorf("config value `%s.versions` must exist", path)
	}

	if len(s.Versions) == 0 {
		return fmt.Errorf("config value `%s.versions` must contain at least one element", path)
	}

	return nil
}

func ProvideSecretsManagerConfig() (*SecretsManagerConfig, error) {
	configFile := shared.GetEnvironmentString("SECRETS_MANAGER_CONFIG_FILE", "")

	var cfg *SecretsManagerConfig
	var err error

	if configFile != "" {
		cfg, err = getSecretsManagerConfigFromFile(configFile)
	} else {
		cfg, err = getSecretsManagerConfigFromEnvironment()
	}

	if err != nil {
		return nil, fmt.Errorf("failed creating application config: %v", err)
	}

	if err = cfg.Validate(); err != nil {
		return nil, fmt.Errorf("failed validation application config: %v", err)
	}

	return cfg, nil
}

func getSecretsManagerConfigFromFile(file string) (*SecretsManagerConfig, error) {
	fileContent, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	cfg := &SecretsManagerConfig{}
	if err = yaml.Unmarshal(fileContent, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func getSecretsManagerConfigFromEnvironment() (*SecretsManagerConfig, error) {
	secretManagerSecrets := shared.GetEnvironmentString("SECRETS_MANAGER_SECRETS", "")

	secretsMap := make([]*SecretsManagerSecret, 0)
	if secretManagerSecrets != "" {
		secrets := strings.Split(secretManagerSecrets, ";")
		for _, secret := range secrets {
			secretsParts := strings.Split(secret, "=")
			if len(secretsParts) != 2 {
				return nil, fmt.Errorf("invalid secret format: %s, must be key=value1[,value2,value3]", secret)
			}

			secretVersions := strings.Split(secretsParts[1], ",")
			secretObj := &SecretsManagerSecret{
				Name:     secretsParts[0],
				Versions: secretVersions,
			}

			secretsMap = append(secretsMap, secretObj)
		}
	}

	return &SecretsManagerConfig{
		Secrets: secretsMap,
		Region:  shared.GetEnvironmentString("AWS_REGION", "eu-central-1"),
		Account: shared.GetEnvironmentString("AWS_ACCOUNT_ID", "000000000000"),
	}, nil
}
