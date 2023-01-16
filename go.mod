module gitlab.shoplazza.site/shoplaza/cobra

require (
	bou.ke/monkey v1.0.2
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/Masterminds/squirrel v1.5.3
	github.com/Pallinder/go-randomdata v1.2.0
	github.com/Shopify/sarama v1.27.2
	github.com/caarlos0/env/v6 v6.6.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/getsentry/raven-go v0.2.0
	github.com/gin-contrib/sentry v0.0.0-20191119142041-ff0e9556d1b7
	github.com/gin-gonic/gin v1.8.1
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gofrs/uuid v4.3.1+incompatible // indirect
	github.com/golang/mock v1.6.0
	github.com/google/uuid v1.3.0
	github.com/gosimple/slug v1.12.0
	github.com/jarcoal/httpmock v1.2.0
	github.com/jinzhu/gorm v1.9.2
	github.com/joho/godotenv v1.4.0
	github.com/json-iterator/go v1.1.12
	github.com/keepeye/logrus-filename v0.0.0-20190711075016-ce01a4391dd1
	github.com/kikyousky/gormigrate v1.3.1
	github.com/pkg/errors v0.9.1
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/shopspring/decimal v1.3.1
	github.com/sirupsen/logrus v1.9.0
	github.com/smartystreets/goconvey v1.7.2
	github.com/spf13/cast v1.5.0
	github.com/stretchr/testify v1.8.0
	github.com/thinkeridea/go-extend v1.3.2
	github.com/urfave/cli v1.22.1
	gopkg.in/h2non/gock.v1 v1.1.2
)

go 1.14

replace (
	golang.org/x/net => golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4
	golang.org/x/text v0.3.7 => golang.org/x/text v0.3.8
)
