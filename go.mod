module github.com/mocheer/charon

go 1.16

require (
	github.com/ajstarks/svgo v0.0.0-20210406150507-75cfd577ce75 // indirect
	github.com/form3tech-oss/jwt-go v3.2.2+incompatible
	github.com/gofiber/fiber/v2 v2.8.0
	github.com/gofiber/jwt/v2 v2.2.1
	github.com/gofrs/uuid v4.0.0+incompatible // indirect
	github.com/jackc/pgmock v0.0.0-20201204152224-4fe30f7445fd // indirect
	github.com/jackc/pgproto3/v2 v2.0.7 // indirect
	github.com/jackc/pgx/v4 v4.11.0 // indirect
	github.com/klauspost/compress v1.12.1 // indirect
	github.com/lib/pq v1.10.0 // indirect
	github.com/mattn/go-sqlite3 v1.14.6 // indirect
	github.com/mocheer/pluto v0.0.0-20210419103954-a857fd1137fc
	github.com/shopspring/decimal v1.2.0 // indirect
	github.com/stretchr/testify v1.7.0 // indirect
	github.com/tdewolff/canvas v0.0.0-20210420033948-0f660268c5ad // direct
	github.com/tdewolff/minify/v2 v2.9.16 // indirect
	github.com/tdewolff/parse/v2 v2.5.15 // indirect
	github.com/twpayne/go-geom v1.3.6
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/crypto v0.0.0-20210415154028-4f45737414dc
	golang.org/x/image v0.0.0-20210220032944-ac19c3e999fb // indirect
	golang.org/x/text v0.3.6 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	gorm.io/datatypes v1.0.1
	gorm.io/driver/postgres v1.0.8
	gorm.io/gorm v1.21.8
)

replace github.com/mocheer/pluto => ../pluto
