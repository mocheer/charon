module github.com/mocheer/charon

go 1.16

require (
	github.com/form3tech-oss/jwt-go v3.2.3+incompatible
	github.com/gofiber/fiber/v2 v2.13.0
	github.com/gofiber/helmet/v2 v2.1.6
	github.com/gofiber/jwt/v2 v2.2.3
	github.com/gofiber/websocket/v2 v2.0.6
	github.com/mocheer/pluto v1.0.0
	github.com/mocheer/vesta v0.0.0-00010101000000-000000000000
	github.com/sirupsen/logrus v1.8.1
	github.com/tdewolff/canvas v0.0.0-20210624142932-07559ac5d4f8
	github.com/tidwall/gjson v1.8.0
	github.com/twpayne/go-geom v1.4.0
	golang.org/x/crypto v0.0.0-20210616213533-5ff15b29337e
	gorm.io/datatypes v1.0.1
	gorm.io/driver/postgres v1.1.0
	gorm.io/gorm v1.21.11
)

replace (
	github.com/mocheer/pluto => ../pluto
	github.com/mocheer/vesta => ../vesta
)
