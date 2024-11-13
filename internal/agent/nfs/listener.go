package nfs

import (
	"net"
	"strings"
)

type FilteredListener struct {
	net.Listener
	AllowedIP string
}

func (fl *FilteredListener) Accept() (net.Conn, error) {
	for {
		conn, err := fl.Listener.Accept()
		if err != nil {
			return nil, err
		}

		if strings.Contains(conn.RemoteAddr().String(), fl.AllowedIP) {
			return conn, nil
		}

		conn.Close()
	}
}
