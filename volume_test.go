/*
 * knoxite
 *     Copyright (c) 2016, Christian Muehlhaeuser <muesli@gmail.com>
 *
 *   For license see LICENSE.txt
 */

package knoxite

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestCreateVolume(t *testing.T) {
	testPassword := "this_is_a_password"

	dir, err := ioutil.TempDir("", "knoxite")
	if err != nil {
		t.Errorf("Failed creating temporary dir for repository: %s", err)
		return
	}
	defer os.RemoveAll(dir)

	vol, verr := NewVolume("test_name", "test_description")
	{
		r, err := NewRepository(dir, testPassword)
		if err != nil {
			t.Errorf("Failed creating repository: %s", err)
			return
		}

		if verr == nil {
			verr = r.AddVolume(vol)
			if verr != nil {
				t.Errorf("Failed creating volume: %s", verr)
				return
			}

			serr := r.Save()
			if serr != nil {
				t.Errorf("Failed saving volume: %s", serr)
				return
			}
		}
	}

	{
		r, err := OpenRepository(dir, testPassword)
		if err != nil {
			t.Errorf("Failed opening repository: %s", err)
			return
		}

		volume, err := r.FindVolume(vol.ID)
		if err != nil {
			t.Errorf("Failed finding volume: %s", err)
			return
		}
		if volume.Name != vol.Name {
			t.Errorf("Failed verifying volume name: %s != %s", vol.Name, volume.Name)
		}
		if volume.Description != vol.Description {
			t.Errorf("Failed verifying volume description: %s != %s", vol.Description, volume.Description)
		}
	}
}

func TestFindUnknownVolume(t *testing.T) {
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

	_, err = r.FindVolume("invalidID")
	if err != ErrVolumeNotFound {
		t.Errorf("Expected %v, got %v", ErrVolumeNotFound, err)
	}
}
