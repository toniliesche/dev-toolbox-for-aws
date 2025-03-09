package model

type SecretsManagerListSecretsRequest struct {
	MaxResults int    `json:"MaxResults,omitempty"`
	NextToken  string `json:"NextToken,omitempty"`
}
