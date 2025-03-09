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

type SQSLambdaTrigger struct {
	config     *config.SQSLambdaTriggerConfig
	queueUrl   string
	sqsService *sqs.SQS
	httpClient *http.Client
}

func (t *SQSLambdaTrigger) Run() error {
	if err := t.setup(); err != nil {
		return err
	}

	for {
		msgResult, err := t.sqsService.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            &t.queueUrl,
			MaxNumberOfMessages: aws.Int64(t.config.SQSConfig.MessageBatchSize),
			WaitTimeSeconds:     aws.Int64(t.config.SQSConfig.WaitTime),
		})

		if err != nil {
			log.Printf("Error receiving message: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		t.handle(msgResult)
	}
}

func (t *SQSLambdaTrigger) setup() error {
	awsSession := session.Must(
		session.NewSession(
			&aws.Config{
				Region:      aws.String(t.config.SQSConfig.Region),
				Endpoint:    aws.String(t.config.SQSConfig.Endpoint),
				Credentials: credentials.NewStaticCredentials("x", "x", "x"),
			},
		),
	)

	t.sqsService = sqs.New(awsSession)
	result, err := t.sqsService.GetQueueUrl(
		&sqs.GetQueueUrlInput{
			QueueName: aws.String(t.config.SQSConfig.QueueName),
		},
	)

	if err != nil {
		return fmt.Errorf("error getting queue URL: %v", err)
	}

	t.queueUrl = *result.QueueUrl
	t.httpClient = &http.Client{
		Timeout: 30 * time.Minute,
	}

	return nil
}

func (t *SQSLambdaTrigger) handle(result *sqs.ReceiveMessageOutput) {
	var payload []byte
	var err error

	semaphore := make(chan int, t.config.LambdaConfig.Concurrency)

	for _, msg := range result.Messages {
		semaphore <- 1

		go func(msg *sqs.Message) {
			defer func() {
				<-semaphore
			}()

			log.Printf("Received message: %s", *msg.Body)

			if t.config.PayloadFormat == "sqs-event" {
				payload, err = t.wrapSQSEvent(msg)
			} else {
				payload, err = []byte(*msg.Body), nil
			}

			log.Printf("Sending event: %s", payload)

			if err != nil {
				log.Printf("Error marshalling message: %v", err)
				return
			}

			var success bool
			if success, err = t.invokeLambda(payload); err != nil {
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

func (t *SQSLambdaTrigger) invokeLambda(payload []byte) (bool, error) {
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

func (t *SQSLambdaTrigger) wrapSQSEvent(msg *sqs.Message) ([]byte, error) {
	return json.Marshal(
		map[string]interface{}{
			"Records": []map[string]interface{}{
				{
					"messageId":         *msg.MessageId,
					"receiptHandle":     *msg.ReceiptHandle,
					"body":              *msg.Body,
					"attributes":        msg.Attributes,
					"messageAttributes": msg.MessageAttributes,
					"eventSource":       "aws:sqs",
					"eventSourceARN":    t.queueUrl,
					"awsRegion":         t.config.SQSConfig.Region,
				},
			},
		},
	)
}

func ProvideSQSLambdaTrigger(cfg *config.SQSLambdaTriggerConfig) (*SQSLambdaTrigger, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("failed validation application config: %v", err)
	}

	return &SQSLambdaTrigger{
		config: cfg,
	}, nil
}
