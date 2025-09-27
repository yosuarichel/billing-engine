package config

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretspb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/yosuarichel/billing-engine/pkg/utils"
	"google.golang.org/api/option"
	"gopkg.in/yaml.v3"
)

// AppConfig holds both config + secrets
type AppConfig struct {
	AppName  string      `json:"app_name" yaml:"app_name"`
	Env      string      `json:"env" yaml:"env"`
	DB       DBConfig    `json:"db" yaml:"db"`
	Redis    RedisConfig `json:"redis" yaml:"redis"`
	HTTPPort int         `json:"http_port" yaml:"http_port"`
	RPCPort  int         `json:"rpc_port" yaml:"rpc_port"`
	NSQ      NSQConfig   `json:"nsq" yaml:"nsq"`
	// DBDSN         string   `json:"db_dsn" yaml:"db_dsn"`
	// RedisAddr     string   `json:"redis_addr" yaml:"redis_addr"`
	// RedisPassword string   `json:"redis_password" yaml:"redis_password"`
	// NSQAddr       string   `json:"nsq_addr" yaml:"nsq_addr"`
	// EtcdEndpoints []string `json:"etcd_endpoints" yaml:"etcd_endpoints"`

}

type DBConfig struct {
	Driver       string        `yaml:"driver"`
	PSM          string        `yaml:"psm"`
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	User         string        `yaml:"user"`
	Password     string        `yaml:"password"`
	DBName       string        `yaml:"db_name"`
	MaxIdleConns int           `yaml:"max_idle_conns"`
	MaxOpenConns int           `yaml:"max_open_conns"`
	Timeout      time.Duration `yaml:"timeout"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	DBCharset    string        `yaml:"db_charset"`
	SSLMode      string        `yaml:"sslmode"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type NSQConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

var cfg = &AppConfig{}

func LoadConfig(ctx context.Context) (*AppConfig, error) {
	env := utils.GetEnvWithKey(utils.APP_ENV, "dev")
	projectID := utils.GetEnvWithKey("GCP_PROJECT_ID", "")
	secretID := utils.GetEnvWithKey("GCP_SECRET_ID", "")
	credFile := utils.GetEnvWithKey("GOOGLE_APPLICATION_CREDENTIALS", "")

	// Try secret manager
	secretmanager, err := loadFromGSM(ctx, projectID, secretID, credFile)
	if err == nil {
		return secretmanager, nil
	}
	// Log / notify fallback
	klog.Infof("⚠️ Warning: using fallback local config due to secret manager error: %v\n", err)

	// Fallback
	localPath := fmt.Sprintf("./conf/app_config_%s.yaml", env)
	data, err := os.ReadFile(localPath)
	if err != nil {
		return nil, fmt.Errorf("read local yaml: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("unmarshal local yaml: %w", err)
	}
	return cfg, nil
}

// loadFromGSM fetches the JSON payload stored in the secret and unmarshals it
func loadFromGSM(ctx context.Context, projectID, secretID, credFilePath string) (*AppConfig, error) {
	// Build the secret version name: “latest”
	name := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectID, secretID)

	// Create client with optional credentials file
	var client *secretmanager.Client
	var err error
	if credFilePath != "" {
		client, err = secretmanager.NewClient(ctx, option.WithCredentialsFile(credFilePath))
	} else {
		client, err = secretmanager.NewClient(ctx)
	}
	if err != nil {
		return nil, fmt.Errorf("secretmanager.NewClient: %w", err)
	}
	defer client.Close()

	// Access secret version
	req := &secretspb.AccessSecretVersionRequest{
		Name: name,
	}
	resp, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("AccessSecretVersion: %w", err)
	}
	if resp.Payload == nil || resp.Payload.Data == nil {
		return nil, errors.New("secretmanager: empty payload")
	}

	// Unmarshal secret JSON into AppConfig
	if err := json.Unmarshal(resp.Payload.Data, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal secret payload: %w", err)
	}

	return cfg, nil
}

func GetAppCfg() *AppConfig {
	return cfg
}
