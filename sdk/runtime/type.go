package runtime

import (
	"github.com/gin-gonic/gin"
	"github.com/xuanlingzi/go-admin-core/block_chain"
	"github.com/xuanlingzi/go-admin-core/lbs"
	"github.com/xuanlingzi/go-admin-core/message"
	"github.com/xuanlingzi/go-admin-core/third_party"
	"net/http"

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
	GetFileStoreKey(key string) storage.AdapterFileStore

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
	GetLocationBasedServiceKey(key string) lbs.AdapterLocationBasedService

	SetBlockChainAdapter(string, block_chain.AdapterBroker)
	GetBlockChainAdapter() block_chain.AdapterBroker
	GetBlockChainAdapters() map[string]block_chain.AdapterBroker
	GetBlockChainKey(key string) block_chain.AdapterBroker

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
