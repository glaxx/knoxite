/*
 * knoxite
 *     Copyright (c) 2016, Christian Muehlhaeuser <muesli@gmail.com>
 *
 *   For license see LICENSE.txt
 */

package knoxite

import (
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func shasumFile(path string) (string, error) {
	hasher := sha256.New()
	s, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	hasher.Write(s)
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func TestCreateSnapshot(t *testing.T) {
	testPassword := "this_is_a_password"

	dir, err := ioutil.TempDir("", "knoxite")
	if err != nil {
		t.Errorf("Failed creating temporary dir for repository: %s", err)
		return
	}
	defer os.RemoveAll(dir)

	snapshotOriginal := Snapshot{}
	{
		r, err := NewRepository(dir, testPassword)
		if err != nil {
			t.Errorf("Failed creating repository: %s", err)
			return
		}
		vol, err := NewVolume("test_name", "test_description")
		if err != nil {
			t.Errorf("Failed creating volume: %s", err)
			return
		}
		err = r.AddVolume(vol)
		if err != nil {
			t.Errorf("Failed creating volume: %s", err)
			return
		}
		snapshot, err := NewSnapshot("test_snapshot")
		if err != nil {
			t.Errorf("Failed creating snapshot: %s", err)
			return
		}

		wd, err := os.Getwd()
		if err != nil {
			t.Errorf("Failed getting working dir: %s", err)
			return
		}
		progress, err := snapshot.Add(wd, []string{"snapshot_test.go"}, r, false, true, 1, 0)
		if err != nil {
			t.Errorf("Failed adding to snapshot: %s", err)
		}
		for range progress {
		}

		err = snapshot.Save(&r)
		if err != nil {
			t.Errorf("Failed saving snapshot: %s", err)
		}
		err = vol.AddSnapshot(snapshot.ID)
		if err != nil {
			t.Errorf("Failed adding snapshot to volume: %s", err)
		}
		err = r.Save()
		if err != nil {
			t.Errorf("Failed saving volume: %s", err)
			return
		}

		snapshotOriginal = snapshot
	}

	{
		r, err := OpenRepository(dir, testPassword)
		if err != nil {
			t.Errorf("Failed opening repository: %s", err)
			return
		}

		_, snapshot, err := r.FindSnapshot(snapshotOriginal.ID)
		if err != nil {
			t.Errorf("Failed finding snapshot: %s", err)
			return
		}
		if !snapshot.Date.Equal(snapshotOriginal.Date) {
			t.Errorf("Failed verifying snapshot date: %v != %v", snapshot.Date, snapshotOriginal.Date)
		}
		if snapshot.Description != snapshotOriginal.Description {
			t.Errorf("Failed verifying snapshot description: %s != %s", snapshot.Description, snapshotOriginal.Description)
		}

		for i, item := range snapshot.Items {
			if item.Path != snapshotOriginal.Items[i].Path {
				t.Errorf("Failed verifying snapshot item: %s != %s", item.Path, snapshotOriginal.Items[i].Path)
				return
			}
			if item.Size != snapshotOriginal.Items[i].Size {
				t.Errorf("Failed verifying snapshot item size: %d != %d", item.Size, snapshotOriginal.Items[i].Size)
				return
			}
		}

		targetdir, err := ioutil.TempDir("", "knoxite.target")
		if err != nil {
			t.Errorf("Failed creating temporary dir for restore: %s", err)
			return
		}
		defer os.RemoveAll(targetdir)

		progress, err := DecodeSnapshot(r, *snapshot, targetdir)
		if err != nil {
			t.Errorf("Failed restoring snapshot: %s", err)
			return
		}
		for range progress {
		}

		for i, item := range snapshot.Items {
			file1 := filepath.Join(targetdir, item.Path)
			sha1, err := shasumFile(file1)
			if err != nil {
				t.Errorf("Failed generating shasum for %s: %s", file1, err)
				return
			}
			sha2, err := shasumFile(snapshotOriginal.Items[i].Path)
			if err != nil {
				t.Errorf("Failed generating shasum for %s: %s", snapshotOriginal.Items[i].Path, err)
				return
			}
			if sha1 != sha2 {
				t.Errorf("Failed verifying shasum: %s != %s", sha1, sha2)
				return
			}
		}
	}
}

func TestFindUnknownSnapshot(t *testing.T) {
	testPassword := "this_is_a_password"

	dir, err := ioutil.TempDir("", "knoxite")
	if err != nil {
		t.Errorf("Failed creating temporary dir for repository: %s", err)
		return
	}
	defer os.RemoveAll(dir)

	r, err := NewRepository(dir, testPassword)
	if err != nil {
		t.Errorf("Failed creating repository: %s", err)
		return
	}

	vol, err := NewVolume("test", "")
	if err != nil {
		t.Errorf("Failed creating volume: %s", err)
		return
	}
	r.AddVolume(vol)

	_, _, err = r.FindSnapshot("invalidID")
	if err != ErrSnapshotNotFound {
		t.Errorf("Expected %v, got %v", ErrSnapshotNotFound, err)
	}
}
