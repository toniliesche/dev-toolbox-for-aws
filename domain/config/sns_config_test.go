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

func TestValidateSNSConfigFailsOnEmptyRegion(t *testing.T) {
	cfg := &config.SNSConfig{
		Region: "",
		Topics: []*config.SNSTopic{
			{
				Name: "test-topic",
				Subscriptions: []*config.SNSSubscription{
					{
						Identifier: "test-subscription",
						Type:       "sqs",
					},
				},
			},
		},
	}

	err := cfg.Validate("sns")
	assert.Error(t, err)
	assert.Equal(t, "config value `sns.region` must exist", err.Error())
}

func TestValidateSNSConfigFailsOnMissingTopicList(t *testing.T) {
	cfg := &config.SNSConfig{
		Region: "eu-central-1",
	}

	err := cfg.Validate("sns")
	assert.Error(t, err)
	assert.Equal(t, "config value `sns.topics` must exist", err.Error())
}

func TestValidateSNSConfigFailsOnEmptyTopicList(t *testing.T) {
	cfg := &config.SNSConfig{
		Region: "eu-central-1",
		Topics: []*config.SNSTopic{},
	}

	err := cfg.Validate("sns")
	assert.Error(t, err)
	assert.Equal(t, "config value `sns.topics` must contain at least one element", err.Error())
}

func TestValidateSNSConfigSucceeds(t *testing.T) {
	cfg := &config.SNSConfig{
		Region: "eu-central-1",
		Topics: []*config.SNSTopic{
			{
				Name: "test-topic",
				Subscriptions: []*config.SNSSubscription{
					{
						Identifier: "test-subscription",
						Type:       "sqs",
					},
				},
			},
		},
	}

	err := cfg.Validate("sns")
	assert.NoError(t, err)
}

func TestValidateSNSTopicFailsOnEmptyName(t *testing.T) {
	cfg := &config.SNSTopic{
		Name: "",
		Subscriptions: []*config.SNSSubscription{
			{
				Identifier: "test-subscription",
				Type:       "sqs",
			},
		},
	}

	err := cfg.Validate("sns.topics[0]")
	assert.Error(t, err)
	assert.Equal(t, "config value `sns.topics[0].name` must exist", err.Error())
}

func testValidateSNSTopicFailsOnMissingSubscriptionList(t *testing.T) {
	cfg := &config.SNSTopic{
		Name: "test-topic",
	}

	err := cfg.Validate("sns.topics[0]")
	assert.Error(t, err)
	assert.Equal(t, "config value `sns.topics[0].subscriptions` must exist", err.Error())
}

func TestValidateSNSTopicFailsOnEmptySubscriptionList(t *testing.T) {
	cfg := &config.SNSTopic{
		Name:          "test-topic",
		Subscriptions: []*config.SNSSubscription{},
	}

	err := cfg.Validate("sns.topics[0]")
	assert.Error(t, err)
	assert.Equal(t, "config value `sns.topics[0].subscriptions` must contain at least one element", err.Error())
}

func TestValidateSNSTopicSucceeds(t *testing.T) {
	cfg := &config.SNSTopic{
		Name: "test-topic",
		Subscriptions: []*config.SNSSubscription{
			{
				Identifier: "test-subscription",
				Type:       "sqs",
			},
		},
	}

	err := cfg.Validate("sns.topics[0]")
	assert.NoError(t, err)
}

func TestValidateSNSSubscriptionFailsOnEmptyType(t *testing.T) {
	cfg := &config.SNSSubscription{
		Identifier: "test-subscription",
		Type:       "",
	}

	err := cfg.Validate("sns.topics[0].subscriptions[0]")
	assert.Error(t, err)
	assert.Equal(t, "config value `sns.topics[0].subscriptions[0].type` must exist", err.Error())
}

func TestValidateSNSSubscriptionFailsOnInvalidType(t *testing.T) {
	cfg := &config.SNSSubscription{
		Identifier: "test-subscription",
		Type:       "invalid",
	}

	err := cfg.Validate("sns.topics[0].subscriptions[0]")
	assert.Error(t, err)
	assert.Equal(t, "config value `sns.topics[0].subscriptions[0].type` must be `sqs`", err.Error())
}

func TestValidateSNSSubscriptionFailsOnEmptyIdentifier(t *testing.T) {
	cfg := &config.SNSSubscription{
		Identifier: "",
		Type:       "sqs",
	}

	err := cfg.Validate("sns.topics[0].subscriptions[0]")
	assert.Error(t, err)
	assert.Equal(t, "config value `sns.topics[0].subscriptions[0].identifier` must exist", err.Error())
}

func TestValidateSNSSubscriptionSucceeds(t *testing.T) {
	cfg := &config.SNSSubscription{
		Identifier: "test-subscription",
		Type:       "sqs",
	}

	err := cfg.Validate("sns.topics[0].subscriptions[0]")
	assert.NoError(t, err)
}
