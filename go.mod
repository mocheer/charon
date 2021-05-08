module github.com/mocheer/charon

go 1.16

require (
	github.com/ajstarks/svgo v0.0.0-20210406150507-75cfd577ce75 // indirect
	github.com/andybalholm/brotli v1.0.2 // indirect
	github.com/form3tech-oss/jwt-go v3.2.2+incompatible
	github.com/gofiber/fiber/v2 v2.9.0
	github.com/gofiber/helmet/v2 v2.1.2 // indirect
	github.com/gofiber/jwt/v2 v2.2.1
	github.com/gofrs/uuid v4.0.0+incompatible // indirect
	github.com/jackc/pgmock v0.0.0-20201204152224-4fe30f7445fd // indirect
	github.com/jackc/pgproto3/v2 v2.0.7 // indirect
	github.com/klauspost/compress v1.12.2 // indirect
	github.com/lib/pq v1.10.1 // indirect
	github.com/mattn/go-sqlite3 v1.14.7 // indirect
	github.com/mocheer/pluto v0.0.0-20210429095047-c4e6c07bef14
	github.com/shopspring/decimal v1.2.0 // indirect
	github.com/stretchr/testify v1.7.0 // indirect
	github.com/tdewolff/canvas v0.0.0-20210505183520-925a1f0167be // direct
	github.com/twpayne/go-geom v1.4.0
	github.com/valyala/fasthttp v1.24.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/crypto v0.0.0-20210506145944-38f3c27a63bf
	golang.org/x/image v0.0.0-20210504121937-7319ad40d33e // indirect
	golang.org/x/sys v0.0.0-20210507161434-a76c4d0a0096 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	gorm.io/datatypes v1.0.1
	gorm.io/driver/postgres v1.1.0
	gorm.io/gorm v1.21.9
)

replace github.com/mocheer/pluto => ../pluto
