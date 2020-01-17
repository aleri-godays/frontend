package config

//Config holds all configuration for this service
type Config struct {
	ServiceName       string
	Version           string
	LogLevel          string
	HTTPPort          int
	OAuthClientID     string
	OAuthClientSecret string
	JWTSecret         string
	ProjectURI        string
	TimeTrackingURI   string
}
