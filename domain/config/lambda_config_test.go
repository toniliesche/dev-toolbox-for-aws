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

func TestValidateLambdaConfigFailsOnEmptyEndpoint(t *testing.T) {
	cfg := &config.LambdaConfig{
		Endpoint:       "",
		Concurrency:    1,
		InvocationType: "RequestResponse",
	}

	err := cfg.Validate("lambda")
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "config value `lambda.endpoint` must exist")
}

func TestValidateLambdaConfigFailsOnEmptyInvocationType(t *testing.T) {
	cfg := &config.LambdaConfig{
		Endpoint:       "http://localhost:8080",
		Concurrency:    1,
		InvocationType: "",
	}

	err := cfg.Validate("lambda")
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "config value `lambda.invocation_type` must exist")
}

func TestValidateLambdaConfigFailsOnInvalidInvocationType(t *testing.T) {
	cfg := &config.LambdaConfig{
		Endpoint:       "http://localhost:8080",
		Concurrency:    1,
		InvocationType: "Invalid",
	}

	err := cfg.Validate("lambda")
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "config value `lambda.invocation_type` must be one of [Event, RequestResponse]")
}

func TestValidateLambdaConfigFailsOnInvalidConcurrency(t *testing.T) {
	cfg := &config.LambdaConfig{
		Endpoint:       "http://localhost:8080",
		Concurrency:    0,
		InvocationType: "RequestResponse",
	}

	err := cfg.Validate("lambda")
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "config value `lambda.concurrency` must be greater than 0")
}

func TestValidateLambdaConfigSucceeds(t *testing.T) {
	cfg := &config.LambdaConfig{
		Endpoint:       "http://localhost:8080",
		Concurrency:    1,
		InvocationType: "RequestResponse",
	}

	err := cfg.Validate("lambda")
	assert.NoError(t, err)
}
