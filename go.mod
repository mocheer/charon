module github.com/mocheer/charon

go 1.16

require (
	github.com/form3tech-oss/jwt-go v3.2.3+incompatible
	github.com/gofiber/fiber/v2 v2.12.0
	github.com/gofiber/helmet/v2 v2.1.4
	github.com/gofiber/jwt/v2 v2.2.2
	github.com/gofiber/websocket/v2 v2.0.4
	github.com/mocheer/pluto v1.0.0
	github.com/mocheer/vesta v0.0.0-00010101000000-000000000000
	github.com/sirupsen/logrus v1.8.1
	github.com/tdewolff/canvas v0.0.0-20210602122216-44a8ab7d5e7f
	github.com/tidwall/gjson v1.8.0
	github.com/twpayne/go-geom v1.4.0
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a
	gorm.io/datatypes v1.0.1
	gorm.io/driver/postgres v1.1.0
	gorm.io/gorm v1.21.10
)

replace (
	github.com/mocheer/pluto => ../pluto
	github.com/mocheer/vesta => ../vesta
)
