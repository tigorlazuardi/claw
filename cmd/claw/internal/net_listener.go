package internal

import (
	"net"
	"strings"
)

type NetListener struct {
	net.Listener
}

func (ne *NetListener) UnmarshalText(text []byte) error {
	var err error
	tgt := string(text)
	if tgt == "" {
		ne.Listener, err = net.Listen("tcp", ":8000")
		return err
	}
	if strings.HasPrefix(tgt, "unix://") || strings.HasPrefix(tgt, "unix:") {
		tgt = strings.TrimPrefix(tgt, "unix://")
		tgt = strings.TrimPrefix(tgt, "unix:")
		ne.Listener, err = net.Listen("unix", tgt)
		return err
	}
	ne.Listener, err = net.Listen("tcp", tgt)
	return err
}

func (ne *NetListener) MarshalText() ([]byte, error) {
	if ne == nil || ne.Listener == nil {
		return []byte(""), nil
	}
	addr := ne.Listener.Addr()
	if addr == nil {
		return []byte(""), nil
	}
	if addr.Network() == "unix" {
		return []byte("unix://" + addr.String()), nil
	}
	return []byte(addr.String()), nil
}
