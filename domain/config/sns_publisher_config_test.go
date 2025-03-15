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

func TestValidateSNSPublisherConfigWithoutSNSConfigFails(t *testing.T) {
	cfg := &config.SNSPublisherConfig{
		SqsConfig: &config.SqsConfig{
			Region:    "eu-central-1",
			Endpoint:  "http://localhost:4566",
			QueueName: "my-queue",
		},
	}

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Equal(t, "config section `sns` must exist", err.Error())
}

func TestValidateSNSPublisherConfigWithoutSqsConfigFails(t *testing.T) {
	cfg := &config.SNSPublisherConfig{
		SNSConfig: &config.SNSConfig{
			Region: "eu-central-1",
			Topics: []*config.SNSTopic{
				{
					Name: "topic",
					Subscriptions: []*config.SNSSubscription{
						{
							Identifier: "identifier",
							Type:       "sqs",
						},
					},
				},
			},
		},
	}

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Equal(t, "config section `sqs` must exist", err.Error())
}

func TestValidateSNSPublisherConfigSucceeds(t *testing.T) {
	cfg := &config.SNSPublisherConfig{
		SNSConfig: &config.SNSConfig{
			Region: "eu-central-1",
			Topics: []*config.SNSTopic{
				{
					Name: "topic",
					Subscriptions: []*config.SNSSubscription{
						{
							Identifier: "identifier",
							Type:       "sqs",
						},
					},
				},
			},
		},
		SqsConfig: &config.SqsConfig{
			Region:    "eu-central-1",
			Endpoint:  "http://localhost:4566",
			QueueName: "my-queue",
		},
	}

	err := cfg.Validate()
	assert.NoError(t, err)
}
