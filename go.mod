module github.com/mocheer/charon

go 1.16

require (
	github.com/adrg/strutil v0.2.3 // indirect
	github.com/andybalholm/brotli v1.0.3 // indirect
	github.com/benoitkugler/textlayout v0.0.1 // indirect
	github.com/fasthttp/websocket v1.4.3 // indirect
	github.com/form3tech-oss/jwt-go v3.2.3+incompatible
	github.com/gofiber/fiber/v2 v2.12.0
	github.com/gofiber/helmet/v2 v2.1.4
	github.com/gofiber/jwt/v2 v2.2.2
	github.com/gofiber/websocket/v2 v2.0.4
	github.com/jackc/pgproto3/v2 v2.0.7 // indirect
	github.com/klauspost/compress v1.13.0 // indirect
	github.com/mocheer/pluto v1.0.0
	github.com/mocheer/vesta v0.0.0-20210604224936-355751c684de
	github.com/savsgio/gotils v0.0.0-20210520110740-c57c45b83e0a // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/tdewolff/canvas v0.0.0-20210602122216-44a8ab7d5e7f
	github.com/tidwall/gjson v1.8.0
	github.com/tidwall/pretty v1.1.1 // indirect
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
