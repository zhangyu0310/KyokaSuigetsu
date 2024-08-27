package binlog

import (
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
)

type Reader struct {
	position mysql.Position

	streamer *replication.BinlogStreamer
}

func (r *Reader) SetPosition(pos mysql.Position) {
	r.position = pos
}

func (r *Reader) AutoPosition(set *mysql.MysqlGTIDSet) {

}

func (r *Reader) Start() (*replication.BinlogStreamer, error) {

}
