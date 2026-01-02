package ssmutil

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

var (
	clientOnce sync.Once
	client     *ssm.Client
	clientErr  error
)

func getClient() (*ssm.Client, error) {
	clientOnce.Do(func() {
		cfg, clientErr := config.LoadDefaultConfig(context.Background())
		if clientErr != nil {
			return
		}
		client = ssm.NewFromConfig(cfg)
	})
	return client, clientErr
}

// GetParameter fetches a SecureString/String parameter with decryption.
func GetParameter(name string) (string, error) {
	c, err := getClient()
	if err != nil {
		return "", fmt.Errorf("load SSM client: %w", err)
	}

	withDecryption := true
	out, err := c.GetParameter(context.Background(), &ssm.GetParameterInput{
		Name:           &name,
		WithDecryption: &withDecryption,
	})
	if err != nil {
		return "", fmt.Errorf("get parameter %s: %w", name, err)
	}

	if out.Parameter == nil || out.Parameter.Value == nil {
		return "", fmt.Errorf("parameter %s has no value", name)
	}

	return *out.Parameter.Value, nil
}
