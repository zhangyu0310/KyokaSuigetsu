package handler

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sync/atomic"

	"github.com/go-mysql-org/go-mysql/replication"
	log "github.com/zhangyu0310/zlogger"
)

type ReplicateMode uint8

const (
	ReplicateModeNone      ReplicateMode = iota
	ReplicateModeBinlogPos               // binlog position
	ReplicateModeGTID                    // gtid
)

var (
	ErrFirstEventInvalid = errors.New("first event invalid")
	ErrBrokenPipe        = errors.New("broken pipe")
)

type DataSource interface {
	io.Reader
	io.Seeker
}

type BinlogReader struct {
	binlogDir string

	binlogName string
	binlogPos  int64
	dataSource DataSource

	parser *replication.BinlogParser

	handler *Handler

	stop atomic.Value
}

func NewBinlogReader(h *Handler, binlogDir string) *BinlogReader {
	r := &BinlogReader{
		binlogDir: binlogDir,
		parser:    replication.NewBinlogParser(),
		handler:   h,
	}
	r.stop.Store(true)
	return r
}

func (r *BinlogReader) SetReplicateBinlogInfo(name string, pos int64) {
	r.binlogName = name
	r.binlogPos = pos
}

func (r *BinlogReader) SetDataSource(ds DataSource) {
	r.dataSource = ds
}

func (r *BinlogReader) Start() error {
	var err error
	// TODO: Set data source by SetDataSource.
	binlogPath := fmt.Sprintf("%s/%s", r.binlogDir, r.binlogName)
	r.dataSource, err = os.OpenFile(binlogPath, os.O_RDONLY, 0666)
	if err != nil {
		log.ErrorF("Open binlog file %s failed, err: %s",
			binlogPath, err)
		return err
	}
	// Seek pos 4 to get FORMAT_DESCRIPTION_EVENT
	_, err = r.dataSource.Seek(4, io.SeekStart)
	if err != nil {
		log.Error("Seek binlog to 4 failed, err:", err)
		return err
	}
	eof, err := r.parser.ParseSingleEvent(r.dataSource,
		func(event *replication.BinlogEvent) error {
			if event.Header.EventType != replication.FORMAT_DESCRIPTION_EVENT {
				log.Error("First event is not FORMAT_DESCRIPTION_EVENT.")
				return ErrFirstEventInvalid
			}
			return nil
		})
	if err != nil {
		log.Error("Parse First event failed, err:", err)
		return err
	}
	if eof {
		return io.EOF
	}
	// Seek to replicate position.
	_, err = r.dataSource.Seek(r.binlogPos, io.SeekStart)
	if err != nil {
		log.ErrorF("Seek binlog to %d failed, err: %s",
			r.binlogPos, err)
	}
	r.stop.Store(false)
	// FIXME: 需要能够cover各种情况的错误，无法恢复再终止
	streamer := r.handler.binlogStreamer
	for !r.stop.Load().(bool) {
		r.dataSource.Read()
		r.parser.Parse()
		eof, err = r.parser.ParseSingleEvent(r.dataSource,
			func(event *replication.BinlogEvent) error {
				return streamer.AddEventToStreamer(event)
			})
		if err != nil {
			if err == ErrBrokenPipe {
				log.Error("Pipe broken, can't continue.")
				return err
			}
		}
		if eof {
			return io.EOF
		}
	}
	return nil
}

func (r *BinlogReader) Stop() {
	r.stop.Store(true)
}
