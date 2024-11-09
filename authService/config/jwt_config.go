package config

func GetJWTConfig() string {
	return GetConfig().JWTSecretKey
}