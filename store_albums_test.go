package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"sort"
	"testing"
)

func addAlbum(album *Album) error {
	buffer, err := json.Marshal(album)
	if err != nil {
		return err
	}

	resp, err := http.Post(
		TEST_SERVER_END_POINT+"addAlbum",
		"application/x-www-form-urlencoded",
		bytes.NewReader(buffer),
	)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("Expected 200 OK but got " + resp.Status)
	}

	return nil
}

func deleteAlbum(id string) error {
	buffer, err := json.Marshal(id)
	if err != nil {
		return err
	}

	resp, err := http.Post(
		TEST_SERVER_END_POINT+"deleteAlbum",
		"application/x-www-form-urlencoded",
		bytes.NewReader(buffer),
	)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("Expected 200 OK but got " + resp.Status)
	}

	return nil
}

func getAlbum(id string) (*Album, error) {
	buffer, err := json.Marshal(id)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(
		TEST_SERVER_END_POINT+"getAlbum",
		"application/x-www-form-urlencoded",
		bytes.NewReader(buffer),
	)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("Expected 200 OK but got " + resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	album := new(Album)
	err = json.Unmarshal(body, album)
	if err != nil {
		return nil, err
	}

	return album, nil
}

func getAllAlbums() ([]string, error) {
	resp, err := http.Post(
		TEST_SERVER_END_POINT+"getAllAlbums",
		"application/x-www-form-urlencoded",
		nil,
	)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("Expected 200 OK but got " + resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	albums := make([]string, 0)
	err = json.Unmarshal(body, &albums)
	if err != nil {
		return nil, err
	}

	return albums, nil
}

func getArtistAlbums(artistId string) ([]string, error) {
	buffer, err := json.Marshal(artistId)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(
		TEST_SERVER_END_POINT+"getArtistAlbums",
		"application/x-www-form-urlencoded",
		bytes.NewReader(buffer),
	)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("Expected 200 OK but got " + resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	albums := make([]string, 0)
	err = json.Unmarshal(body, &albums)
	if err != nil {
		return nil, err
	}

	return albums, nil
}

func updateAlbum(album *Album) error {
	buffer, err := json.Marshal(album)
	if err != nil {
		return err
	}

	resp, err := http.Post(
		TEST_SERVER_END_POINT+"updateAlbum",
		"application/x-www-form-urlencoded",
		bytes.NewReader(buffer),
	)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("Expected 200 OK but got " + resp.Status)
	}

	return nil
}

func TestAddAlbum(test *testing.T) {
	album := Album{
		Id:       "testAddId",
		Name:     "testAdd",
		Price:    "100",
		ArtistId: "testAddAlbumArtist",
	}

	test.Log("Adding album")
	err := addAlbum(&album)
	if err != nil {
		test.Errorf("Unable to add album %#v: %s", album, err)
		test.FailNow()
	}
}

func TestDeleteAlbum(test *testing.T) {
	// First, add the album.
	album := Album{
		Id:       "testDeleteAlbumId",
		Name:     "testDelete",
		Price:    "100",
		ArtistId: "testAddAlbumArtist",
	}

	err := addAlbum(&album)
	if err != nil {
		test.Errorf("Unable to add album %#v: %s", album, err)
		test.FailNow()
	}

	// Now, delete the album.
	err = deleteAlbum(album.Id)
	if err != nil {
		test.Errorf("Unable to delete album %#v: %s", album, err)
		test.FailNow()
	}

	// Try to get the album, should fail.
	_, err = getAlbum(album.Id)
	if err == nil {
		test.Errorf("Album %s should have been deleted, it was not.", album.Id)
		test.FailNow()
	}
}

func TestGetAlbum(test *testing.T) {
	albumI := Album{
		Id:       "testGetId",
		Name:     "testGet",
		Price:    "100",
		ArtistId: "testAddAlbumArtist",
	}

	err := addAlbum(&albumI)
	if err != nil {
		test.Errorf("Unable to add album %#v: %s", albumI, err)
		test.FailNow()
	}

	albumF, err := getAlbum(albumI.Id)
	if err != nil {
		test.Errorf("Unable to get album %s", albumI.Id)
		test.FailNow()
	}

	matches := true

	if albumF.Id != albumI.Id {
		test.Errorf("Album id's did not match: %s != %s", albumI.Id, albumF.Id)
		matches = false
	}
	if albumF.Name != albumI.Name {
		test.Errorf("Album id's did not match: %s != %s", albumI.Name, albumF.Name)
		matches = false
	}
	if albumF.Price != albumI.Price {
		test.Errorf("Album id's did not match: %s != %s", albumI.Price, albumF.Price)
		matches = false
	}
	if albumF.ArtistId != albumI.ArtistId {
		test.Errorf("Album id's did not match: %s != %s", albumI.ArtistId, albumF.ArtistId)
		matches = false
	}

	if !matches {
		test.FailNow()
	}
}

func TestGetAllAlbums(test *testing.T) {
	albumI := Album{
		Id:       "testGetAllAlbumsId",
		Name:     "testGetAllAlbums",
		Price:    "testGetAllAlbumsPrice",
		ArtistId: "testGetAllAlbumsArtistId",
	}

	err := addAlbum(&albumI)
	if err != nil {
		test.Errorf("Unable to add album %#v: %s", albumI, err)
		test.FailNow()
	}

	albums, err := getAllAlbums()
	if err != nil {
		test.Errorf("Unable to get album %s: %s", albumI.Id, err)
		test.FailNow()
	}

	if len(albums) <= 0 {
		test.Errorf("GetAllAlbums did not return at least one album.")
		test.FailNow()
	}
}

func TestGetArtistAlbums(test *testing.T) {
	album0 := Album{
		Id:       "testGetArtistAlbumsId0",
		Name:     "testGetArtistAlbums0",
		Price:    "100",
		ArtistId: "testGetArtistAlbumsArtist",
	}
	album1 := Album{
		Id:       "testGetArtistAlbumsId1",
		Name:     "testGetArtistAlbums1",
		Price:    "100",
		ArtistId: "testGetArtistAlbumsArtist",
	}

	// Add the albums first.
	err := addAlbum(&album0)
	if err != nil {
		test.Errorf("Unable to add album %#v: %s", album0, err)
		test.FailNow()
	}
	err = addAlbum(&album1)
	if err != nil {
		test.Errorf("Unable to add album %#v: %s", album1, err)
		test.FailNow()
	}

	// Now list the albums.
	albums, err := getArtistAlbums(album0.ArtistId)
	if err != nil {
		test.Errorf("Unable to get artist albums %s", album0.ArtistId)
		test.FailNow()
	}

	expectedAlbums := []string{album0.Id, album1.Id}

	sort.Strings(expectedAlbums)
	sort.Strings(albums)

	eq := true
	for i := range albums {
		if expectedAlbums[i] != albums[i] {
			eq = false
			break
		}
	}

	if !eq {
		test.Errorf("Albums did not match what was stored. %#v != %#v", albums, expectedAlbums)
		test.FailNow()
	}
}

func TestUpdateAlbum(test *testing.T) {
	album := Album{
		Id:       "testUpdateId",
		Name:     "testUpdate",
		Price:    "100",
		ArtistId: "testAddAlbumArtist",
	}

	// First, add the album.
	err := addAlbum(&album)
	if err != nil {
		test.Errorf("Unable to add album %#v: %s", album, err)
		test.FailNow()
	}

	// Update the album now.
	album.Name = "testUpdate_updated"
	err = updateAlbum(&album)
	if err != nil {
		test.Errorf("Unable to update album %#v: %s", album, err)
		test.FailNow()
	}

	albumF, err := getAlbum(album.Id)
	if err != nil {
		test.Errorf("Unable to get album %s: %s", album.Id, err)
		test.FailNow()
	}

	// Check to make sure the update occurred.
	if album.Name != albumF.Name {
		test.Errorf("Update failed, field 'name' was not updated: %s != %s", album.Name, albumF.Name)
		test.FailNow()
	}
}
