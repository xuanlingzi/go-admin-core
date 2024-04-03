package cronjob

import (
	"github.com/robfig/cron/v3"
	"time"
)

// NewWithSeconds returns a Cron with the seconds field enabled.
func NewWithSeconds(loc *time.Location) *cron.Cron {
	secondParser := cron.NewParser(cron.Second | cron.Minute |
		cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)
	if loc == nil {
		loc = time.Local
	}
	return cron.New(cron.WithParser(secondParser), cron.WithChain(), cron.WithLocation(loc))
}
