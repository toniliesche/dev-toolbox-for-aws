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

package apps

import (
	"dev-toolbox-for-aws/domain/config"
	"dev-toolbox-for-aws/domain/model"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
	"strings"
)

type SNSPublisher struct {
	config         *config.SNSPublisherConfig
	sqsService     *sqs.SQS
	queueUrls      map[string]string
	sqsSubscribers map[string][]string
}

func (p *SNSPublisher) Run() error {
	if err := p.setup(); err != nil {
		return err
	}

	r := mux.NewRouter()
	r.HandleFunc("/Publish", p.handleSNSRequests).Methods(http.MethodPost)

	return http.ListenAndServe(":8080", r)
}

func (p *SNSPublisher) handleSNSRequests(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := request.ParseForm(); err != nil {
		http.Error(writer, fmt.Sprintf("Error parsing query parameters: %v", err), http.StatusBadRequest)
		return
	}

	message := request.FormValue("Message")
	topicArn := request.FormValue("TopicArn")
	messageStructure := request.FormValue("MessageStructure")

	if message == "" || topicArn == "" {
		http.Error(writer, "Message and TopicArn are required", http.StatusBadRequest)
		return
	}

	messageAttributes := p.parseMessageAttributes(request.Form)
	var err error

	topicParts := strings.Split(topicArn, ":")
	if len(topicParts) != 6 {
		http.Error(writer, "Invalid TopicArn", http.StatusBadRequest)
		return
	}

	topic := topicParts[5]
	if _, ok := p.sqsSubscribers[topic]; ok {
		if err = p.sendToSQSSubscribers(topic, message, messageStructure, messageAttributes); err != nil {
			http.Error(writer, fmt.Sprintf("Error sending message to SQS subscribers: %v", err), http.StatusInternalServerError)
			return
		}
	}

	messageID := uuid.New().String()
	resp := &model.SNSPublishResponse{
		MessageId: messageID,
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(resp)
}

func (p *SNSPublisher) setup() error {
	awsSession := session.Must(session.NewSession(
		&aws.Config{
			Region:      aws.String(p.config.SQSConfig.Region),
			Endpoint:    aws.String(p.config.SQSConfig.Endpoint),
			Credentials: credentials.NewStaticCredentials("x", "x", "x"),
		}))

	p.sqsService = sqs.New(awsSession)

	for _, topic := range p.config.SNSConfig.Topics {
		for _, subscription := range topic.Subscriptions {
			if _, ok := p.queueUrls[subscription.Identifier]; ok {
				continue
			}

			result, err := p.sqsService.GetQueueUrl(&sqs.GetQueueUrlInput{
				QueueName: aws.String(subscription.Identifier),
			})

			if err != nil {
				return fmt.Errorf("error getting queue URL: %v", err)
			}

			p.queueUrls[subscription.Identifier] = *result.QueueUrl
		}
	}

	for _, topic := range p.config.SNSConfig.Topics {
		p.sqsSubscribers[topic.Name] = make([]string, 0)

		for _, subscription := range topic.Subscriptions {
			if subscription.Type == "sqs" {
				p.sqsSubscribers[topic.Name] = append(p.sqsSubscribers[topic.Name], subscription.Identifier)
			}
		}
	}

	return nil
}

func (p *SNSPublisher) parseMessageAttributes(form url.Values) map[string]map[string]string {
	attributes := make(map[string]map[string]string)

	for key, value := range form {
		if len(value) == 0 {
			continue
		}

		if len(key) > 24 && key[:24] == "MessageAttributes.entry." {
			parts := p.splitAttributeKey(key)

			if parts.EntryID == "" || parts.Field == "" {
				continue
			}

			if _, exists := attributes[parts.EntryID]; !exists {
				attributes[parts.EntryID] = make(map[string]string)
			}

			attributes[parts.EntryID][parts.Field] = value[0]
		}
	}

	return attributes
}

func (p *SNSPublisher) splitAttributeKey(key string) model.AttributeKeyParts {
	var parts model.AttributeKeyParts

	attributeKey := key[24:]
	for i := 0; i < len(attributeKey); i++ {
		if attributeKey[i] == '.' {
			parts.EntryID = attributeKey[:i]
			parts.Field = attributeKey[i+1:]
			break
		}
	}

	return parts
}

func (p *SNSPublisher) sendToSQSQueue(queueUrl string, message string, attributes map[string]*sqs.MessageAttributeValue) error {
	_, err := p.sqsService.SendMessage(
		&sqs.SendMessageInput{
			MessageBody:       aws.String(message),
			QueueUrl:          aws.String(queueUrl),
			MessageAttributes: attributes,
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func (p *SNSPublisher) mapSqsMessageAttributes(attributes map[string]map[string]string) map[string]*sqs.MessageAttributeValue {
	sqsAttributes := make(map[string]*sqs.MessageAttributeValue)
	for _, attr := range attributes {
		sqsAttributes[attr["Name"]] = &sqs.MessageAttributeValue{
			DataType:    aws.String(attr["DataType"]),
			StringValue: aws.String(attr["StringValue"]),
		}
	}

	return sqsAttributes
}

func (p *SNSPublisher) extractMessage(structure string, subscriberType string, message string) (string, error) {
	if structure == "" {
		return message, nil
	}

	if structure != "json" {
		return "", fmt.Errorf("unsupported message structure: %s", structure)
	}

	var messageMap map[string]string

	err := json.Unmarshal([]byte(message), &messageMap)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling message: %v", err)
	}

	if msg, exists := messageMap[subscriberType]; exists {
		return msg, nil
	}

	if msg, exists := messageMap["default"]; exists {
		return msg, nil
	}

	return "", fmt.Errorf("message not found for subscriber type: %s", subscriberType)
}

func (p *SNSPublisher) sendToSQSSubscribers(topic string, message string, structure string, attributes map[string]map[string]string) error {
	sqsAttributes := p.mapSqsMessageAttributes(attributes)

	sqsMessage, err := p.extractMessage(structure, "sqs", message)
	if err != nil {
		return fmt.Errorf("error extracting message structure: %v", err)
	}

	var queueUrl string
	var exists bool
	for _, queue := range p.sqsSubscribers[topic] {
		if queueUrl, exists = p.queueUrls[queue]; !exists {
			return fmt.Errorf("queue %s not found", queue)
		}

		if err = p.sendToSQSQueue(queueUrl, sqsMessage, sqsAttributes); err != nil {
			return err
		}
	}

	return nil
}

func ProvideSNSPublisher(cfg *config.SNSPublisherConfig) (*SNSPublisher, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("failed validation application config: %v", err)
	}

	return &SNSPublisher{
		config:         cfg,
		queueUrls:      make(map[string]string),
		sqsSubscribers: make(map[string][]string),
	}, nil
}
