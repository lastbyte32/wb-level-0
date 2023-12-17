package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type App struct {
	Mode     string `env:"MODE" env-default:"local"`
	DataBase DataBaseCfg
	Nats     NatsCfg
}

type DataBaseCfg struct {
	URL string `env:"DATABASE_URL" env-default:"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"`
}
type NatsCfg struct {
	URL       string `env:"NATS_URL" env-default:"nats://localhost:4222"`
	ClusterID string `env:"NATS_CLUSTER_ID" env-default:"test-cluster"`
	ClientID  string `env:"NATS_CLIENT_ID" env-default:"test-subscriber"`
	Subject   string `env:"NATS_SUBJECT" env-default:"test-subject"`
}

func New() (*App, error) {
	var instance App
	err := cleanenv.ReadEnv(&instance)
	if err != nil {
		return nil, fmt.Errorf("failed to read environment variables: %w", err)
	}
	return &instance, err
}
