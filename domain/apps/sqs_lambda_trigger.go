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
	"bytes"
	"dev-toolbox-for-aws/domain/config"
	"dev-toolbox-for-aws/domain/model"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"io"
	"log"
	"net/http"
	"time"
)

type SqsLambdaTrigger struct {
	config     *config.SqsLambdaTriggerConfig
	queueUrl   string
	sqsService *sqs.SQS
	httpClient *http.Client
	semaphore  chan int
}

func (t *SqsLambdaTrigger) Run() error {
	if err := t.setup(); err != nil {
		return err
	}

	for {
		time.Sleep(5 * time.Second)

		t.semaphore <- 1
		msgResult, err := t.sqsService.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            &t.queueUrl,
			MaxNumberOfMessages: aws.Int64(t.config.MessageBatchSize),
			WaitTimeSeconds:     aws.Int64(t.config.WaitTime),
		})

		if err != nil {
			<-t.semaphore
			log.Printf("Error receiving message: %v", err)
			continue
		}

		if msgResult == nil || msgResult.Messages == nil || len(msgResult.Messages) == 0 {
			<-t.semaphore
			continue
		}

		t.handle(msgResult)
	}
}

func (t *SqsLambdaTrigger) setup() error {
	awsSession := session.Must(
		session.NewSession(
			&aws.Config{
				Region:      aws.String(t.config.SqsConfig.Region),
				Endpoint:    aws.String(t.config.SqsConfig.Endpoint),
				Credentials: credentials.NewStaticCredentials("x", "x", "x"),
			},
		),
	)

	t.sqsService = sqs.New(awsSession)
	result, err := t.sqsService.GetQueueUrl(
		&sqs.GetQueueUrlInput{
			QueueName: aws.String(t.config.SqsConfig.QueueName),
		},
	)

	if err != nil {
		return fmt.Errorf("error getting queue URL: %v", err)
	}

	t.queueUrl = *result.QueueUrl
	t.httpClient = &http.Client{
		Timeout: 30 * time.Minute,
	}

	t.semaphore = make(chan int, t.config.LambdaConfig.Concurrency)

	return nil
}

func (t *SqsLambdaTrigger) handle(result *sqs.ReceiveMessageOutput) {
	if t.config.PayloadFormat == "sqs-event" {
		t.handleSqsMessageBatch(result)
	} else {
		t.handleActionMessage(result)
	}
}

func (t *SqsLambdaTrigger) handleSqsMessageBatch(result *sqs.ReceiveMessageOutput) {
	go func(messageBatch []*sqs.Message) {
		defer func() {
			<-t.semaphore
		}()

		if len(messageBatch) == 0 {
			log.Println("No messages to process")

			return
		}

		log.Printf("Received batch of %d messages", len(messageBatch))

		payload, err := t.wrapSqsMessageBatch(messageBatch)

		log.Printf("Sending event: %s", payload)

		if err != nil {
			log.Printf("Error marshalling message: %v", err)
			return
		}

		var successfulMessages []*sqs.Message
		if successfulMessages, err = t.invokeLambdaWithSqsMessageBatch(payload, messageBatch); err != nil {
			log.Printf("Error invoking lambda: %v", err)
			return
		}

		for _, msg := range successfulMessages {
			log.Printf("Successfully processed message: %s", *msg.MessageId)
			_, err = t.sqsService.DeleteMessage(
				&sqs.DeleteMessageInput{
					QueueUrl:      &t.queueUrl,
					ReceiptHandle: msg.ReceiptHandle,
				},
			)
		}
	}(result.Messages)
}

func (t *SqsLambdaTrigger) handleActionMessage(result *sqs.ReceiveMessageOutput) {
	for _, msg := range result.Messages {
		go func(msg *sqs.Message) {
			defer func() {
				<-t.semaphore
			}()

			log.Printf("Received message: %s", *msg.Body)

			payload := []byte(*msg.Body)

			log.Printf("Sending event: %s", payload)

			var success bool
			var err error
			if success, err = t.invokeLambdaWithActionMessage(payload); err != nil {
				log.Printf("Error invoking lambda: %v", err)
				return
			}

			if success {
				log.Printf("Successfully invoked lambda")
				_, err = t.sqsService.DeleteMessage(
					&sqs.DeleteMessageInput{
						QueueUrl:      &t.queueUrl,
						ReceiptHandle: msg.ReceiptHandle,
					},
				)

			}
		}(msg)
	}
}

func (t *SqsLambdaTrigger) invokeLambdaWithSqsMessageBatch(payload []byte, messages []*sqs.Message) ([]*sqs.Message, error) {
	req, err := http.NewRequest("POST", t.config.LambdaConfig.Endpoint+"/2015-03-31/functions/function/invocations", bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("X-Amz-Invocation-Type", t.config.LambdaConfig.InvocationType)
	req.Header.Add("Content-Type", "application/json")

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error invoking lambda: %v", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("error invoking lambda: response is nil")
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error invoking lambda: %v", resp.Status)
	}

	body, _ := io.ReadAll(resp.Body)
	log.Printf("Lambda response: %s", string(body))

	lambdaResponse := &model.LambdaResponse{}

	err = json.Unmarshal(body, lambdaResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling lambda response: %v", err)
	}

	failedMessages := make(map[string]bool)
	for _, failure := range lambdaResponse.BatchItemFailures {
		failedMessages[failure.ItemIdentifier] = true
	}

	successfulMessages := make([]*sqs.Message, 0)
	for _, msg := range messages {
		if _, ok := failedMessages[*msg.MessageId]; !ok {
			successfulMessages = append(successfulMessages, msg)
		}
	}

	return successfulMessages, nil
}

func (t *SqsLambdaTrigger) invokeLambdaWithActionMessage(payload []byte) (bool, error) {
	req, err := http.NewRequest("POST", t.config.LambdaConfig.Endpoint+"/2015-03-31/functions/function/invocations", bytes.NewReader(payload))
	if err != nil {
		return false, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("X-Amz-Invocation-Type", t.config.LambdaConfig.InvocationType)
	req.Header.Add("Content-Type", "application/json")

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("error invoking lambda: %v", err)
	}

	if resp == nil {
		return false, fmt.Errorf("error invoking lambda: response is nil")
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("error invoking lambda: %v", resp.Status)
	}

	body, _ := io.ReadAll(resp.Body)
	log.Printf("Lambda response: %s", string(body))

	return true, nil
}

func (t *SqsLambdaTrigger) wrapSqsMessageBatch(batch []*sqs.Message) ([]byte, error) {
	records := make([]map[string]interface{}, len(batch))

	for _, msg := range batch {
		records = append(records, map[string]interface{}{
			"messageId":         msg.MessageId,
			"receiptHandle":     msg.ReceiptHandle,
			"body":              msg.Body,
			"attributes":        msg.Attributes,
			"messageAttributes": msg.MessageAttributes,
			"md5OfBody":         msg.MD5OfBody,
			"eventSource":       "aws:sqs",
			"eventSourceARN":    t.queueUrl,
			"awsRegion":         t.config.SqsConfig.Region,
		})
	}

	wrappedMsg := map[string][]map[string]interface{}{
		"Records": records,
	}

	return json.Marshal(wrappedMsg)
}

func ProvideSqsLambdaTrigger(cfg *config.SqsLambdaTriggerConfig) (*SqsLambdaTrigger, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("failed validation application config: %v", err)
	}

	return &SqsLambdaTrigger{
		config: cfg,
	}, nil
}
