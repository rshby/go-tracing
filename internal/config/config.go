package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

func GetEnv(key string) string {
	if err := godotenv.Load("./env"); err != nil {
		return ""
	}

	return os.Getenv(key)
}

func AppPort() int {
	if port := GetEnv("APP_PORT"); port != "" {
		appPort, err := strconv.Atoi(port)
		if err != nil {
			return DefaultAppPort
		}

		return appPort
	}

	return DefaultAppPort
}

func DatabaseHost() string {
	if host := GetEnv("DB_HOST"); host != "" {
		return host
	}

	return DefaultDatabaseHost
}

func DatabaseUser() string {
	if user := GetEnv("DB_USER"); user != "" {
		return user
	}

	return DefaultDatabaseUser
}

func DatabasePassword() string {
	if password := GetEnv("DB_PASSWORD"); password != "" {
		return password
	}

	return DefaultDatabasePassword
}

func DatabasePort() int {
	if port := GetEnv("DB_PORT"); port != "" {
		dbPort, err := strconv.Atoi(port)
		if err != nil {
			return DefaultDatabasePort
		}

		return dbPort
	}

	return DefaultDatabasePort
}

func DatabaseName() string {
	if name := GetEnv("DB_NAME"); name != "" {
		return name
	}

	return DefaultDatabaseName
}

func DatabaseDSN() string {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=UTC",
		DatabaseUser(), DatabasePassword(), DatabaseHost(), DatabasePort(), DatabaseName())

	return dsn
}

func OtelExporter() string {
	return GetEnv("OTEL_EXPORTER")
}
