package runtime

import (
	"net/http"
	"sync"

	"github.com/xuanlingzi/go-admin-core/core/pkg/utils"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/xuanlingzi/go-admin-core/lbs"
	"github.com/xuanlingzi/go-admin-core/message"
	"github.com/xuanlingzi/go-admin-core/moderation"
	"github.com/xuanlingzi/go-admin-core/payment"
	"github.com/xuanlingzi/go-admin-core/rtc"
	"github.com/xuanlingzi/go-admin-core/third_party"

	"github.com/xuanlingzi/go-admin-core/logger"
	"github.com/xuanlingzi/go-admin-core/storage"
	"github.com/xuanlingzi/go-admin-core/storage/queue"
	"gorm.io/gorm"
)

type Application struct {
	dbs         map[string]*gorm.DB
	casbins     map[string]*casbin.SyncedEnforcer
	engine      http.Handler
	crontab     map[string]*cron.Cron
	mux         sync.RWMutex
	middlewares map[string]interface{}
	cache       storage.AdapterCache
	queues      map[string]storage.AdapterQueue
	locker      storage.AdapterLocker
	memoryQueue storage.AdapterQueue
	fileStores  map[string]storage.AdapterFileStore
	moderation  map[string]moderation.AdapterModeration
	sms         map[string]message.AdapterSms
	mail        map[string]message.AdapterMail
	amqp        map[string]message.AdapterAmqp
	mqtt        map[string]message.AdapterMqtt
	thirdParty  map[string]third_party.AdapterThirdParty
	lbs         map[string]lbs.AdapterLocationBasedService
	payment     map[string]payment.AdapterPaymentService
	leshua      map[string]payment.AdapterLeshuaService
	zego        map[string]rtc.AdapterZegoService
	handler     map[string][]func(r *gin.RouterGroup, hand ...*gin.HandlerFunc)
	routers     []Router
	configs     map[string]interface{} // 系统参数
	appRouters  []func()               // app路由
}

type Router struct {
	HttpMethod, RelativePath, Handler string
}

type Routers struct {
	List []Router
}

// SetDb 设置对应key的db
func (e *Application) SetDb(key string, db *gorm.DB) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.dbs[key] = db
}

// GetDb 获取所有map里的db数据
func (e *Application) GetDb() map[string]*gorm.DB {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.dbs
}

// GetDbByKey 根据key获取db
func (e *Application) GetDbByKey(key string) *gorm.DB {
	e.mux.RLock()
	defer e.mux.RUnlock()
	if db, ok := e.dbs["*"]; ok {
		return db
	}
	return e.dbs[key]
}

func (e *Application) SetCasbin(key string, enforcer *casbin.SyncedEnforcer) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.casbins[key] = enforcer
}

func (e *Application) GetCasbin() map[string]*casbin.SyncedEnforcer {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.casbins
}

// GetCasbinKey 根据key获取casbin
func (e *Application) GetCasbinKey(key string) *casbin.SyncedEnforcer {
	e.mux.RLock()
	defer e.mux.RUnlock()
	if e, ok := e.casbins["*"]; ok {
		return e
	}
	return e.casbins[key]
}

// SetEngine 设置路由引擎
func (e *Application) SetEngine(engine http.Handler) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.engine = engine
}

// GetEngine 获取路由引擎
func (e *Application) GetEngine() http.Handler {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.engine
}

// GetRouter 获取路由表
func (e *Application) GetRouter() []Router {
	return e.setRouter()
}

// setRouter 设置路由表
func (e *Application) setRouter() []Router {
	e.mux.Lock()
	defer e.mux.Unlock()

	e.routers = e.routers[:0]
	engine, ok := e.engine.(*gin.Engine)
	if !ok {
		return e.routers
	}

	routers := engine.Routes()
	for _, router := range routers {
		e.routers = append(e.routers, Router{RelativePath: router.Path, Handler: router.Handler, HttpMethod: router.Method})
	}
	return e.routers
}

// SetLogger 设置日志组件
func (e *Application) SetLogger(l logger.Logger) {
	logger.DefaultLogger = l
}

// GetLogger 获取日志组件
func (e *Application) GetLogger() logger.Logger {
	return logger.DefaultLogger
}

// NewConfig 默认值
func NewConfig() *Application {
	return &Application{
		dbs:         make(map[string]*gorm.DB),
		casbins:     make(map[string]*casbin.SyncedEnforcer),
		crontab:     make(map[string]*cron.Cron),
		middlewares: make(map[string]interface{}),
		queues:      make(map[string]storage.AdapterQueue),
		memoryQueue: queue.NewMemory(10000),
		fileStores:  make(map[string]storage.AdapterFileStore),
		moderation:  make(map[string]moderation.AdapterModeration),
		sms:         make(map[string]message.AdapterSms),
		mail:        make(map[string]message.AdapterMail),
		amqp:        make(map[string]message.AdapterAmqp),
		mqtt:        make(map[string]message.AdapterMqtt),
		thirdParty:  make(map[string]third_party.AdapterThirdParty),
		lbs:         make(map[string]lbs.AdapterLocationBasedService),
		payment:     make(map[string]payment.AdapterPaymentService),
		leshua:      make(map[string]payment.AdapterLeshuaService),
		zego:        make(map[string]rtc.AdapterZegoService),
		handler:     make(map[string][]func(r *gin.RouterGroup, hand ...*gin.HandlerFunc)),
		routers:     make([]Router, 0),
		configs:     make(map[string]interface{}),
	}
}

// SetCrontab 设置对应key的crontab
func (e *Application) SetCrontab(key string, crontab *cron.Cron) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.crontab[key] = crontab
}

// GetCrontab 获取所有map里的crontab数据
func (e *Application) GetCrontab() map[string]*cron.Cron {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.crontab
}

// GetCrontabKey 根据key获取crontab
func (e *Application) GetCrontabKey(key string) *cron.Cron {
	e.mux.RLock()
	defer e.mux.RUnlock()
	if e, ok := e.crontab["*"]; ok {
		return e
	}
	return e.crontab[key]
}

// SetMiddleware 设置中间件
func (e *Application) SetMiddleware(key string, middleware interface{}) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.middlewares[key] = middleware
}

// GetMiddleware 获取所有中间件
func (e *Application) GetMiddleware() map[string]interface{} {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.middlewares
}

// GetMiddlewareKey 获取对应key的中间件
func (e *Application) GetMiddlewareKey(key string) interface{} {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.middlewares[key]
}

// SetCacheAdapter 设置缓存
func (e *Application) SetCacheAdapter(c storage.AdapterCache) {
	e.cache = c
}

// GetCacheAdapter 获取缓存
func (e *Application) GetCacheAdapter() storage.AdapterCache {
	return NewCache("", e.cache)
}

// GetCachePrefix 获取带租户标记的cache
func (e *Application) GetCachePrefix(key string) storage.AdapterCache {
	return NewCache(key, e.cache)
}

// SetQueueAdapter 设置队列适配器
func (e *Application) SetQueueAdapter(key string, c storage.AdapterQueue) {
	e.mux.Lock()
	defer e.mux.Unlock()
	if utils.StringIsEmpty(key) {
		key = "*"
	}
	e.queues[key] = c
}

// GetQueueAdapter 获取队列适配器
func (e *Application) GetQueueAdapter() storage.AdapterQueue {
	return e.GetQueueKey("*")
}

func (e *Application) GetQueueKey(key string) storage.AdapterQueue {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.queues[key]
}

// GetQueuePrefix 获取带租户标记的queue
func (e *Application) GetQueuePrefix(key string, prefix string) storage.AdapterQueue {
	return NewQueue(prefix, e.GetQueueKey(key))
}

// SetLockerAdapter 设置分布式锁
func (e *Application) SetLockerAdapter(c storage.AdapterLocker) {
	e.locker = c
}

// GetLockerAdapter 获取分布式锁
func (e *Application) GetLockerAdapter() storage.AdapterLocker {
	return NewLocker("", e.locker)
}

func (e *Application) GetLockerPrefix(key string) storage.AdapterLocker {
	return NewLocker(key, e.locker)
}

// SetSmsAdapter 设置缓存
func (e *Application) SetSmsAdapter(key string, c message.AdapterSms) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.sms[key] = c
}

// GetSmsAdapter 获取缓存
func (e *Application) GetSmsAdapter() message.AdapterSms {
	return e.GetSmsKey("*")
}

// GetSmsAdapters 获取缓存
func (e *Application) GetSmsAdapters() map[string]message.AdapterSms {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.sms
}

// GetSmsKey 获取带租户标记的sms
func (e *Application) GetSmsKey(key string) message.AdapterSms {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.sms[key]
}

// SetMailAdapter 设置缓存
func (e *Application) SetMailAdapter(key string, c message.AdapterMail) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.mail[key] = c
}

// GetMailAdapter 获取缓存
func (e *Application) GetMailAdapter() message.AdapterMail {
	return e.GetMailKey("*")
}

// GetMailAdapters 获取缓存
func (e *Application) GetMailAdapters() map[string]message.AdapterMail {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.mail
}

// GetMailKey 获取带租户标记的mail
func (e *Application) GetMailKey(key string) message.AdapterMail {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.mail[key]
}

// SetFileStoreAdapter 设置缓存
func (e *Application) SetFileStoreAdapter(key string, c storage.AdapterFileStore) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.fileStores[key] = c
}

// GetFileStoreAdapter 获取缓存
func (e *Application) GetFileStoreAdapter() storage.AdapterFileStore {
	return e.GetFileStoreKey("*")
}

// GetFileStoreAdapters 获取缓存
func (e *Application) GetFileStoreAdapters() map[string]storage.AdapterFileStore {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.fileStores
}

// GetFileStoreKey 获取带租户标记的cos
func (e *Application) GetFileStoreKey(key string) storage.AdapterFileStore {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.fileStores[key]
}

// SetModerationAdapter 设置缓存
func (e *Application) SetModerationAdapter(key string, c moderation.AdapterModeration) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.moderation[key] = c
}

// GetModerationAdapter 获取缓存
func (e *Application) GetModerationAdapter() moderation.AdapterModeration {
	return e.GetModerationKey("*")
}

// GetModerationAdapters 获取缓存
func (e *Application) GetModerationAdapters() map[string]moderation.AdapterModeration {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.moderation
}

// GetModerationKey 获取带租户标记的cos
func (e *Application) GetModerationKey(key string) moderation.AdapterModeration {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.moderation[key]
}

// SetAmqpAdapter 设置缓存
func (e *Application) SetAmqpAdapter(key string, c message.AdapterAmqp) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.amqp[key] = c
}

// GetAmqpAdapter 获取缓存
func (e *Application) GetAmqpAdapter() message.AdapterAmqp {
	return e.GetAmqpKey("*")
}

// GetAmqpAdapters 获取缓存
func (e *Application) GetAmqpAdapters() map[string]message.AdapterAmqp {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.amqp
}

// GetAmqpKey 获取带租户标记的amqp
func (e *Application) GetAmqpKey(key string) message.AdapterAmqp {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.amqp[key]
}

// SetMqttAdapter 设置MQTT适配器
func (e *Application) SetMqttAdapter(key string, c message.AdapterMqtt) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.mqtt[key] = c
}

// GetMqttAdapter 获取MQTT适配器
func (e *Application) GetMqttAdapter() message.AdapterMqtt {
	return e.GetMqttKey("*")
}

// GetMqttAdapters 获取MQTT适配器列表
func (e *Application) GetMqttAdapters() map[string]message.AdapterMqtt {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.mqtt
}

// GetMqttKey 获取对应key的MQTT适配器
func (e *Application) GetMqttKey(key string) message.AdapterMqtt {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.mqtt[key]
}

// SetThirdPartyAdapter 设置缓存
func (e *Application) SetThirdPartyAdapter(key string, c third_party.AdapterThirdParty) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.thirdParty[key] = c
}

// GetThirdPartyAdapter 获取缓存
func (e *Application) GetThirdPartyAdapter() third_party.AdapterThirdParty {
	return e.GetThirdPartyKey("*")
}

// GetThirdPartyAdapters 获取缓存
func (e *Application) GetThirdPartyAdapters() map[string]third_party.AdapterThirdParty {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.thirdParty
}

// GetThirdPartyKey 获取带租户标记的amqp
func (e *Application) GetThirdPartyKey(key string) third_party.AdapterThirdParty {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.thirdParty[key]
}

// SetLocationBasedServiceAdapter 设置LBS
func (e *Application) SetLocationBasedServiceAdapter(key string, c lbs.AdapterLocationBasedService) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.lbs[key] = c
}

// GetLocationBasedServiceAdapter 获取LBS
func (e *Application) GetLocationBasedServiceAdapter() lbs.AdapterLocationBasedService {
	return e.GetLocationBasedServiceKey("*")
}

// GetLocationBasedServiceAdapters 获取LBS
func (e *Application) GetLocationBasedServiceAdapters() map[string]lbs.AdapterLocationBasedService {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.lbs
}

// GetLocationBasedServiceKey 获取LBS
func (e *Application) GetLocationBasedServiceKey(key string) lbs.AdapterLocationBasedService {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.lbs[key]
}

// SetPaymentServiceAdapter 设置支付
func (e *Application) SetPaymentServiceAdapter(key string, c payment.AdapterPaymentService) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.payment[key] = c
}

// GetPaymentServiceAdapter 获取支付
func (e *Application) GetPaymentServiceAdapter() payment.AdapterPaymentService {
	return e.GetPaymentServiceKey("*")
}

// GetPaymentServiceAdapters 获取支付
func (e *Application) GetPaymentServiceAdapters() map[string]payment.AdapterPaymentService {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.payment
}

// GetPaymentServiceKey 获取支付
func (e *Application) GetPaymentServiceKey(key string) payment.AdapterPaymentService {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.payment[key]
}

// SetLeshuaServiceAdapter 设置乐刷
func (e *Application) SetLeshuaServiceAdapter(key string, c payment.AdapterLeshuaService) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.leshua[key] = c
}

// GetLeshuaServiceAdapter 获取乐刷
func (e *Application) GetLeshuaServiceAdapter() payment.AdapterLeshuaService {
	return e.GetLeshuaServiceKey("*")
}

// GetLeshuaServiceAdapters 获取乐刷
func (e *Application) GetLeshuaServiceAdapters() map[string]payment.AdapterLeshuaService {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.leshua
}

// GetLeshuaServiceKey 获取乐刷
func (e *Application) GetLeshuaServiceKey(key string) payment.AdapterLeshuaService {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.leshua[key]
}

// SetZegoServiceAdapter 设置即构
func (e *Application) SetZegoServiceAdapter(key string, c rtc.AdapterZegoService) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.zego[key] = c
}

// GetZegoServiceAdapter 获取即构
func (e *Application) GetZegoServiceAdapter() rtc.AdapterZegoService {
	return e.GetZegoServiceKey("*")
}

// GetZegoServiceAdapters 获取即构
func (e *Application) GetZegoServiceAdapters() map[string]rtc.AdapterZegoService {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.zego
}

// GetZegoServiceKey 获取即构
func (e *Application) GetZegoServiceKey(key string) rtc.AdapterZegoService {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.zego[key]
}

func (e *Application) SetHandler(key string, routerGroup func(r *gin.RouterGroup, hand ...*gin.HandlerFunc)) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.handler[key] = append(e.handler[key], routerGroup)
}

func (e *Application) GetHandler() map[string][]func(r *gin.RouterGroup, hand ...*gin.HandlerFunc) {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.handler
}

func (e *Application) GetHandlerPrefix(key string) []func(r *gin.RouterGroup, hand ...*gin.HandlerFunc) {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.handler[key]
}

// GetStreamMessage 获取队列需要用的message
func (e *Application) GetStreamMessage(id, stream string, value map[string]interface{}) (storage.Messager, error) {
	message := &queue.Message{}
	message.SetID(id)
	message.SetStream(stream)
	message.SetValues(value)
	return message, nil
}

func (e *Application) GetMemoryQueue(prefix string) storage.AdapterQueue {
	return NewQueue(prefix, e.memoryQueue)
}

// SetConfig 设置对应key的config
func (e *Application) SetConfig(key string, value interface{}) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.configs[key] = value
}

// GetConfig 获取对应key的config
func (e *Application) GetConfig(key string) interface{} {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.configs[key]
}

// SetAppRouters 设置app的路由
func (e *Application) SetAppRouters(appRouters func()) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.appRouters = append(e.appRouters, appRouters)
}

// GetAppRouters 获取app的路由
func (e *Application) GetAppRouters() []func() {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.appRouters
}
