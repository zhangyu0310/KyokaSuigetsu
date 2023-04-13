package handler

import (
	"errors"
	"io"

	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
	"github.com/pingcap/tidb/parser"
	_ "github.com/pingcap/tidb/parser/test_driver"
	log "github.com/zhangyu0310/zlogger"
)

var (
	ErrNotRegisterSlave = errors.New("not register slave")
	ErrBinlogNameEmpty  = errors.New("binlog name empty")
)

type Handler struct {
	// parser & visitor are used to parse sql
	parser  *parser.Parser
	visitor *Visitor

	// session variables
	variables *Variables

	// For replication.
	registerSlave  bool
	slaveServerId  uint32
	replicateMode  ReplicateMode
	binlogStreamer *replication.BinlogStreamer
	binlogReader   *BinlogReader
}

func NewHandler() *Handler {
	h := &Handler{
		parser:  parser.New(),
		visitor: nil,

		variables: NewVariables(),

		registerSlave:  false,
		slaveServerId:  0,
		replicateMode:  ReplicateModeNone,
		binlogStreamer: replication.NewBinlogStreamer(),
	}
	h.visitor = NewVisitor(h)
	// TODO: binlog dir can be config.
	h.binlogReader = NewBinlogReader(h, "./test")
	return h
}

func (h *Handler) GetVariable(key string) interface{} {
	return h.variables.GetVariable(key)
}

func (h *Handler) SetVariable(key string, value interface{}) {
	h.variables.SetVariable(key, value)
}

func (h *Handler) UseDB(dbName string) error {
	log.Info("use db:", dbName)
	return nil
}

// HandleQuery handle COM_QUERY command, like SELECT, INSERT, UPDATE, etc...
// If Result has a Resultset (SELECT, SHOW, etc...), we will send this as the response, otherwise, we will send Result
func (h *Handler) HandleQuery(query string) (*mysql.Result, error) {
	stmt, err := h.parser.ParseOneStmt(query, "", "")
	if err != nil {
		log.ErrorF("Parse query [%s] error: %s", query, err)
		return nil, err
	}
	h.visitor.Clean()
	stmt.Accept(h.visitor)
	if h.visitor.Error != nil {
		log.ErrorF("Parse query [%s] error: %s", query, h.visitor.Error)
		return nil, h.visitor.Error
	}
	result := h.visitor.Result
	return result, nil
}

// HandleFieldList handle COM_FILED_LIST command
func (h *Handler) HandleFieldList(_ string, _ string) ([]*mysql.Field, error) {
	return nil, nil
}

// HandleStmtPrepare handle COM_STMT_PREPARE, params is the param number for this statement, columns is the column number
// context will be used later for statement execute
func (h *Handler) HandleStmtPrepare(_ string) (params int, columns int, context interface{}, err error) {
	return 0, 0, nil, nil
}

// HandleStmtExecute handle COM_STMT_EXECUTE, context is the previous one set in prepare
// query is the statement prepare query, and args is the params for this statement
func (h *Handler) HandleStmtExecute(_ interface{}, _ string, _ []interface{}) (*mysql.Result, error) {
	return nil, nil
}

// HandleStmtClose handle COM_STMT_CLOSE, context is the previous one set in prepare
// this handler has no response
func (h *Handler) HandleStmtClose(_ interface{}) error {
	return nil
}

// HandleOtherCommand handle any other command that is not currently handled by the library,
// default implementation for this method will return an ER_UNKNOWN_ERROR
func (h *Handler) HandleOtherCommand(_ byte, _ []byte) error {
	return nil
}

// HandleRegisterSlave handle COM_REGISTER_SLAVE command
func (h *Handler) HandleRegisterSlave(data []byte) error {
	// assert h.registerSlave == false
	serverId, pos := GetInt4(data, 0)
	log.Info("Register slave serverId:", serverId)
	reportHost, pos := GetLengthString(data, pos)
	log.Info("Register slave reportHost:", reportHost)
	reportUser, pos := GetLengthString(data, pos)
	log.Info("Register slave reportUser:", reportUser)
	reportPassword, pos := GetLengthString(data, pos)
	log.Info("Register slave reportPassword:", reportPassword)
	reportPort, pos := GetInt2(data, pos)
	log.Info("Register slave reportPort:", reportPort)
	rplRecoveryRank, pos := GetInt4(data, pos)
	log.Info("Register slave rplRecoveryRank:", rplRecoveryRank)
	masterId, pos := GetInt4(data, pos)
	log.Info("Register slave masterId:", masterId)
	h.registerSlave = true
	h.slaveServerId = serverId
	return nil
}

// HandleBinlogDump handle COM_BINLOG_DUMP command
// Start a binlog file reader to put binlog event into BinlogStreamer
// Note: Only support binlog which smaller than 4G
func (h *Handler) HandleBinlogDump(pos mysql.Position) (*replication.BinlogStreamer, error) {
	if !h.registerSlave {
		log.Error("Not register slave!")
		return nil, ErrNotRegisterSlave
	}
	h.replicateMode = ReplicateModeBinlogPos
	if pos.Name == "" {
		return nil, ErrBinlogNameEmpty
	} else {
		h.binlogReader.SetReplicateBinlogInfo(pos.Name, int64(pos.Pos))
	}
	go func() {
		err := h.binlogReader.Start()
		if err != nil {
			if err == io.EOF {
				log.Info("Binlog read EOF.")
			} else {
				log.Error("binlog reader return err:", err)
			}
		}
	}()
	return h.binlogStreamer, nil
}

// HandleBinlogDumpGTID handle COM_BINLOG_DUMP_GTID command
// Start a binlog file reader to put binlog event into BinlogStreamer
func (h *Handler) HandleBinlogDumpGTID(gtidSet *mysql.MysqlGTIDSet) (*replication.BinlogStreamer, error) {
	if !h.registerSlave {
		log.Error("Not register slave!")
		return nil, ErrNotRegisterSlave
	}
	h.replicateMode = ReplicateModeGTID
	return h.binlogStreamer, nil
}
