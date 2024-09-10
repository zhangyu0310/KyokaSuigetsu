package fakemaster

import (
	"KyokaSuigetsu/pkg/recovery/binlog"
	"KyokaSuigetsu/pkg/util"
	"errors"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
	"go.uber.org/zap"
)

var (
	ErrNotRegisterSlave = errors.New("not register slave")
	ErrBinlogNameEmpty  = errors.New("binlog name empty")
)

type Session struct {
	fakeMaster *FakeMaster
	visitor    *Visitor
	// Session Variables
	Variable Variables

	// For replication.
	registerSlave bool
	slaveServerId uint32
	binlogReader  *binlog.Reader

	log *zap.Logger
}

func NewSession() *Session {
	s := &Session{
		fakeMaster:    nil,
		visitor:       nil,
		Variable:      Variables{},
		registerSlave: false,
		slaveServerId: 0,
		binlogReader:  nil,
		log:           nil,
	}

	return s
}

func (h *Session) UseDB(dbName string) error {
	h.log.Info("use db.", zap.String("dbName", dbName))
	return nil
}

// HandleQuery handle COM_QUERY command, like SELECT, INSERT, UPDATE, etc...
// If Result has a ResultSet (SELECT, SHOW, etc...), we will send this as the response, otherwise, we will send Result
func (h *Session) HandleQuery(query string) (*mysql.Result, error) {
	stmt, err := h.fakeMaster.parser.ParseOneStmt(query, "", "")
	if err != nil {
		h.log.Error("Parse query error",
			zap.String("query", query),
			zap.Error(err))
		return nil, err
	}
	h.visitor.Clean()
	stmt.Accept(h.visitor)
	if h.visitor.Error != nil {
		h.log.Error("Parse query error",
			zap.String("query", query),
			zap.Error(h.visitor.Error))
		return nil, h.visitor.Error
	}
	result := h.visitor.Result
	return result, nil
}

// HandleFieldList handle COM_FILED_LIST command
func (h *Session) HandleFieldList(_ string, _ string) ([]*mysql.Field, error) {
	return nil, nil
}

// HandleStmtPrepare handle COM_STMT_PREPARE, params is the param number for this statement, columns is the column number
// context will be used later for statement execute
func (h *Session) HandleStmtPrepare(_ string) (params int, columns int, context interface{}, err error) {
	return 0, 0, nil, nil
}

// HandleStmtExecute handle COM_STMT_EXECUTE, context is the previous one set in prepare
// query is the statement prepare query, and args is the params for this statement
func (h *Session) HandleStmtExecute(_ interface{}, _ string, _ []interface{}) (*mysql.Result, error) {
	return nil, nil
}

// HandleStmtClose handle COM_STMT_CLOSE, context is the previous one set in prepare
// this handler has no response
func (h *Session) HandleStmtClose(_ interface{}) error {
	return nil
}

// HandleOtherCommand handle any other command that is not currently handled by the library,
// default implementation for this method will return an ER_UNKNOWN_ERROR
func (h *Session) HandleOtherCommand(_ byte, _ []byte) error {
	return nil
}

// HandleRegisterSlave handle COM_REGISTER_SLAVE command
func (h *Session) HandleRegisterSlave(data []byte) error {
	// assert h.registerSlave == false
	serverId, pos := util.GetInt4(data, 0)
	h.log.Info("Register slave serverId.", zap.Uint32("Server ID", serverId))
	reportHost, pos := util.GetLengthString(data, pos)
	h.log.Info("Register slave reportHost.", zap.String("Slave Host", reportHost))
	reportUser, pos := util.GetLengthString(data, pos)
	h.log.Info("Register slave reportUser.", zap.String("Slave User", reportUser))
	reportPassword, pos := util.GetLengthString(data, pos)
	h.log.Info("Register slave reportPassword.", zap.String("Slave Password", reportPassword))
	reportPort, pos := util.GetInt2(data, pos)
	h.log.Info("Register slave reportPort.", zap.Uint16("Slave Port", reportPort))
	rplRecoveryRank, pos := util.GetInt4(data, pos)
	h.log.Info("Register slave rplRecoveryRank.", zap.Uint32("Slave Rank", rplRecoveryRank))
	masterId, pos := util.GetInt4(data, pos)
	h.log.Info("Register slave masterId.", zap.Uint32("Master ID", masterId))
	h.registerSlave = true
	h.slaveServerId = serverId
	return nil
}

// HandleBinlogDump handle COM_BINLOG_DUMP command
// Start a binlog file reader to put binlog event into BinlogStreamer
// Note: Only support binlog which smaller than 4G
func (h *Session) HandleBinlogDump(pos mysql.Position) (*replication.BinlogStreamer, error) {
	if !h.registerSlave {
		h.log.Error("Not register slave!")
		return nil, ErrNotRegisterSlave
	}
	if pos.Name == "" {
		return nil, ErrBinlogNameEmpty
	} else {
		h.binlogReader.SetPosition(pos)
	}
	return h.binlogReader.Start()
}

// HandleBinlogDumpGTID handle COM_BINLOG_DUMP_GTID command
// Start a binlog file reader to put binlog event into BinlogStreamer
func (h *Session) HandleBinlogDumpGTID(gtidSet *mysql.MysqlGTIDSet) (*replication.BinlogStreamer, error) {
	if !h.registerSlave {
		h.log.Error("Not register slave!")
		return nil, ErrNotRegisterSlave
	}
	h.binlogReader.AutoPosition(gtidSet)
	return h.binlogReader.Start()
}
