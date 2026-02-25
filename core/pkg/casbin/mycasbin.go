package mycasbin

import (
	"github.com/redis/go-redis/v9"
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/log"
	"github.com/casbin/casbin/v2/model"
	"github.com/xuanlingzi/go-admin-core/core"
	"github.com/xuanlingzi/go-admin-core/core/config"
	"github.com/xuanlingzi/go-admin-core/logger"
	redisWatcher "github.com/xuanlingzi/redis-watcher/v2"
	"gorm.io/gorm"

	gormAdapter "github.com/xuanlingzi/gorm-adapter/v3"
)

// Initialize the model from a string.
var text = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && (keyMatch2(r.obj, p.obj) || keyMatch(r.obj, p.obj)) && (r.act == p.act || p.act == "*")
`

var (
	enforcer *casbin.SyncedEnforcer
	once     sync.Once
)

func Setup(db *gorm.DB, _ string) (*casbin.SyncedEnforcer, error) {
	var setupErr error
	once.Do(func() {
		Apter, err := gormAdapter.NewAdapterByDBUseTableName(db, "sys", "casbin_rule")
		if err != nil && err.Error() != "invalid DDL" {
			setupErr = err
			return
		}

		m, err := model.NewModelFromString(text)
		if err != nil {
			setupErr = err
			return
		}
		enforcer, err = casbin.NewSyncedEnforcer(m, Apter)
		if err != nil {
			setupErr = err
			return
		}
		err = enforcer.LoadPolicy()
		if err != nil {
			setupErr = err
			return
		}
		// set redis watcher if redis config is not nil
		if config.CacheConfig.Redis != nil {
			w, err := redisWatcher.NewWatcher(config.CacheConfig.Redis.Addr, redisWatcher.WatcherOptions{
				Options: redis.Options{
					Network:  "tcp",
					Password: config.CacheConfig.Redis.Password,
				},
				Channel:    "/casbin",
				IgnoreSelf: false,
			})
			if err != nil {
				setupErr = err
				return
			}

			err = w.SetUpdateCallback(updateCallback)
			if err != nil {
				setupErr = err
				return
			}
			err = enforcer.SetWatcher(w)
			if err != nil {
				setupErr = err
				return
			}
		}

		log.SetLogger(&Logger{})
		enforcer.EnableLog(true)
	})

	if setupErr != nil {
		return nil, setupErr
	}
	return enforcer, nil
}

func updateCallback(msg string) {
	l := logger.NewHelper(core.Runtime.GetLogger())
	l.Infof("casbin updateCallback msg: %v", msg)
	err := enforcer.LoadPolicy()
	if err != nil {
		l.Errorf("casbin LoadPolicy err: %v", err)
	}
}
