package config

import "github.com/go-mysql-org/go-mysql/mysql"

type Config struct {
	ServerVersion     string
	ServerID          int
	ServerPort        int
	ReplicateUser     string
	ReplicatePassword string

	// TODO: Support more data source.
	BinlogDir string

	UntilTimestamp int64
	UntilGTID      mysql.UUIDSet
}
