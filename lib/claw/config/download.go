package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/adrg/xdg"
	"github.com/dustin/go-humanize"
)

type Download struct {
	BaseDir           string       `koanf:"base_dir"`
	TmpDir            string       `koanf:"tmp_dir"`
	StallMonitor      StallMonitor `koanf:"stall_monitor"`
	FilenameMaxLength int          `koanf:"filename_max_length"`
	SanityCheck       SanityCheck  `koanf:"sanity_check"`
}

func (do Download) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("base_dir", do.BaseDir),
		slog.String("tmp_dir", do.TmpDir),
		slog.Int("filename_max_length", do.FilenameMaxLength),
		slog.Any("stall_monitor", do.StallMonitor),
		slog.Any("sanity_check", do.SanityCheck),
	)
}

func DefaultDownload() Download {
	return Download{
		BaseDir:           filepath.Join(xdg.UserDirs.Pictures, "claw"),
		TmpDir:            filepath.Join(os.TempDir(), "claw"),
		StallMonitor:      DefaultStallMonitor(),
		FilenameMaxLength: 100,
		SanityCheck:       DefaultSanityCheck(),
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

func (st StallMonitor) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Bool("enabled", st.Enabled),
		slog.Any("speed", st.Speed),
		slog.Duration("speed_duration", st.SpeedDuration),
		slog.Duration("no_data_received_duration", st.NoDataReceivedDuration),
	)
}

type ByteSize uint64

func (by ByteSize) LogValue() slog.Value {
	return slog.StringValue(by.String())
}

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

type SanityCheck struct {
	Enabled          bool     `koanf:"enabled"`
	MinImageFilesize ByteSize `koanf:"image_filesize"`
}

func DefaultSanityCheck() SanityCheck {
	return SanityCheck{
		Enabled:          true,
		MinImageFilesize: 64 * 1024, // 10 KB
	}
}
