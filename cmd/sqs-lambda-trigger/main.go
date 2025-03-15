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

package main

import (
	"dev-toolbox-for-aws/domain/apps"
	"dev-toolbox-for-aws/domain/config"
	"log"
)

func main() {
	sqsLambdaTriggerConfig, err := config.ProvideSqsLambdaTriggerConfig()
	if err != nil {
		log.Fatalf("Error providing config: %v", err)
	}

	sqsLambdaTrigger, err := apps.ProvideSqsLambdaTrigger(sqsLambdaTriggerConfig)
	if err != nil {
		log.Fatalf("Error providing sqs lambda trigger: %v", err)
	}

	if err = sqsLambdaTrigger.Run(); err != nil {
		log.Fatalf("Error running sqs lambda trigger: %v", err)
	}
}
