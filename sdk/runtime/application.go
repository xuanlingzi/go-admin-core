package runtime

import (
	"github.com/gin-gonic/gin"
	"github.com/xuanlingzi/go-admin-core/block_chain"
	"github.com/xuanlingzi/go-admin-core/lbs"
	"github.com/xuanlingzi/go-admin-core/message"
	"github.com/xuanlingzi/go-admin-core/third_party"
	"net/http"
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/robfig/cron/v3"
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
	queue       storage.AdapterQueue
	locker      storage.AdapterLocker
	memoryQueue storage.AdapterQueue
	fileStores  map[string]storage.AdapterFileStore
	sms         map[string]message.AdapterSms
	mail        map[string]message.AdapterMail
	amqp        map[string]message.AdapterAmqp
	thirdParty  map[string]third_party.AdapterThirdParty
	blockChain  map[string]block_chain.AdapterBroker
	lbs         map[string]lbs.AdapterLocationBasedService
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
	e.mux.Lock()
	defer e.mux.Unlock()
	return e.dbs
}

// GetDbByKey 根据key获取db
func (e *Application) GetDbByKey(key string) *gorm.DB {
	e.mux.Lock()
	defer e.mux.Unlock()
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
	return e.casbins
}

// GetCasbinKey 根据key获取casbin
func (e *Application) GetCasbinKey(key string) *casbin.SyncedEnforcer {
	e.mux.Lock()
	defer e.mux.Unlock()
	if e, ok := e.casbins["*"]; ok {
		return e
	}
	return e.casbins[key]
}

// SetEngine 设置路由引擎
func (e *Application) SetEngine(engine http.Handler) {
	e.engine = engine
}

// GetEngine 获取路由引擎
func (e *Application) GetEngine() http.Handler {
	return e.engine
}

// GetRouter 获取路由表
func (e *Application) GetRouter() []Router {
	return e.setRouter()
}

// setRouter 设置路由表
func (e *Application) setRouter() []Router {
	switch e.engine.(type) {
	case *gin.Engine:
		routers := e.engine.(*gin.Engine).Routes()
		for _, router := range routers {
			e.routers = append(e.routers, Router{RelativePath: router.Path, Handler: router.Handler, HttpMethod: router.Method})
		}
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
		memoryQueue: queue.NewMemory(10000),
		fileStores:  make(map[string]storage.AdapterFileStore),
		sms:         make(map[string]message.AdapterSms),
		mail:        make(map[string]message.AdapterMail),
		amqp:        make(map[string]message.AdapterAmqp),
		thirdParty:  make(map[string]third_party.AdapterThirdParty),
		blockChain:  make(map[string]block_chain.AdapterBroker),
		lbs:         make(map[string]lbs.AdapterLocationBasedService),
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
	e.mux.Lock()
	defer e.mux.Unlock()
	return e.crontab
}

// GetCrontabKey 根据key获取crontab
func (e *Application) GetCrontabKey(key string) *cron.Cron {
	e.mux.Lock()
	defer e.mux.Unlock()
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
	return e.middlewares
}

// GetMiddlewareKey 获取对应key的中间件
func (e *Application) GetMiddlewareKey(key string) interface{} {
	e.mux.Lock()
	defer e.mux.Unlock()
	return e.middlewares[key]
}

// SetCacheAdapter 设置缓存
func (e *Application) SetCacheAdapter(c storage.AdapterCache) {
	e.cache = c
}

// GetCacheAdapter 获取缓存
func (e *Application) GetCacheAdapter() storage.AdapterCache {
	return NewCache("", e.cache, "")
}

// GetCachePrefix 获取带租户标记的cache
func (e *Application) GetCachePrefix(key string) storage.AdapterCache {
	return NewCache(key, e.cache, "")
}

// SetQueueAdapter 设置队列适配器
func (e *Application) SetQueueAdapter(c storage.AdapterQueue) {
	e.queue = c
}

// GetQueueAdapter 获取队列适配器
func (e *Application) GetQueueAdapter() storage.AdapterQueue {
	return NewQueue("", e.queue)
}

// GetQueuePrefix 获取带租户标记的queue
func (e *Application) GetQueuePrefix(key string) storage.AdapterQueue {
	return NewQueue(key, e.queue)
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
	e.mux.Lock()
	defer e.mux.Unlock()
	return e.sms
}

// GetSmsKey 获取带租户标记的sms
func (e *Application) GetSmsKey(key string) message.AdapterSms {
	e.mux.Lock()
	defer e.mux.Unlock()
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
	e.mux.Lock()
	defer e.mux.Unlock()
	return e.mail
}

// GetMailKey 获取带租户标记的mail
func (e *Application) GetMailKey(key string) message.AdapterMail {
	e.mux.Lock()
	defer e.mux.Unlock()
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
	e.mux.Lock()
	defer e.mux.Unlock()
	return e.fileStores
}

// GetFileStoreKey 获取带租户标记的cos
func (e *Application) GetFileStoreKey(key string) storage.AdapterFileStore {
	e.mux.Lock()
	defer e.mux.Unlock()
	return e.fileStores[key]
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
	return e.amqp
}

// GetAmqpKey 获取带租户标记的amqp
func (e *Application) GetAmqpKey(key string) message.AdapterAmqp {
	e.mux.Lock()
	defer e.mux.Unlock()
	return e.amqp[key]
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
	return e.thirdParty
}

// GetThirdPartyKey 获取带租户标记的amqp
func (e *Application) GetThirdPartyKey(key string) third_party.AdapterThirdParty {
	e.mux.Lock()
	defer e.mux.Unlock()
	return e.thirdParty[key]
}

// SetBlockChainAdapter 设置缓存
func (e *Application) SetBlockChainAdapter(key string, c block_chain.AdapterBroker) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.blockChain[key] = c
}

// GetBlockChainAdapter 获取缓存
func (e *Application) GetBlockChainAdapter() block_chain.AdapterBroker {
	return e.GetBlockChainKey("*")
}

// GetBlockChainAdapters 获取缓存
func (e *Application) GetBlockChainAdapters() map[string]block_chain.AdapterBroker {
	return e.blockChain
}

// GetBlockChainKey 获取带租户标记
func (e *Application) GetBlockChainKey(key string) block_chain.AdapterBroker {
	e.mux.Lock()
	defer e.mux.Unlock()
	return e.blockChain[key]
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
	return e.lbs
}

// GetLocationBasedServiceKey 获取LBS
func (e *Application) GetLocationBasedServiceKey(key string) lbs.AdapterLocationBasedService {
	e.mux.Lock()
	defer e.mux.Unlock()
	return e.lbs[key]
}

func (e *Application) SetHandler(key string, routerGroup func(r *gin.RouterGroup, hand ...*gin.HandlerFunc)) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.handler[key] = append(e.handler[key], routerGroup)
}

func (e *Application) GetHandler() map[string][]func(r *gin.RouterGroup, hand ...*gin.HandlerFunc) {
	e.mux.Lock()
	defer e.mux.Unlock()
	return e.handler
}

func (e *Application) GetHandlerPrefix(key string) []func(r *gin.RouterGroup, hand ...*gin.HandlerFunc) {
	e.mux.Lock()
	defer e.mux.Unlock()
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
	e.mux.Lock()
	defer e.mux.Unlock()
	return e.configs[key]
}

// SetAppRouters 设置app的路由
func (e *Application) SetAppRouters(appRouters func()) {
	e.appRouters = append(e.appRouters, appRouters)
}

// GetAppRouters 获取app的路由
func (e *Application) GetAppRouters() []func() {
	return e.appRouters
}
