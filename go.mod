module github.com/linmadan/egglib-go

go 1.14

require (
	github.com/Shopify/sarama v1.23.1
	github.com/astaxie/beego v1.12.0
	github.com/beego/beego/v2 v2.0.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-pg/pg/v10 v10.7.7
	github.com/sirupsen/logrus v1.4.2
	gopkg.in/jcmturner/goidentity.v3 v3.0.0 // indirect
)

replace github.com/linmadan/egglib-go => ../egglib-go
