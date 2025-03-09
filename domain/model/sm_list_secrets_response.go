package model

type SecretsManagerListSecretsResponse struct {
	SecretList []*SecretsManagerSecretSimple `json:"SecretList"`
	NextToken  string                        `json:"NextToken,omitempty"`
}
