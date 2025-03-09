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
	"strings"
)

type SNSConfig struct {
	Region string      `yaml:"region"`
	Topics []*SNSTopic `yaml:"topics"`
}

type SNSTopic struct {
	Name          string             `yaml:"name"`
	Subscriptions []*SNSSubscription `yaml:"subscriptions"`
}

type SNSSubscription struct {
	Identifier string `yaml:"identifier"`
	Type       string `yaml:"type"`
}

func (c *SNSConfig) Validate(path string) error {
	if c.Region == "" {
		return fmt.Errorf("config value `%s.region` must exist", path)
	}

	if c.Topics == nil {
		return fmt.Errorf("config value `%s.topics` must exist", path)
	}

	if len(c.Topics) == 0 {
		return fmt.Errorf("config value `%s.topics` must contain at least one element", path)
	}

	for i, t := range c.Topics {
		if err := t.Validate(fmt.Sprintf("%s.topics[%d]", path, i)); err != nil {
			return err
		}
	}

	return nil
}

func (t *SNSTopic) Validate(path string) error {
	if t.Name == "" {
		return fmt.Errorf("config value `%s.name` must exist", path)
	}

	if t.Subscriptions == nil {
		return fmt.Errorf("config value `%s.subscriptions` must exist", path)
	}

	if len(t.Subscriptions) == 0 {
		return fmt.Errorf("config value `%s.subscriptions` must contain at least one element", path)
	}

	for i, s := range t.Subscriptions {
		if err := s.Validate(fmt.Sprintf("%s.subscriptions[%d]", path, i)); err != nil {
			return err
		}
	}

	return nil
}

func (s *SNSSubscription) Validate(path string) error {
	if s.Identifier == "" {
		return fmt.Errorf("config value `%s.identifier` must exist", path)
	}

	if s.Type == "" {
		return fmt.Errorf("config value `%s.type` must exist", path)
	}

	if s.Type != "sqs" {
		return fmt.Errorf("config value `%s.type` must be `sqs`", path)
	}

	return nil
}

func getSNSConfigFromEnvironmentVariables() (*SNSConfig, error) {
	topics := shared.GetEnvironmentString("SNS_TOPICS", "")
	if topics == "" {
		return nil, fmt.Errorf("SNS_TOPICS environment variable must exist")
	}

	sqsSubscriptions := shared.GetEnvironmentString("SNS_SQS_SUBSCRIPTIONS", "")
	if sqsSubscriptions == "" {
		return nil, fmt.Errorf("SNS_SQS_SUBSCRIPTIONS environment variable must exist")
	}

	topicIdentifiers := strings.Split(topics, ",")
	snsTopics := make([]*SNSTopic, 0, len(topicIdentifiers))
	for _, topicIdentifier := range topicIdentifiers {
		newTopic := SNSTopic{
			Name:          topicIdentifier,
			Subscriptions: make([]*SNSSubscription, 0),
		}

		snsTopics = append(snsTopics, &newTopic)
	}

	sqsSubscriptionIdentifiers := strings.Split(sqsSubscriptions, ",")
	for _, sqsSubscriptionIdentifier := range sqsSubscriptionIdentifiers {
		if sqsSubscriptionIdentifier == "" {
			continue
		}

		sqsSubscriptionDetails := strings.Split(sqsSubscriptionIdentifier, ":")
		if len(sqsSubscriptionDetails) != 2 {
			continue
		}

		found := false
		for _, snsTopic := range snsTopics {
			if snsTopic.Name == sqsSubscriptionDetails[0] {
				newSubscription := SNSSubscription{
					Identifier: sqsSubscriptionDetails[1],
					Type:       "sqs",
				}
				snsTopic.Subscriptions = append(snsTopic.Subscriptions, &newSubscription)

				found = true
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("SNS topic `%s` not found", sqsSubscriptionDetails[0])
		}
	}

	return &SNSConfig{
		Region: shared.GetEnvironmentString("AWS_REGION", "eu-central-1"),
		Topics: snsTopics,
	}, nil
}
