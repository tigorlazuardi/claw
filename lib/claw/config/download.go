package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/adrg/xdg"
	"github.com/dustin/go-humanize"
)

type Download struct {
	BaseDir      string       `koanf:"base_dir"`
	TmpDir       string       `koanf:"tmp_dir"`
	StallMonitor StallMonitor `koanf:"stall_monitor"`
}

func DefaultDownload() Download {
	return Download{
		BaseDir:      filepath.Join(xdg.UserDirs.Pictures, "claw"),
		TmpDir:       filepath.Join(os.TempDir(), "claw"),
		StallMonitor: DefaultStallMonitor(),
	}
}

type StallMonitor struct {
	// Enabled indicates whether the download monitor is enabled
	Enabled bool `koanf:"enabled"`
	// Speed is the speed threshold (in bytes per second) below which a download is considered stalled.
	Speed ByteSize `koanf:"threshold_speed"`
	// SpeedDuration is the duration for which the download speed must remain below the threshold before being considered stalled.
	SpeedDuration time.Duration `koanf:"threshold_duration"`
	// NoDataReceivedDuration is the duration for which no data is received before considering the connection stalled.
	//
	// Default: 10 seconds.
	NoDataReceivedDuration time.Duration `koanf:"no_data_recived_duration"`
}

type ByteSize uint64

func (b *ByteSize) UnmarshalText(text []byte) error {
	size, err := humanize.ParseBytes(string(text))
	if err != nil {
		return fmt.Errorf("failed to parse byte size: %w", err)
	}
	*b = ByteSize(size)
	return nil
}

func (b ByteSize) MarshalText() ([]byte, error) {
	return []byte(b.String()), nil
}

func (b ByteSize) String() string {
	return humanize.Bytes(uint64(b))
}

func DefaultStallMonitor() StallMonitor {
	return StallMonitor{
		Enabled:       true,
		Speed:         10 * 1024, // 10 KB/s
		SpeedDuration: 10 * time.Second,
	}
}
