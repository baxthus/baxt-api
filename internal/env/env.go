package env

import "os"

var (
	Port          string
	AvatarURL     string
	LoggerWebhook string
)

func init() {
	Port = os.Getenv("PORT")
	AvatarURL = os.Getenv("AVATAR_URL")
	LoggerWebhook = os.Getenv("LOGGER_WEBHOOK")
}
