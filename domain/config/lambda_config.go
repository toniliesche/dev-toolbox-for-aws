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
	"slices"
)

type LambdaConfig struct {
	Endpoint       string `yaml:"endpoint"`
	Concurrency    int    `yaml:"concurrency"`
	InvocationType string `yaml:"invocation_type"`
}

func (c *LambdaConfig) Validate(path string) error {
	if c.Endpoint == "" {
		return fmt.Errorf("config value `%s.endpoint` must exist", path)
	}

	if c.Concurrency <= 0 {
		return fmt.Errorf("config value `%s.concurrency` must be greater than 0", path)
	}

	if c.InvocationType == "" {
		return fmt.Errorf("config value `%s.invocation_type` must exist", path)
	}

	if !slices.Contains([]string{"Event", "RequestResponse"}, c.InvocationType) {
		return fmt.Errorf("config value `%s.invocation_type` must be one of [Event, RequestResponse]", path)
	}

	return nil
}

func getLambdaConfigFromEnvironmentVariables() *LambdaConfig {
	return &LambdaConfig{
		Endpoint:       shared.GetEnvironmentString("LAMBDA_ENDPOINT", "http://localhost:8080"),
		Concurrency:    int(shared.GetEnvironmentInt("LAMBDA_CONCURRENCY", 1)),
		InvocationType: shared.GetEnvironmentString("LAMBDA_INVOCATION_TYPE", "RequestResponse"),
	}
}
