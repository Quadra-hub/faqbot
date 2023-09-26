package config

import "os"

func GptApiKey() string {
	return os.Getenv("GPT_API_KEY")
}

func DBConnString() string {
	return os.Getenv("DATABASE_URL")
}
