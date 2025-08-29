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
	BaseDir      string          `koanf:"base_dir"`
	TmpDir       string          `koanf:"tmp_dir"`
	StallMonitor DownloadMonitor `koanf:"stall_monitor"`
}

func DefaultDownload() Download {
	return Download{
		BaseDir:      filepath.Join(xdg.UserDirs.Pictures, "claw"),
		TmpDir:       filepath.Join(os.TempDir(), "claw"),
		StallMonitor: DefaultDownloadMonitor(),
	}
}

type DownloadMonitor struct {
	// Enabled indicates whether the download monitor is enabled
	Enabled bool `koanf:"enabled"`
	// Speed is the speed threshold (in bytes per second) below which a download is considered stalled.
	Speed ByteSize `koanf:"threshold_speed"`
	// Duration is the duration for which the download speed must remain below the threshold before being considered stalled.
	Duration time.Duration `koanf:"threshold_duration"`
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

func DefaultDownloadMonitor() DownloadMonitor {
	return DownloadMonitor{
		Enabled:  true,
		Speed:    10 * 1024, // 10 KB/s
		Duration: 10 * time.Second,
	}
}
