package configs

import (
	"github.com/go-chi/jwtauth"
	"github.com/melkzsiqueira/water-gas-measurement/docs"
	"github.com/spf13/viper"
)

type conf struct {
	DBDriver         string `mapstructure:"DB_DRIVER"`
	DBHost           string `mapstructure:"DB_HOST"`
	DBPort           string `mapstructure:"DB_PORT"`
	DBUser           string `mapstructure:"DB_USER"`
	DBPassword       string `mapstructure:"DB_PASSWORD"`
	DBName           string `mapstructure:"DB_NAME"`
	DBSSLMode        string `mapstructure:"DB_SSL_MODE"`
	DBTimezone       string `mapstructure:"DB_TIMEZONE"`
	WebServerPort    string `mapstructure:"WEB_SERVER_PORT"`
	WebServerHost    string `mapstructure:"WEB_SERVER_HOST"`
	JWTSecret        string `mapstructure:"JWT_SECRET"`
	JWTExpiresIn     int    `mapstructure:"JWT_EXPIRES_IN"`
	APIVersion       string `mapstructure:"API_VERSION"`
	GeminiKey        string `mapstructure:"GEMINI_API_KEY"`
	GeminiModel      string `mapstructure:"GEMINI_MODEL"`
	StorageAPIKey    string `mapstructure:"STORAGE_API_KEY"`
	StorageAPISecret string `mapstructure:"STORAGE_API_SECRET"`
	StorageName      string `mapstructure:"STORAGE_NAME"`
	DBDSN            string
	SwaggerURL       string
	TokenAuth        *jwtauth.JWTAuth
}

func LoadConfig(path string) (*conf, error) {
	var cfg *conf

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
	cfg.DBDSN = "host=" + cfg.DBHost + " port=" + cfg.DBPort + " user=" + cfg.DBUser + " dbname=" + cfg.DBName + " password=" + cfg.DBPassword + " sslmode=" + cfg.DBSSLMode + " TimeZone=" + cfg.DBTimezone
	cfg.SwaggerURL = "http://" + cfg.WebServerHost + ":" + cfg.WebServerPort + "/" + cfg.APIVersion + "/docs/doc.json"

	docs.SwaggerInfo.Host = cfg.WebServerHost + ":" + cfg.WebServerPort
	docs.SwaggerInfo.BasePath = "/" + cfg.APIVersion

	return cfg, err
}
