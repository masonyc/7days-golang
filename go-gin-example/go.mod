module github.com/masonyc/7days-golang/go-gin-example

go 1.14

require (
	github.com/astaxie/beego v1.12.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.6.3
	github.com/go-ini/ini v1.55.0
	github.com/golang/protobuf v1.4.1 // indirect
	github.com/jinzhu/gorm v1.9.12
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/unknwon/com v1.0.1
	golang.org/x/sys v0.0.0-20200509044756-6aff5f38e54f // indirect
)

replace (
	github.com/masonyc/7days-golang/go-gin-example/conf => ./conf
	github.com/masonyc/7days-golang/go-gin-example/db => ./db
	github.com/masonyc/7days-golang/go-gin-example/middleware => ./middleware
	github.com/masonyc/7days-golang/go-gin-example/middleware/jwt => ./middleware/jwt
	github.com/masonyc/7days-golang/go-gin-example/models => ./models
	github.com/masonyc/7days-golang/go-gin-example/pkg/e => ./pkg/e
	github.com/masonyc/7days-golang/go-gin-example/pkg/setting => ./pkg/setting
	github.com/masonyc/7days-golang/go-gin-example/pkg/util => ./pkg/util
	github.com/masonyc/7days-golang/go-gin-example/routers => ./routers
)
