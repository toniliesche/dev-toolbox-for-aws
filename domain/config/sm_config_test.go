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

package config_test

import (
	"dev-toolbox-for-aws/domain/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateSecretsManagerConfigWithoutAccountFails(t *testing.T) {
	cfg := &config.SecretsManagerConfig{
		Region: "eu-central",
		Secrets: []*config.SecretsManagerSecret{
			{
				Name:     "secret",
				Versions: []string{"1"},
			},
		},
	}

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "config value `account` must exist")
}

func TestValidateSecretsManagerConfigWithoutRegionFails(t *testing.T) {
	cfg := &config.SecretsManagerConfig{
		Account: "000000000000",
		Secrets: []*config.SecretsManagerSecret{
			{
				Name:     "secret",
				Versions: []string{"1"},
			},
		},
	}

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "config value `region` must exist")
}

func TestValidateSecretsManagerConfigWithoutSecretsListFails(t *testing.T) {
	cfg := &config.SecretsManagerConfig{
		Account: "000000000000",
		Region:  "eu-central-1",
	}

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "config value `secrets` must exist")
}

func TestValidateSecretsManagerConfigWithEmptySecretsListFails(t *testing.T) {
	cfg := &config.SecretsManagerConfig{
		Account: "000000000000",
		Region:  "eu-central-1",
		Secrets: []*config.SecretsManagerSecret{},
	}

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "config value `secrets` must contain at least one secret")
}

func TestValidateSecretsManagerConfigSucceeds(t *testing.T) {
	cfg := &config.SecretsManagerConfig{
		Account: "000000000000",
		Region:  "eu-central-1",
		Secrets: []*config.SecretsManagerSecret{
			{
				Name:     "secret",
				Versions: []string{"1"},
			},
		},
	}

	err := cfg.Validate()
	assert.NoError(t, err)
}

func TestValidateSecretsManagerSecretWithoutNameFails(t *testing.T) {
	cfg := &config.SecretsManagerSecret{
		Name:     "",
		Versions: []string{"1"},
	}

	err := cfg.Validate("secrets[0]")
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "config value `secrets[0].name` must exist")
}

func TestValidateSecretsManagerSecretWithoutVersionsListFails(t *testing.T) {
	cfg := &config.SecretsManagerSecret{
		Name: "secret",
	}

	err := cfg.Validate("secrets[0]")
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "config value `secrets[0].versions` must exist")
}

func TestValidateSecretsManagerSecretWithEmptyVersionsListFails(t *testing.T) {
	cfg := &config.SecretsManagerSecret{
		Name:     "secret",
		Versions: []string{},
	}

	err := cfg.Validate("secrets[0]")
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "config value `secrets[0].versions` must contain at least one element")
}

func TestValidateSecretsManagerSecretSucceeds(t *testing.T) {
	cfg := &config.SecretsManagerSecret{
		Name:     "secret",
		Versions: []string{"1"},
	}

	err := cfg.Validate("secrets[0]")
	assert.NoError(t, err)
}
