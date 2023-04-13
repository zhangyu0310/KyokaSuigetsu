package main

import (
	"crypto/tls"
	"net"
	"time"

	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/server"
	"github.com/go-mysql-org/go-mysql/test_util/test_keys"
	log "github.com/zhangyu0310/zlogger"

	"KyokaSuigetsu/handler"
)

type RemoteThrottleProvider struct {
	*server.InMemoryProvider
	delay int // in milliseconds
}

func (m *RemoteThrottleProvider) GetCredential(username string) (password string, found bool, err error) {
	time.Sleep(time.Millisecond * time.Duration(m.delay))
	return m.InMemoryProvider.GetCredential(username)
}

// SetGlobalVariables set global variables (TODO: variable can be config)
func SetGlobalVariables() {
	handler.GlobalVariables.SetVariable("server_uuid", "deadbeef-1018-1219-0226-123456789012")
	handler.GlobalVariables.SetVariable("server_id", 26)
	handler.GlobalVariables.SetVariable("binlog_checksum", "CRC32")
	handler.GlobalVariables.SetVariable("gtid_mode", "ON")
}

func main() {
	// TODO: use config
	l, _ := net.Listen("tcp", "0.0.0.0:30226")
	remoteProvider := &RemoteThrottleProvider{server.NewInMemoryProvider(), 10 + 50}
	remoteProvider.AddUser("poppinzhang", "poppinzhang")
	var tlsConf = server.NewServerTLSConfig(test_keys.CaPem, test_keys.CertPem, test_keys.KeyPem, tls.VerifyClientCertIfGiven)
	SetGlobalVariables()
	for {
		c, _ := l.Accept()
		go func() {
			svr := server.NewServer("8.0.12", mysql.DEFAULT_COLLATION_ID, mysql.AUTH_NATIVE_PASSWORD, test_keys.PubPem, tlsConf)
			h := handler.NewHandler()
			conn, err := server.NewCustomizedConn(c, svr, remoteProvider, h)
			if err != nil {
				log.ErrorF("Connection error: %v", err)
				return
			}

			for {
				err = conn.HandleCommand()
				if err != nil {
					log.ErrorF(`Could not handle command: %v`, err)
					return
				}
			}
		}()
	}
}
