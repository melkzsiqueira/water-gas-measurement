package configs

import (
	"github.com/go-chi/jwtauth"
	"github.com/spf13/viper"
)

var cfg *conf

type conf struct {
	DBDriver      string `mapsstructure:"DB_DRIVER"`
	DBHost        string `mapsstructure:"DB_HOST"`
	DBPort        string `mapsstructure:"DB_PORT"`
	DBUser        string `mapsstructure:"DB_USER"`
	DBPassword    string `mapsstructure:"DB_PASSWORD"`
	DBDatabase    string `mapsstructure:"DB_DATABASE"`
	WebServerPort string `mapsstructure:"WEB_SERVER_PORT"`
	JWTSecret     string `mapsstructure:"JWT_SECRET"`
	JWTExpiresIn  int    `mapsstructure:"JWT_EXPIRES_IN"`
	TokenAuth     *jwtauth.JWTAuth
}

func LoadConfig(path string) (*conf, error) {
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()

	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&cfg)

	if err != nil {
		panic(err)
	}

	cfg.TokenAuth = jwtauth.New("HS256", []byte(cfg.JWTSecret), nil)

	return cfg, err
}
