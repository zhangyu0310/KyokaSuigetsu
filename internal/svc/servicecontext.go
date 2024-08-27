package svc

import (
	"KyokaSuigetsu/internal/config"
	"KyokaSuigetsu/pkg/recovery"
)

type ServiceContext struct {
	Config config.Config

	// Back *back.Back

	Recovery *recovery.Recovery
}

func NewServiceContext(c config.Config) *ServiceContext {
	r := recovery.NewRecovery()
	return &ServiceContext{
		Config: c,

		Recovery: r,
	}
}
