package config

type ApplicationConfig struct {
	ServerAddress string
	RedisUrl string
	DatabaseUrl string	
}
func NewApplicationConfig(addr string,redis_url string, databaseUrl string) *ApplicationConfig{
	return &ApplicationConfig{
		ServerAddress: addr,
		RedisUrl: redis_url,
		DatabaseUrl: databaseUrl,
	}
}