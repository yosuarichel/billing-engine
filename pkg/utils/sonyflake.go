package utils

import (
	"fmt"
	"sync"
	"time"

	"github.com/sony/sonyflake/v2"
)

var (
	sf   *sonyflake.Sonyflake
	once sync.Once
)

func initSonyflake() {
	settings := sonyflake.Settings{
		StartTime: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		MachineID: func() (int, error) {
			return GetMachineID()
		},
	}

	res, err := sonyflake.New(settings)
	if err != nil {
		panic("gagal inisialisasi sonyflake")
	}

	sf = res
}

func InitSonyflakeCluster() {
	once.Do(func() {
		initSonyflake()
	})
}

func GenerateSonyflakeID() int64 {
	if sf == nil {
		panic("sonyflake not initialized yet, call InitSonyflakeCluster first")
	}

	id, err := sf.NextID()
	if err != nil {
		panic(fmt.Errorf("failed to generate sonyflake id: %w", err))
	}
	return id
}
