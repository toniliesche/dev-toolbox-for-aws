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
)

type SqsConfig struct {
	Region    string `yaml:"region"`
	Endpoint  string `yaml:"endpoint"`
	QueueName string `yaml:"queue_name"`
}

func (c *SqsConfig) Validate(path string, allowEmptyQueue bool) error {
	if c.Region == "" {
		return fmt.Errorf("config value `%s.region` must exist", path)
	}

	if c.Endpoint == "" {
		return fmt.Errorf("config value `%s.endpoint` must exist", path)
	}

	if c.QueueName == "" && !allowEmptyQueue {
		return fmt.Errorf("config value `%s.queue_name` must exist", path)
	}

	return nil
}

func getSqsConfigFromEnvironmentVariables() *SqsConfig {
	return &SqsConfig{
		Region:    shared.GetEnvironmentString("AWS_REGION", "eu-central-1"),
		Endpoint:  shared.GetEnvironmentString("SQS_ENDPOINT", "http://localhost:4566"),
		QueueName: shared.GetEnvironmentString("SQS_QUEUE_NAME", "my-queue"),
	}
}
