package recovery

import (
	"KyokaSuigetsu/internal/types"
	"KyokaSuigetsu/pkg/recovery/config"
	"KyokaSuigetsu/pkg/recovery/fakemaster"
	"KyokaSuigetsu/pkg/util"
	"github.com/go-mysql-org/go-mysql/mysql"
	"math/rand"
)

type Recovery struct {
	fakeMaster map[string]*fakemaster.FakeMaster
}

func NewRecovery() *Recovery {
	return &Recovery{
		fakeMaster: make(map[string]*fakemaster.FakeMaster),
	}
}

func (r *Recovery) RegisterNewFakeMaster(info *types.RecoveryInfo) (*types.FakeMasterInfo, error) {
	if info.ServerID == 0 {
		info.ServerID = rand.Intn(256)
	}
	heroine, err := util.MaKeRuNa()
	if err != nil {
		return nil, err
	}
	fm := &fakemaster.FakeMaster{
		Config: &config.Config{
			ServerVersion:     info.ServerVersion,
			ServerID:          info.ServerID,
			BinlogDir:         info.BinlogDir,
			ServerPort:        heroine.LuckyNumber,
			ReplicateUser:     heroine.Name,
			ReplicatePassword: heroine.Characteristic,
			UntilTimestamp:    0,
			UntilGTID:         mysql.UUIDSet{},
		},
		Variable: fakemaster.NewVariables(),
		Sessions: make(map[int64]*fakemaster.Session),
	}
	// Set nesseray Global variables
	// TODO: NEED TO MODIFY SERVER UUID!!!
	fm.Variable.SetVariable("server_uuid", "deadbeef-1018-1219-0226-123456789012")
	fm.Variable.SetVariable("server_id", info.ServerID)
	fm.Variable.SetVariable("binlog_checksum", "CRC32")
	fm.Variable.SetVariable("gtid_mode", "ON")

	r.fakeMaster[heroine.Name] = fm

	fakeMaster := &types.FakeMasterInfo{
		User:     heroine.Name,
		Password: heroine.Characteristic,
		Port:     heroine.LuckyNumber,
	}
	return fakeMaster, nil
}

func (r *Recovery) UntilTimestamp(info *types.FakeMasterInfo, ts int64) {
	fake := r.fakeMaster[info.User]
	fake.Config.UntilTimestamp = ts
}

func (r *Recovery) UntilGTID(info *types.FakeMasterInfo, set mysql.UUIDSet) {
	fake := r.fakeMaster[info.User]
	fake.Config.UntilGTID = set
}

func (r *Recovery) Start(info *types.FakeMasterInfo) error {
	fake := r.fakeMaster[info.User]
	err := fake.Run()
	if err != nil {
		return err
	}
	return nil
}
