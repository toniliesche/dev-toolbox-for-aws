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
	"slices"
)

type SQSLambdaTriggerConfig struct {
	LambdaConfig  *LambdaConfig `yaml:"lambda"`
	SQSConfig     *SQSConfig    `yaml:"sqs"`
	PayloadFormat string        `yaml:"payload_format"`
}

func (c *SQSLambdaTriggerConfig) Validate() error {
	if c.LambdaConfig == nil {
		return fmt.Errorf("config section `lambda` must exist")
	}

	if err := c.LambdaConfig.Validate("lambda"); err != nil {
		return err
	}

	if c.SQSConfig == nil {
		return fmt.Errorf("config section `sqs` must exist")
	}

	if err := c.SQSConfig.Validate("sqs", false); err != nil {
		return err
	}

	if c.PayloadFormat == "" {
		return fmt.Errorf("config value `payload_format` must exist")
	}

	if !slices.Contains([]string{"sqs-event", "pure"}, c.PayloadFormat) {
		return fmt.Errorf("config value `payload_format` must be one of [sqs-event, pure]")
	}

	return nil
}

func ProvideSQSLambdaTriggerConfig() (*SQSLambdaTriggerConfig, error) {
	configFile := shared.GetEnvironmentString("SQS_LAMBDA_TRIGGER_CONFIG_FILE", "")

	var cfg *SQSLambdaTriggerConfig
	var err error

	if configFile != "" {
		cfg, err = getSQSLambdaTriggerConfigFromFile(configFile)
	} else {
		cfg, err = getSQSLambdaTriggerConfigFromEnvironment()
	}

	if err != nil {
		return nil, fmt.Errorf("failed creating application config: %v", err)
	}

	if err = cfg.Validate(); err != nil {
		return nil, fmt.Errorf("failed validation application config: %v", err)
	}

	return cfg, nil
}

func getSQSLambdaTriggerConfigFromFile(file string) (*SQSLambdaTriggerConfig, error) {
	fileContent, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	cfg := &SQSLambdaTriggerConfig{}
	if err = yaml.Unmarshal(fileContent, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func getSQSLambdaTriggerConfigFromEnvironment() (*SQSLambdaTriggerConfig, error) {
	return &SQSLambdaTriggerConfig{
		LambdaConfig:  getLambdaConfigFromEnvironmentVariables(),
		SQSConfig:     getSQSConfigFromEnvironmentVariables(),
		PayloadFormat: shared.GetEnvironmentString("PAYLOAD_FORMAT", "sqs-event"),
	}, nil
}
