package runtime

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xuanlingzi/go-admin-core/lbs"
	"github.com/xuanlingzi/go-admin-core/message"
	"github.com/xuanlingzi/go-admin-core/moderation"
	"github.com/xuanlingzi/go-admin-core/payment"
	"github.com/xuanlingzi/go-admin-core/third_party"

	"github.com/casbin/casbin/v2"
	"github.com/robfig/cron/v3"
	"github.com/xuanlingzi/go-admin-core/logger"
	"github.com/xuanlingzi/go-admin-core/storage"
	"gorm.io/gorm"
)

type Runtime interface {
	// SetDb 多db设置，⚠️SetDbs不允许并发,可以根据自己的业务，例如app分库、host分库
	SetDb(key string, db *gorm.DB)
	GetDb() map[string]*gorm.DB
	GetDbByKey(key string) *gorm.DB

	SetCasbin(key string, enforcer *casbin.SyncedEnforcer)
	GetCasbin() map[string]*casbin.SyncedEnforcer
	GetCasbinKey(key string) *casbin.SyncedEnforcer

	// SetEngine 使用的路由
	SetEngine(engine http.Handler)
	GetEngine() http.Handler

	GetRouter() []Router

	// SetLogger 使用go-admin定义的logger，参考来源go-micro
	SetLogger(logger logger.Logger)
	GetLogger() logger.Logger

	// SetCrontab crontab
	SetCrontab(key string, crontab *cron.Cron)
	GetCrontab() map[string]*cron.Cron
	GetCrontabKey(key string) *cron.Cron

	// SetMiddleware middleware
	SetMiddleware(string, interface{})
	GetMiddleware() map[string]interface{}
	GetMiddlewareKey(key string) interface{}

	// SetCacheAdapter cache
	SetCacheAdapter(storage.AdapterCache)
	GetCacheAdapter() storage.AdapterCache
	GetCachePrefix(string) storage.AdapterCache

	GetMemoryQueue(string) storage.AdapterQueue
	SetQueueAdapter(storage.AdapterQueue)
	GetQueueAdapter() storage.AdapterQueue
	GetQueuePrefix(string) storage.AdapterQueue

	SetLockerAdapter(storage.AdapterLocker)
	GetLockerAdapter() storage.AdapterLocker
	GetLockerPrefix(string) storage.AdapterLocker

	SetSmsAdapter(string, message.AdapterSms)
	GetSmsAdapter() message.AdapterSms
	GetSmsAdapters() map[string]message.AdapterSms
	GetSmsKey(key string) message.AdapterSms

	SetMailAdapter(string, message.AdapterMail)
	GetMailAdapter() message.AdapterMail
	GetMailAdapters() map[string]message.AdapterMail
	GetMailKey(key string) message.AdapterMail

	SetFileStoreAdapter(string, storage.AdapterFileStore)
	GetFileStoreAdapter() storage.AdapterFileStore
	GetFileStoreAdapters() map[string]storage.AdapterFileStore
	GetFileStoreKey(string) storage.AdapterFileStore

	SetModerationAdapter(string, moderation.AdapterModeration)
	GetModerationAdapter() moderation.AdapterModeration
	GetModerationAdapters() map[string]moderation.AdapterModeration
	GetModerationKey(string) moderation.AdapterModeration

	SetAmqpAdapter(string, message.AdapterAmqp)
	GetAmqpAdapter() message.AdapterAmqp
	GetAmqpAdapters() map[string]message.AdapterAmqp
	GetAmqpKey(key string) message.AdapterAmqp

	SetThirdPartyAdapter(string, third_party.AdapterThirdParty)
	GetThirdPartyAdapter() third_party.AdapterThirdParty
	GetThirdPartyAdapters() map[string]third_party.AdapterThirdParty
	GetThirdPartyKey(key string) third_party.AdapterThirdParty

	SetLocationBasedServiceAdapter(string, lbs.AdapterLocationBasedService)
	GetLocationBasedServiceAdapter() lbs.AdapterLocationBasedService
	GetLocationBasedServiceAdapters() map[string]lbs.AdapterLocationBasedService
	GetLocationBasedServiceKey(string) lbs.AdapterLocationBasedService

	SetPaymentServiceAdapter(string, payment.AdapterPaymentService)
	GetPaymentServiceAdapter() payment.AdapterPaymentService
	GetPaymentServiceAdapters() map[string]payment.AdapterPaymentService
	GetPaymentServiceKey(string) payment.AdapterPaymentService

	SetHandler(key string, routerGroup func(r *gin.RouterGroup, hand ...*gin.HandlerFunc))
	GetHandler() map[string][]func(r *gin.RouterGroup, hand ...*gin.HandlerFunc)
	GetHandlerPrefix(key string) []func(r *gin.RouterGroup, hand ...*gin.HandlerFunc)

	GetStreamMessage(id, stream string, value map[string]interface{}) (storage.Messager, error)

	GetConfig(key string) interface{}
	SetConfig(key string, value interface{})

	// SetAppRouters set AppRouter
	SetAppRouters(appRouters func())
	GetAppRouters() []func()
}
