module github.com/mocheer/charon

go 1.16

require (
	github.com/form3tech-oss/jwt-go v3.2.3+incompatible
	github.com/gofiber/fiber/v2 v2.10.0
	github.com/gofiber/helmet/v2 v2.1.3
	github.com/gofiber/jwt/v2 v2.2.2
	github.com/gofiber/websocket/v2 v2.0.4
	github.com/mocheer/pluto v1.0.0
	github.com/mocheer/vesta v0.0.0-20210603235947-af0492b5793a
	github.com/sirupsen/logrus v1.8.1
	github.com/tdewolff/canvas v0.0.0-20210531202529-1e808269d181
	github.com/tidwall/gjson v1.8.0
	github.com/twpayne/go-geom v1.4.0
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a
	golang.org/x/sys v0.0.0-20210603125802-9665404d3644 // indirect
	gorm.io/datatypes v1.0.1
	gorm.io/driver/postgres v1.1.0
	gorm.io/gorm v1.21.10
)

replace (
	github.com/mocheer/pluto => ../pluto
	github.com/mocheer/vesta => ../vesta
)
