module github.com/xuanlingzi/go-admin-core/plugins/logger/zap

go 1.22.0

toolchain go1.23.4

require (
	github.com/xuanlingzi/go-admin-core v1.7.14
	go.uber.org/zap v1.27.0
)

require go.uber.org/multierr v1.11.0 // indirect

//replace github.com/xuanlingzi/go-admin-core v1.7.14 => ../../../
