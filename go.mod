module github.com/rivernews/k8s-cluster-head-service/v2

go 1.14

// add line below to apply go version in heroku deployment
// https://github.com/heroku/heroku-buildpack-go/issues/301#issuecomment-471032174
// +heroku goVersion go1.14

require (
	github.com/braintree/manners v0.0.0-20160418043613-82a8879fc5fd // indirect
	github.com/gin-gonic/gin v1.6.3
	github.com/go-playground/validator/v10 v10.3.0 // indirect
	github.com/gocraft/web v0.0.0-20190207150652-9707327fb69b // indirect
	github.com/gocraft/work v0.5.1
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/gomodule/redigo v1.8.2
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/robfig/cron v1.2.0 // indirect
	github.com/stretchr/testify v1.6.1
	golang.org/x/sys v0.0.0-20200620081246-981b61492c35 // indirect
	google.golang.org/protobuf v1.24.0
	gopkg.in/yaml.v2 v2.3.0 // indirect
)
