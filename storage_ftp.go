/*
 * knoxite
 *     Copyright (c) 2016, Stefan Luecke <glaxx@glaxx.net>
 *   For license see LICENSE.txt
 */

package knoxite

import (
	"net/url"
	"time"

	"github.com/jlaffaye/ftp"
)

// StorageFTP stores data on a remote FTP
type StorageFTP struct {
	url      url.URL
	c        *ftp.ServerConn
	loggedIn bool
}

// NewStorageFTP establishs a FTP connection and returns a StorageFTP object.
func NewStorageFTP(u url.URL) (*StorageFTP, error) {
	conn, err := ftp.DialTimeout(u.Host, 30*time.Second)
	if err != nil {
		return nil, err
	}
	loggedIn := false
	if u.User.Username() != "" {
		if pw, set := u.User.Password(); set {
			err := conn.Login(u.User.Username(), pw)
			if err != nil {
				return nil, err
			}
			loggedIn = true
		} else {
			err := conn.Login(u.User.Username(), "")
			if err != nil {
				return nil, err
			}
			loggedIn = true
		}
	} else {
		err := conn.Login("anonymous", "anonymous")
		if err != nil {
			return nil, err
		}
		loggedIn = true
	}

	return &StorageFTP{
		url:      u,
		c:        conn,
		loggedIn: loggedIn,
	}, nil
}

// Location returns the type and location of the repository
func (backend *StorageFTP) Location() string {
	return backend.url.String()
}

// Close the backend
func (backend *StorageFTP) Close() error {
	if backend.loggedIn {
		err := backend.c.Logout()
		if err != nil {
			return err
		}
	}

	return backend.c.Quit()
}

// Protocols returns the Protocol Schemes supported by this backend
func (backend *StorageFTP) Protocols() []string {
	return []string{"ftp"}
}

// Description returns a user-friendly description for this backend
func (backend *StorageFTP) Description() string {
	return "FTP Storage"
}

// LoadChunk loads a Chunk from network
func (backend *StorageFTP) LoadChunk(shasum string, part, totalParts uint) (*[]byte, error) {
	return &[]byte{}, ErrChunkNotFound
}

// StoreChunk stores a single Chunk on network
func (backend *StorageFTP) StoreChunk(shasum string, part, totalParts uint, data *[]byte) (size uint64, err error) {
	return 0, ErrStoreChunkFailed
}

// LoadSnapshot loads a snapshot
func (backend *StorageFTP) LoadSnapshot(id string) ([]byte, error) {
	return []byte{}, ErrSnapshotNotFound
}

// SaveSnapshot stores a snapshot
func (backend *StorageFTP) SaveSnapshot(id string, data []byte) error {
	return ErrStoreSnapshotFailed
}

// InitRepository creates a new repository
func (backend *StorageFTP) InitRepository() error {
	return nil
}

// LoadRepository reads the metadata for a repository
func (backend *StorageFTP) LoadRepository() ([]byte, error) {
	return []byte{}, ErrLoadRepositoryFailed
}

// SaveRepository stores the metadata for a repository
func (backend *StorageFTP) SaveRepository(data []byte) error {
	return ErrStoreRepositoryFailed
}
