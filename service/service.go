/* Date:10/4/23-2023
 * Author:fanzhao
 */
package service

import (
	"context"
	"sync"
)

// Service service.
type Service struct {
	cmsClient     cms.CooperationServiceClient
	jwtClient     *JWTClient
	snowSvr       *snow.Snow //ID gen SDK
	ac            *paladin.Map
	dao           dao.Dao
	accountClient account.AccountClient
	OSS           boss.Client
	conf          *configs.Config
	httpClient    *blademaster.Client
	wg            sync.WaitGroup
}

// New new a service and return.
func New(d dao.Dao) (s *Service, cf func(), err error) {
	var (
		rpcConf struct {
			CMSrpcConfig     *warden.ClientConfig
			AccountrpcConfig *warden.ClientConfig
		}
		// cfg configs.Config
	)
	s = &Service{
		ac:  &paladin.TOML{},
		dao: d,
	}
	cf = s.Close
	if err = paladin.Get("client.toml").UnmarshalTOML(&rpcConf); err != nil {
		panic(err)
	}
	// hot load
	if err = paladin.Watch("application.toml", reloader.NewTomlReloader(&s.conf)); err != nil {
		panic(err)
	}
	if s.cmsClient, err = cms.NewClient(rpcConf.CMSrpcConfig); err != nil {
		panic(err)
	}
	if s.jwtClient = NewJWTClient(s.conf); err != nil {
		panic(err)
	}
	if s.accountClient, err = account.NewClient(rpcConf.AccountrpcConfig); err != nil {
		panic(err)
	}
	s.snowSvr = s.NewSnow(s.conf)
	s.OSS = boss.New(s.conf.Oss)
	s.httpClient = blademaster.NewClient(s.conf.BmClientCfg)
	return
}

func (s *Service) NewSnow(cfg *configs.Config) (snowSvr *snow.Snow) {
	var (
		bizTag   = cfg.SnowFlake.BizTag
		bizToken = cfg.SnowFlake.BizToken
	)
	snowSvr = snow.NewSnow(bizTag, bizToken)
	return
}

// Ping ping the resource.
func (s *Service) Ping(ctx context.Context, e *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, s.dao.Ping(ctx)
}

// Close close the resource.
func (s *Service) Close() {
}
