package binlog

import (
	"KyokaSuigetsu/pkg/recovery/config"
	"errors"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
	"os"
	"path"
)

type PositionMode uint8

const (
	PositionModeNone      PositionMode = iota
	PositionModeBinlogPos              // binlog position
	PositionModeGTID                   // gtid
)

type Reader struct {
	Config *config.Config

	mode     PositionMode
	position mysql.Position
	gtidSet  *mysql.MysqlGTIDSet

	streamer *replication.BinlogStreamer
}

func NewReader(cfg *config.Config) *Reader {
	return &Reader{
		Config: cfg,

		mode: PositionModeNone,
	}
}

func (r *Reader) SetPosition(pos mysql.Position) {
	r.mode = PositionModeBinlogPos
	r.position = pos
}

func (r *Reader) AutoPosition(set *mysql.MysqlGTIDSet) {
	r.mode = PositionModeGTID
	r.gtidSet = set
}

func (r *Reader) Start() (*replication.BinlogStreamer, error) {
	streamer := replication.NewBinlogStreamer()
	r.streamer = streamer
	// find first position.
	entries, err := os.ReadDir(r.Config.BinlogDir)
	if err != nil {
		return nil, err
	}
	targetFiles := make([]os.DirEntry, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() && entry.Type().IsRegular() {
			targetFiles = append(targetFiles, entry)
		}
	}
	if len(targetFiles) == 0 {
		return nil, errors.New("empty binlog dir")
	}

	go func() {
		parser := replication.NewBinlogParser()
		// If we don't give fuck about database info. Just use RawMode.
		parser.SetRawMode(true)
		targetIndex := 0
		for {
			if targetIndex {
				break
			}
			f, err := os.Open(path.Join(r.Config.BinlogDir, targetFiles[targetIndex].Name()))
			if err != nil {
				r.streamer.AddErrorToStreamer(err)
			} else {
				err = parser.ParseReader(f, func(event *replication.BinlogEvent) error {
					if event.Header.Timestamp {

					}
					if event.Header.EventType {

					}
					return r.streamer.AddEventToStreamer(event)
				})
				if err != nil {
					r.streamer.AddErrorToStreamer(err)
				}
			}
		}
	}()
	return streamer, nil
}
