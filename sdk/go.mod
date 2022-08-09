module github.com/xuanlingzi/go-admin-core/sdk

go 1.16

require (
	github.com/bsm/redislock v0.7.2
	github.com/bytedance/go-tagexpr/v2 v2.9.2
	github.com/casbin/casbin/v2 v2.47.1
	github.com/chanxuehong/wechat v0.0.0-20211009063332-41a5c6d8b38b
	github.com/gin-gonic/gin v1.7.7
	github.com/go-playground/locales v0.14.0
	github.com/go-playground/universal-translator v0.18.0
	github.com/go-playground/validator/v10 v10.11.0
	github.com/go-redis/redis/v8 v8.11.5
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/google/uuid v1.3.0
	github.com/gorilla/websocket v1.5.0
	github.com/mojocn/base64Captcha v1.3.5
	github.com/nsqio/go-nsq v1.1.0
	github.com/robfig/cron/v3 v3.0.1
	github.com/robinjoseph08/redisqueue/v2 v2.1.0
	github.com/shamsher31/goimgext v1.0.0
	github.com/slok/go-http-metrics v0.10.0
	github.com/smartystreets/goconvey v1.7.2
	github.com/xuanlingzi/go-admin-core v1.3.14
	github.com/xuanlingzi/go-admin-core/plugins/logger/zap v1.3.5-rc.1
	github.com/xuanlingzi/gorm-adapter/v3 v3.2.1-0.20220523055446-8b16ec988e0b
	golang.org/x/crypto v0.0.0-20220518034528-6f7dac969898
	gorm.io/gorm v1.23.5
)

require (
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/henrylee2cn/goutil v0.0.0-20210818094442-ed2b3cfe804b // indirect
	github.com/nyaruka/phonenumbers v1.0.75 // indirect
	github.com/smartystreets/assertions v1.13.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	go.uber.org/zap v1.21.0 // indirect
)

//replace github.com/xuanlingzi/go-admin-core => ../
