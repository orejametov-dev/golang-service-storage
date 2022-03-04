package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"os"
	"time"
)

const (
	readTimeout        time.Duration = 10 * time.Second
	writeTimeout       time.Duration = 10 * time.Second
	maxHeaderMegabytes int           = 1
)

type Config struct {
	APP    APPConfig
	DB     DBConfig
	HTTP   HTTPConfig
	MinIO  MinIOConfig
}

type MinIOConfig struct {
	Endpoint string
	AccessKey string
	SecretKey string
}

type APPConfig struct {
	Env string
	Name string
	Debug bool
}

type DBConfig struct {
	Connection string
	Host       string
	Port       string
	Database   string
	Username   string
	Password   string
}

type HTTPConfig struct {
	Port               string
	Host               string
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	MaxHeaderMegabytes int
}

func LoadConfig(path string) (config *Config, err error) {
	err = parseConfigFile(path)
	if err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	if err = setFromEnv(config); err != nil {
		return nil, err
	}
	return config, nil
}

func parseConfigFile(path string) error {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return viper.MergeInConfig()
}


func setFromEnv(config *Config) error {

	err := godotenv.Load()
	if err != nil {
		return err
	}

	config.APP.Env = os.Getenv("APP_ENV")
	if os.Getenv("APP_DEBUG") == "true" {
		config.APP.Debug = true
	}
	config.APP.Name = os.Getenv("APP_NAME")

	config.HTTP.Port = os.Getenv("DEFAULT_HTTP_PORT")
	config.HTTP.Host = os.Getenv("DEFAULT_HTTP_HOST")
	config.HTTP.ReadTimeout = readTimeout
	config.HTTP.WriteTimeout = writeTimeout
	config.HTTP.MaxHeaderMegabytes = maxHeaderMegabytes

	config.DB.Connection = os.Getenv("DB_CONNECTION")
	config.DB.Port = os.Getenv("DB_PORT")
	config.DB.Host = os.Getenv("DB_HOST")
	config.DB.Database = os.Getenv("DB_DATABASE")
	config.DB.Username = os.Getenv("DB_USERNAME")
	config.DB.Password = os.Getenv("DB_PASSWORD")

	config.MinIO.AccessKey = os.Getenv("MINIO_ACCESS_KEY")
	config.MinIO.Endpoint = os.Getenv("MINIO_ENDPOINT")
	config.MinIO.SecretKey = os.Getenv("MINIO_SECRET_KEY")

	return nil
}

// func getEnv(key string, defaultVal string) string {
// 	if value, exists := os.LookupEnv(key); exists {
// 		return value
// 	}
// 	return defaultVal
// }

