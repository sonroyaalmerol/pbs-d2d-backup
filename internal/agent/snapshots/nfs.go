//go:build windows
// +build windows

package snapshots

import (
	"fmt"
	"log"
	"net"
	"net/url"

	"github.com/go-git/go-billy/v5/osfs"
	nfs2 "github.com/sonroyaalmerol/pbs-plus/internal/agent/nfs"
	"github.com/willscott/go-nfs"
	"github.com/willscott/go-nfs/helpers"
	"golang.org/x/sys/windows/registry"
)

func (snapshot *WinVSSSnapshot) Serve(port int) error {
	if snapshot.SnapshotPath == "" {
		return fmt.Errorf("Snapshot path is empty")
	}

	baseKey, _, err := registry.CreateKey(registry.LOCAL_MACHINE, "Software\\PBSPlus\\Config", registry.QUERY_VALUE)
	if err != nil {
		return fmt.Errorf("Unable to create registry key -> %v", err)
	}

	defer baseKey.Close()

	var server string
	if server, _, err = baseKey.GetStringValue("ServerURL"); err != nil {
		return fmt.Errorf("Unable to get server url -> %v", err)
	}

	serverUrl, err := url.Parse(server)
	if err != nil {
		return fmt.Errorf("failed to parse server IP: %v", err)
	}

	listenAt := fmt.Sprintf("0.0.0.0:%d", port)
	listener, err := net.Listen("tcp", listenAt)
	if err != nil {
		return fmt.Errorf("Port is already in use! Failed to listen on %s: %v", listenAt, err)
	}

	listener = &nfs2.FilteredListener{Listener: listener, AllowedIP: serverUrl.Hostname()}

	defer listener.Close()

	fs := osfs.New(snapshot.SnapshotPath)
	readOnlyFs := nfs2.NewROFS(fs)
	nfsHandler := helpers.NewNullAuthHandler(readOnlyFs)

	for {
		done := make(chan struct{})
		go func() {
			err := nfs.Serve(listener, nfsHandler)
			if err != nil {
				log.Printf("NFS server error: %v\n", err)
			}
			close(done)
		}()

		select {
		case <-snapshot.Ctx.Done():
			listener.Close()
			return nil
		case <-done:
		}
	}
}
