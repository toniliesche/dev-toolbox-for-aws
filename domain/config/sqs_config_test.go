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

func TestValidateSqsConfigFailsOnEmptyRegion(t *testing.T) {
	cfg := &config.SqsConfig{
		Region:    "",
		Endpoint:  "http://localhost:4566",
		QueueName: "my-queue",
	}

	err := cfg.Validate("sqs", false)
	assert.Error(t, err)
	assert.Equal(t, "config value `sqs.region` must exist", err.Error())
}

func TestValidateSqsConfigFailsOnEmptyQueueName(t *testing.T) {
	cfg := &config.SqsConfig{
		Region:    "eu-central-1",
		Endpoint:  "http://localhost:4566",
		QueueName: "",
	}

	err := cfg.Validate("sqs", false)
	assert.Error(t, err)
	assert.Equal(t, "config value `sqs.queue_name` must exist", err.Error())
}

func TestValidateSqsConfigFailsOnEmptyEndpoint(t *testing.T) {
	cfg := &config.SqsConfig{
		Region:    "eu-central-1",
		Endpoint:  "",
		QueueName: "my-queue",
	}

	err := cfg.Validate("sqs", false)
	assert.Error(t, err)
	assert.Equal(t, "config value `sqs.endpoint` must exist", err.Error())
}

func TestValidateSqsConfigSucceeds(t *testing.T) {
	cfg := &config.SqsConfig{
		Region:    "eu-central-1",
		Endpoint:  "http://localhost:4566",
		QueueName: "my-queue",
	}

	err := cfg.Validate("sqs", false)
	assert.NoError(t, err)
}
func TestValidateSqsConfigSucceedsWithAllowEmptyQueue(t *testing.T) {
	cfg := &config.SqsConfig{
		Region:    "eu-central-1",
		Endpoint:  "http://localhost:4566",
		QueueName: "",
	}

	err := cfg.Validate("sqs", true)
	assert.NoError(t, err)
}
