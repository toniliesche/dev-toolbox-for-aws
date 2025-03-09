package services

import (
	"dev-toolbox-for-aws/domain/config"
	"dev-toolbox-for-aws/domain/model"
	"fmt"
	"github.com/google/uuid"
	"sync"
)

type SecretsManagerSecretsStore struct {
	store map[string]*model.SecretsManagerSecret
	mutex *sync.RWMutex
}

func (s *SecretsManagerSecretsStore) GetSecret(secret string, version string, stage string) (*model.SecretsManagerSecret, *model.SecretsManagerSecretVersion, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	secretObj, ok := s.store[secret]
	if !ok {
		return nil, nil, fmt.Errorf("secret not found")
	}

	if version != "" {
		for _, secretVersion := range secretObj.Versions {
			if secretVersion.VersionId == version {
				return secretObj, secretVersion, nil
			}
		}

		return nil, nil, fmt.Errorf("secretVersion not found")
	}

	if stage == "" {
		stage = "AWSCURRENT"
	}

	for _, secretVersion := range secretObj.Versions {
		if secretVersion.VersionStage == stage {
			return secretObj, secretVersion, nil
		}
	}

	return nil, nil, fmt.Errorf("stage not found")
}

func (s *SecretsManagerSecretsStore) InitSecretsFromConfig(cfg *config.SecretsManagerConfig) error {
	for _, secret := range cfg.Secrets {
		arn := fmt.Sprintf("arn:aws:secretsmanager:%s:%s:secret:%s", cfg.Region, cfg.Account, secret.Name)
		secretVersions := make([]*model.SecretsManagerSecretVersion, 0)

		for index, version := range secret.Versions {
			var stage string
			if index == 0 {
				stage = "AWSCURRENT"
			} else if index == 1 {
				stage = "AWSPREVIOUS"
			}

			secretVersion := &model.SecretsManagerSecretVersion{
				SecretString: version,
				VersionStage: stage,
				CreatedDate:  0,
				VersionId:    fmt.Sprintf("%s-%s", secret.Name, uuid.New().String()),
			}

			secretVersions = append(secretVersions, secretVersion)
		}

		s.store[secret.Name] = &model.SecretsManagerSecret{
			ARN:      arn,
			Name:     secret.Name,
			Versions: secretVersions,
		}
	}

	return nil
}

func (s *SecretsManagerSecretsStore) ListSecrets(maxResults int, startToken string) ([]*model.SecretsManagerSecretSimple, string, error) {
	s.mutex.RLock()
	secrets := make([]*model.SecretsManagerSecret, 0, len(s.store))
	for _, secret := range s.store {
		secrets = append(secrets, secret)
	}
	s.mutex.RUnlock()

	startIndex := 0
	if startToken != "" {
		fmt.Sscanf(startToken, "%d", &startIndex)
	}

	endIndex := startIndex + maxResults
	if endIndex > len(secrets) {
		endIndex = len(secrets)
	}

	nextToken := ""
	if endIndex < len(secrets) {
		nextToken = fmt.Sprintf("%d", endIndex)
	}

	simpleSecrets := make([]*model.SecretsManagerSecretSimple, 0, maxResults)
	for _, secret := range secrets[startIndex:endIndex] {
		simpleSecrets = append(simpleSecrets, &model.SecretsManagerSecretSimple{
			ARN:         secret.ARN,
			Name:        secret.Name,
			CreatedDate: secret.Versions[0].CreatedDate,
		})
	}

	return simpleSecrets, nextToken, nil
}

func ProvideSecretsManagerSecretsStore() (*SecretsManagerSecretsStore, error) {
	return &SecretsManagerSecretsStore{
		store: make(map[string]*model.SecretsManagerSecret),
		mutex: &sync.RWMutex{},
	}, nil
}
