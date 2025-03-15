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
)

type SNSPublisherConfig struct {
	SNSConfig *SNSConfig `yaml:"sns"`
	SqsConfig *SqsConfig `yaml:"sqs"`
}

func (c *SNSPublisherConfig) Validate() error {
	if c.SNSConfig == nil {
		return fmt.Errorf("config section `sns` must exist")
	}

	if err := c.SNSConfig.Validate("sns"); err != nil {
		return err
	}

	if c.SqsConfig == nil {
		return fmt.Errorf("config section `sqs` must exist")
	}

	if err := c.SqsConfig.Validate("sqs", true); err != nil {
		return err
	}

	return nil
}

func ProvideSNSPublisherConfig() (*SNSPublisherConfig, error) {
	configFile := shared.GetEnvironmentString("SNS_PUBLISHER_CONFIG_FILE", "")

	var cfg *SNSPublisherConfig
	var err error

	if configFile != "" {
		cfg, err = getSNSPublisherConfigFromFile(configFile)
	} else {
		cfg, err = getSNSPublisherConfigFromEnvironment()
	}

	if err != nil {
		return nil, fmt.Errorf("failed creating application config: %v", err)
	}

	if err = cfg.Validate(); err != nil {
		return nil, fmt.Errorf("failed validating application config: %v", err)
	}

	return cfg, nil
}

func getSNSPublisherConfigFromFile(file string) (*SNSPublisherConfig, error) {
	fileContent, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	cfg := &SNSPublisherConfig{}
	if err = yaml.Unmarshal(fileContent, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func getSNSPublisherConfigFromEnvironment() (*SNSPublisherConfig, error) {
	snsConfig, err := getSNSConfigFromEnvironmentVariables()
	if err != nil {
		return nil, err
	}

	return &SNSPublisherConfig{
		SNSConfig: snsConfig,
		SqsConfig: getSqsConfigFromEnvironmentVariables(),
	}, nil
}
