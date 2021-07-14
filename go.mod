module github.com/mocheer/charon

go 1.16

require (
	github.com/form3tech-oss/jwt-go v3.2.3+incompatible
	github.com/gofiber/fiber/v2 v2.14.0
	github.com/gofiber/helmet/v2 v2.1.7
	github.com/gofiber/jwt/v2 v2.2.4
	github.com/gofiber/websocket/v2 v2.0.7
	github.com/mocheer/pluto v1.0.0
	github.com/mocheer/vesta v0.0.0-00010101000000-000000000000
	github.com/mocheer/xena v0.0.0-00010101000000-000000000000
	github.com/rubenv/topojson v0.0.0-20180822134236-13be738db397
	github.com/sirupsen/logrus v1.8.1
	github.com/tdewolff/canvas v0.0.0-20210629004453-6d68f6ba1416
	github.com/tidwall/gjson v1.8.1
	github.com/twpayne/go-geom v1.4.0
	golang.org/x/crypto v0.0.0-20210616213533-5ff15b29337e
	gorm.io/datatypes v1.0.1
	gorm.io/driver/postgres v1.1.0
	gorm.io/gorm v1.21.11
)

replace (
	github.com/mocheer/pluto => ../pluto
	github.com/mocheer/vesta => ../vesta
	github.com/mocheer/xena => ../xena
)
