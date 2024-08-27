package fakemaster

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/server"
	"github.com/go-mysql-org/go-mysql/test_util/test_keys"
	"github.com/pingcap/tidb/pkg/parser"
	"go.uber.org/zap"
	"net"
)

type Config struct {
	ServerVersion     string
	ServerID          int
	ServerPort        int
	ReplicateUser     string
	ReplicatePassword string

	BinlogDir string

	UntilTimestamp int64
	UntilGTID      mysql.UUIDSet
}

type FakeMaster struct {
	Config *Config

	// parser are used to parse sql
	parser *parser.Parser

	// Global Variables in Server
	Variable *Variables

	Sessions map[int64]*Session
	Ctx      context.Context

	log *zap.Logger
}

type RemoteThrottleProvider struct {
	*server.InMemoryProvider
	delay int // in milliseconds
}

func (fake *FakeMaster) Run() error {
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", fake.Config.ServerPort))
	if err != nil {
		return err
	}
	remoteProvider := &RemoteThrottleProvider{server.NewInMemoryProvider(), 10 + 50}
	remoteProvider.AddUser(fake.Config.ReplicateUser, fake.Config.ReplicatePassword)
	var tlsConf = server.NewServerTLSConfig(test_keys.CaPem, test_keys.CertPem, test_keys.KeyPem, tls.VerifyClientCertIfGiven)

	for {
		c, err := l.Accept()
		if err != nil {
			fake.log.Error("Accept new connection failed.", zap.Error(err))
			continue
		}
		go func() {
			svr := server.NewServer("8.0.12", mysql.DEFAULT_COLLATION_ID, mysql.AUTH_NATIVE_PASSWORD, test_keys.PubPem, tlsConf)
			session := NewSession()
			conn, err := server.NewCustomizedConn(c, svr, remoteProvider, session)
			if err != nil {
				fake.log.Error("Connection error on go mysql", zap.Error(err))
				return
			}

			for {
				err = conn.HandleCommand()
				if err != nil {
					fake.log.Error("Could not handle command", zap.Error(err))
					return
				}
			}
		}()
	}
}

func (fake *FakeMaster) work() {

}
