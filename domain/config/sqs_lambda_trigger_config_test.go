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

func TestValidateSqsLambdaTriggerConfigWithoutLambdaConfigFails(t *testing.T) {
	cfg := &config.SqsLambdaTriggerConfig{
		SqsConfig: &config.SqsConfig{
			Region:    "eu-central-1",
			Endpoint:  "http://localhost:4566",
			QueueName: "my-queue",
		},
		PayloadFormat:    "sqs-event",
		MessageBatchSize: 10,
		WaitTime:         10,
	}

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "config section `lambda` must exist")
}

func TestValidateSqsLambdaTriggerConfigWithoutSqsConfigFails(t *testing.T) {
	cfg := &config.SqsLambdaTriggerConfig{
		LambdaConfig: &config.LambdaConfig{
			Endpoint:       "http://localhost:8080",
			Concurrency:    1,
			InvocationType: "RequestResponse",
		},
		PayloadFormat:    "sqs-event",
		MessageBatchSize: 10,
		WaitTime:         10,
	}

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "config section `sqs` must exist")
}

func TestValidateSqsLambdaTriggerConfigWithoutPayloadFormatFails(t *testing.T) {
	cfg := &config.SqsLambdaTriggerConfig{
		LambdaConfig: &config.LambdaConfig{
			Endpoint:       "http://localhost:8080",
			Concurrency:    1,
			InvocationType: "RequestResponse",
		},
		SqsConfig: &config.SqsConfig{
			Region:    "eu-central-1",
			Endpoint:  "http://localhost:4566",
			QueueName: "my-queue",
		},
		MessageBatchSize: 10,
		WaitTime:         10,
	}

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "config value `payload_format` must exist")
}

func TestValidateSqsLambdaTriggerConfigWithInvalidPayloadFormatFails(t *testing.T) {
	cfg := &config.SqsLambdaTriggerConfig{
		LambdaConfig: &config.LambdaConfig{
			Endpoint:       "http://localhost:8080",
			Concurrency:    1,
			InvocationType: "RequestResponse",
		},
		SqsConfig: &config.SqsConfig{
			Region:    "eu-central-1",
			Endpoint:  "http://localhost:4566",
			QueueName: "my-queue",
		},
		PayloadFormat:    "invalid",
		MessageBatchSize: 10,
		WaitTime:         10,
	}

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "config value `payload_format` must be one of [sqs-event, pure]")
}

func TestValidateSqsLambdaTriggerConfigWithInvalidMessageBatchSizeFails(t *testing.T) {
	cfg := &config.SqsLambdaTriggerConfig{
		LambdaConfig: &config.LambdaConfig{
			Endpoint:       "http://localhost:8080",
			Concurrency:    1,
			InvocationType: "RequestResponse",
		},
		SqsConfig: &config.SqsConfig{
			Region:    "eu-central-1",
			Endpoint:  "http://localhost:4566",
			QueueName: "my-queue",
		},
		PayloadFormat:    "sqs-event",
		MessageBatchSize: 0,
		WaitTime:         10,
	}

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "config value `message_batch_size` must be greater than 0")
}

func TestValidateSqsLambdaTriggerConfigWithInvalidWaitTimeFails(t *testing.T) {
	cfg := &config.SqsLambdaTriggerConfig{
		LambdaConfig: &config.LambdaConfig{
			Endpoint:       "http://localhost:8080",
			Concurrency:    1,
			InvocationType: "RequestResponse",
		},
		SqsConfig: &config.SqsConfig{
			Region:    "eu-central-1",
			Endpoint:  "http://localhost:4566",
			QueueName: "my-queue",
		},
		PayloadFormat:    "sqs-event",
		MessageBatchSize: 10,
		WaitTime:         0,
	}

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "config value `wait_time` must be greater than 0")
}

func TestValidateSqsLambdaTriggerConfigSucceeds(t *testing.T) {
	cfg := &config.SqsLambdaTriggerConfig{
		LambdaConfig: &config.LambdaConfig{
			Endpoint:       "http://localhost:8080",
			Concurrency:    1,
			InvocationType: "RequestResponse",
		},
		SqsConfig: &config.SqsConfig{
			Region:    "eu-central-1",
			Endpoint:  "http://localhost:4566",
			QueueName: "my-queue",
		},
		PayloadFormat:    "sqs-event",
		MessageBatchSize: 10,
		WaitTime:         10,
	}

	err := cfg.Validate()
	assert.NoError(t, err)
}
