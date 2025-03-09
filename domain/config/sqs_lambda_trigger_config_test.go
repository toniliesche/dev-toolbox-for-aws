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

func TestValidateSQSLambdaTriggerConfigWithoutLambdaConfigFails(t *testing.T) {
	cfg := &config.SQSLambdaTriggerConfig{
		SQSConfig: &config.SQSConfig{
			Region:           "eu-central-1",
			Endpoint:         "http://localhost:4566",
			QueueName:        "my-queue",
			MessageBatchSize: 10,
			WaitTime:         10,
		},
		PayloadFormat: "sqs-event",
	}

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "config section `lambda` must exist")
}

func TestValidateSQSLambdaTriggerConfigWithoutSQSConfigFails(t *testing.T) {
	cfg := &config.SQSLambdaTriggerConfig{
		LambdaConfig: &config.LambdaConfig{
			Endpoint:       "http://localhost:8080",
			Concurrency:    1,
			InvocationType: "RequestResponse",
		},
		PayloadFormat: "sqs-event",
	}

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "config section `sqs` must exist")
}

func TestValidateSQSLambdaTriggerConfigWithoutPayloadFormatFails(t *testing.T) {
	cfg := &config.SQSLambdaTriggerConfig{
		LambdaConfig: &config.LambdaConfig{
			Endpoint:       "http://localhost:8080",
			Concurrency:    1,
			InvocationType: "RequestResponse",
		},
		SQSConfig: &config.SQSConfig{
			Region:           "eu-central-1",
			Endpoint:         "http://localhost:4566",
			QueueName:        "my-queue",
			MessageBatchSize: 10,
			WaitTime:         10,
		},
	}

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "config value `payload_format` must exist")
}

func TestValidateSQSLambdaTriggerConfigWithInvalidPayloadFormatFails(t *testing.T) {
	cfg := &config.SQSLambdaTriggerConfig{
		LambdaConfig: &config.LambdaConfig{
			Endpoint:       "http://localhost:8080",
			Concurrency:    1,
			InvocationType: "RequestResponse",
		},
		SQSConfig: &config.SQSConfig{
			Region:           "eu-central-1",
			Endpoint:         "http://localhost:4566",
			QueueName:        "my-queue",
			MessageBatchSize: 10,
			WaitTime:         10,
		},
		PayloadFormat: "invalid",
	}

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "config value `payload_format` must be one of [sqs-event, pure]")
}

func TestValidateSQSLambdaTriggerConfigSucceeds(t *testing.T) {
	cfg := &config.SQSLambdaTriggerConfig{
		LambdaConfig: &config.LambdaConfig{
			Endpoint:       "http://localhost:8080",
			Concurrency:    1,
			InvocationType: "RequestResponse",
		},
		SQSConfig: &config.SQSConfig{
			Region:           "eu-central-1",
			Endpoint:         "http://localhost:4566",
			QueueName:        "my-queue",
			MessageBatchSize: 10,
			WaitTime:         10,
		},
		PayloadFormat: "sqs-event",
	}

	err := cfg.Validate()
	assert.NoError(t, err)
}
