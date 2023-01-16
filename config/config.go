package config

import (
	"os"

	env "github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type cfg struct {
	Env                               string `env:"ENV" envDefault:"test"`
	BasePath                          string `env:"BASE_PATH" envDefault:"."`
	Port                              int    `env:"PORT" envDefault:"80"`
	InternalPort                      int    `env:"INTERNAL_PORT" envDefault:"3000"`
	LogLevel                          string `env:"LOG_LEVEL" envDefault:"warn"`
	DmAccessKeyId                     string `env:"DM_ACCESS_KEY_ID"`
	DmAccessKeySecret                 string `env:"DM_ACCESS_KEY_SECRET"`
	UpdateLocaleAndKeyHighlightSecret string `env:"UPDATE_LOCALE_AND_KEYHIGHLIGHT_SECRET" envDefault:"f25a2fc72690b780b2a14e140ef6a9e0"`

	DBHostname string `env:"DB_HOSTNAME" envDefault:"localhost"`
	DBPort     int    `env:"DB_PORT" envDefault:"3306"`
	DBUsername string `env:"DB_USERNAME" envDefault:"root"`
	DBPassword string `env:"DB_PASSWORD" envDefault:"123456"`
	DBDatabase string `env:"DB_DATABASE" envDefault:"oauth2_production"`

	Scopes            []string `env:"SCOPES"`
	IconPrefix        string   `env:"ICON_PREFIX" envDefault:"https://cdn.shoplazza.com"`
	StoreService      string   `env:"SERVICE_URL_TOTORO"`
	ShanksService     string   `env:"SERVICE_URL_SHANKS"`
	StoreServiceKoala string   `env:"SERVICE_URL_KOALA"`
	PufferServiceHost string   `env:"SERVICE_URL_PUFFER" envDefault:"http://localhost:3000"`

	LlamaServiceHost string   `env:"SERVICE_URL_LLAMA" envDefault:"https://127.0.0.1:3001"`
	SentryDSN        string   `env:"SENTRY_DSN"`
	KafkaHost        []string `env:"KAFKA_HOST" envDefault:"[]string{'localhost:9092'}"`

	CookieName                     string `env:"COOKIE_NAME" envDefault:"awesomev2"`
	AppStoreHost                   string `env:"APP_STORE_HOST" envDefault:"https://appstore1024.shoplazza.com"`
	SSOCallbackPath                string `env:"SSO_CALLBACK_PATH" envDefault:"/api/auth/callback"`
	SSOClientId                    string `env:"SSO_CLIENT_ID" envDefault:"4bf2418f-a2f6-425a-a11d-327e5cc70291"`
	SSOClientSecret                string `env:"SSO_CLIENT_SECRET" envDefault:"033639c52f5748279f2be64282c045206377f9e2aeb14ba4bb3da10d1bd62983"`
	SSOHost                        string `env:"SSO_HOST" envDefault:"https://sso1024.shoplazza.com"`
	StateCookieKey                 string `env:"STATE_COOKIE_KEY" envDefault:"state"`
	FeebackLinkSecret              string `env:"FEEDBACK_LINK_SECRET" envDefault:"21371a76-7894-4b8d-935f-9d70c55d7d4a"`
	AppSupportFbLinkExpired        int    `env:"APP_SUPPORT_FB_LINK_EXPIRED" envDefault:"30"`
	SendFeedbackEmailInterval      int    `env:"SEND_FEEDBACK_EMAIL_INTERVAL" envDefault:"3"`
	AppReviewsInviteTaskStartDay   string `env:"EMAIL_GUIDES_MERCHANT_REVIEWS_START_DAY" envDefault:"2022-04-18"`
	AppReviewsInviteAfterInstalled int    `env:"EMAIL_GUIDES_MERCHANT_REVIEWS_INTERVAL"  envDefault:"168"`
	NotifyAppSupportFbAfterCreated int    `env:"NOTIFY_APP_SUPPORT_FB_AFTER_CREATED" envDefault:"3"`
	StoreCreatedAt                 string `env:"STORE_CREATED_AT" envDefault:"2022-06-14 17:00:00"`
}

var Cfg *cfg

func init() {
	if os.Getenv("ENV") == "locale" {
		// if no ENV specified

		// local debugging
		err := godotenv.Load()
		if err != nil {
			panic(err)
		}

		logrus.Info("dot env loaded")
	}

	Cfg = &cfg{}
	if err := env.Parse(Cfg); err != nil {
		panic(err)
	}
	logrus.Infof("config ===>>> %+v", Cfg)

	InitSSO()
}
