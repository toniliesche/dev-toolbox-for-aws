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
	"dev-toolbox-for-aws/domain/services"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"sync"
)

type SecretsManager struct {
	config      *config.SecretsManagerConfig
	mutex       sync.Mutex
	secretStore *services.SecretsManagerSecretsStore
}

func (m *SecretsManager) Run() error {
	if err := m.setup(); err != nil {
		return err
	}

	r := mux.NewRouter()
	r.HandleFunc("/", m.handleRequest).Methods(http.MethodPost)

	return http.ListenAndServe(":8080", r)
}

func (m *SecretsManager) setup() error {
	return m.secretStore.InitSecretsFromConfig(m.config)
}

func (m *SecretsManager) handleRequest(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		m.WriteError(writer, http.StatusMethodNotAllowed, "MethodNotAllowed", "Method not allowed")
		return
	}

	target := request.Header.Get("X-Amz-Target")
	switch target {
	case "secretsmanager.GetSecretValue":
		m.handleGetSecretValue(writer, request)
		break
	case "secretsmanager.ListSecrets":
		m.handleListSecrets(writer, request)
		break
	case "":
		m.WriteError(writer, http.StatusBadRequest, "MissingHeaderException", "X-Amz-Target header is required")
		break
	default:
		m.WriteError(writer, http.StatusBadRequest, "InvalidTargetTypeException", "Unknown target type")
		break
	}
}

func (m *SecretsManager) handleGetSecretValue(writer http.ResponseWriter, request *http.Request) {
	req := &model.SecretsManagerGetSecretRequest{}

	body, err := io.ReadAll(request.Body)
	if err != nil {
		m.WriteError(writer, http.StatusBadRequest, "InvalidRequestException", "Invalid JSON request body")
		return
	}

	err = json.Unmarshal(body, req)
	if err != nil {
		m.WriteError(writer, http.StatusBadRequest, "InvalidRequestException", "Invalid JSON request body")
		return
	}

	if req.SecretId == "" {
		m.WriteError(writer, http.StatusBadRequest, "InvalidRequestException", "Missing required parameter: SecretId")
		return
	}

	secret, secretVersion, err := m.secretStore.GetSecret(req.SecretId, req.VersionId, req.VersionStage)
	if err != nil {
		m.WriteError(writer, http.StatusBadRequest, "ResourceNotFoundException", "Secrets Manager can't find the specified secret.")
		return
	}

	response := model.SecretsManagerGetSecretResponse{
		ARN:          secret.ARN,
		CreatedDate:  secretVersion.CreatedDate,
		Name:         secret.Name,
		SecretString: secretVersion.SecretString,
		VersionId:    secretVersion.VersionId,
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		m.WriteError(writer, http.StatusInternalServerError, "InternalServiceErrorException", "An error occurred on the server side")
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(responseJson)
}

func (m *SecretsManager) handleListSecrets(writer http.ResponseWriter, request *http.Request) {
	req := &model.SecretsManagerListSecretsRequest{}

	body, err := io.ReadAll(request.Body)
	if err != nil {
		m.WriteError(writer, http.StatusBadRequest, "InvalidRequestException", "Invalid JSON request body")
		return
	}

	err = json.Unmarshal(body, req)
	if err != nil {
		m.WriteError(writer, http.StatusBadRequest, "InvalidRequestException", "Invalid JSON request body")
		return
	}

	if req.MaxResults == 0 {
		req.MaxResults = 10
	}

	secrets, nextToken, err := m.secretStore.ListSecrets(req.MaxResults, req.NextToken)
	if err != nil {
		m.WriteError(writer, http.StatusInternalServerError, "InternalServiceErrorException", "An error occurred on the server side")
		return
	}

	response := model.SecretsManagerListSecretsResponse{
		SecretList: secrets,
		NextToken:  nextToken,
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		m.WriteError(writer, http.StatusInternalServerError, "InternalServiceErrorException", "An error occurred on the server side")
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(responseJson)
}

func (m *SecretsManager) WriteError(writer http.ResponseWriter, httpStatusCode int, errorType string, errorMessage string) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(httpStatusCode)
	errorResponse := model.ErrorResponse{
		Type:    errorType,
		Message: errorMessage,
	}

	jsonError, _ := json.Marshal(errorResponse)
	writer.Write(jsonError)
}

func ProvideSecretsManager(cfg *config.SecretsManagerConfig) (*SecretsManager, error) {
	secretsStore, err := services.ProvideSecretsManagerSecretsStore()
	if err != nil {
		return nil, err
	}

	return &SecretsManager{
		secretStore: secretsStore,
		config:      cfg,
	}, nil
}
