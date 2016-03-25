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

func addSong(song *Song) error {
	buffer, err := json.Marshal(song)
	if err != nil {
		return err
	}

	resp, err := http.Post(
		TEST_SERVER_END_POINT+"addSong",
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

func deleteSong(id string) error {
	buffer, err := json.Marshal(id)
	if err != nil {
		return err
	}

	resp, err := http.Post(
		TEST_SERVER_END_POINT+"deleteSong",
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

func getSong(id string) (*Song, error) {
	buffer, err := json.Marshal(id)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(
		TEST_SERVER_END_POINT+"getSong",
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

	song := new(Song)
	err = json.Unmarshal(body, song)
	if err != nil {
		return nil, err
	}

	return song, nil
}

func getAlbumSongs(albumId string) ([]string, error) {
	buffer, err := json.Marshal(albumId)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(
		TEST_SERVER_END_POINT+"getAlbumSongs",
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

	songs := make([]string, 0)
	err = json.Unmarshal(body, &songs)
	if err != nil {
		return nil, err
	}

	return songs, nil
}

func getArtistSongs(artistId string) ([]string, error) {
	buffer, err := json.Marshal(artistId)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(
		TEST_SERVER_END_POINT+"getArtistSongs",
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

	songs := make([]string, 0)
	err = json.Unmarshal(body, &songs)
	if err != nil {
		return nil, err
	}

	return songs, nil
}

func updateSong(song *Song) error {
	buffer, err := json.Marshal(song)
	if err != nil {
		return err
	}

	resp, err := http.Post(
		TEST_SERVER_END_POINT+"updateSong",
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

func TestAddSong(test *testing.T) {
	song := Song{
		Id:       "testAddId",
		Name:     "testAdd",
		Genre:    "testAddGenre",
		Time:     "testAddSongTime",
		Price:    "testAddSongPrice",
		AlbumId:  "testAddSongAlbumId",
		ArtistId: "testAddSongArtistId",
	}

	test.Log("Adding song")
	err := addSong(&song)
	if err != nil {
		test.Errorf("Unable to add song %#v: %s", song, err)
		test.FailNow()
	}
}

func TestDeleteSong(test *testing.T) {
	// First, add the song.
	song := Song{
		Id:       "testDeleteId",
		Name:     "testDelete",
		Genre:    "testDeleteGenre",
		Time:     "testDeleteSongTime",
		Price:    "testDeleteSongPrice",
		AlbumId:  "testDeleteSongAlbumId",
		ArtistId: "testDeleteSongArtistId",
	}

	err := addSong(&song)
	if err != nil {
		test.Errorf("Unable to add song %#v: %s", song, err)
		test.FailNow()
	}

	// Now, delete the song.
	err = deleteSong(song.Id)
	if err != nil {
		test.Errorf("Unable to delete song %#v: %s", song, err)
		test.FailNow()
	}

	// Try to get the song, should fail.
	_, err = getSong(song.Id)
	if err == nil {
		test.Errorf("Song %s should have been deleted, it was not.", song.Id)
		test.FailNow()
	}
}

func TestGetSong(test *testing.T) {
	songI := Song{
		Id:       "testGetId",
		Name:     "testGet",
		Genre:    "testGetGenre",
		Time:     "testGetSongTime",
		Price:    "testGetSongPrice",
		AlbumId:  "testGetSongAlbumId",
		ArtistId: "testGetSongArtistId",
	}

	err := addSong(&songI)
	if err != nil {
		test.Errorf("Unable to add song %#v: %s", songI, err)
		test.FailNow()
	}

	songF, err := getSong(songI.Id)
	if err != nil {
		test.Errorf("Unable to get song %s", songI.Id)
		test.FailNow()
	}

	matches := true

	if songF.Id != songI.Id {
		test.Errorf("Song id's did not match: %s != %s", songI.Id, songF.Id)
		matches = false
	}
	if songF.Name != songI.Name {
		test.Errorf("Song id's did not match: %s != %s", songI.Name, songF.Name)
		matches = false
	}
	if songF.Price != songI.Price {
		test.Errorf("Song id's did not match: %s != %s", songI.Price, songF.Price)
		matches = false
	}
	if songF.ArtistId != songI.ArtistId {
		test.Errorf("Song id's did not match: %s != %s", songI.ArtistId, songF.ArtistId)
		matches = false
	}

	if !matches {
		test.FailNow()
	}
}

func TestGetAlbumSongs(test *testing.T) {
	song0 := Song{
		Id:       "testGetAlbumSongsId0",
		Name:     "testGetAlbumSongs0",
		Genre:    "testGetAlbumSongsGenre0",
		Time:     "testGetAlbumSongsTime0",
		Price:    "testGetAlbumSongsPrice0",
		AlbumId:  "testGetAlbumSongsAlbumId",
		ArtistId: "testGetAlbumSongsArtistId",
	}
	song1 := Song{
		Id:       "testGetAlbumSongsId1",
		Name:     "testGetAlbumSongs1",
		Genre:    "testGetAlbumSongsGenre1",
		Time:     "testGetAlbumSongsTime1",
		Price:    "testGetAlbumSongsPrice1",
		AlbumId:  "testGetAlbumSongsAlbumId",
		ArtistId: "testGetAlbumSongsArtistId",
	}

	// Add the songs first.
	err := addSong(&song0)
	if err != nil {
		test.Errorf("Unable to add song %#v: %s", song0, err)
		test.FailNow()
	}
	err = addSong(&song1)
	if err != nil {
		test.Errorf("Unable to add song %#v: %s", song1, err)
		test.FailNow()
	}

	// Now list the songs.
	songs, err := getAlbumSongs(song0.AlbumId)
	if err != nil {
		test.Errorf("Unable to get album songs %s", song0.ArtistId)
		test.FailNow()
	}

	expectedSongs := []string{song0.Id, song1.Id}

	sort.Strings(expectedSongs)
	sort.Strings(songs)

	eq := true
	for i := range songs {
		if expectedSongs[i] != songs[i] {
			eq = false
			break
		}
	}

	if !eq {
		test.Errorf("Songs did not match what was stored. %#v != %#v", songs, expectedSongs)
		test.FailNow()
	}
}

func TestGetArtistSongs(test *testing.T) {
	song0 := Song{
		Id:       "testGetArtistSongsId0",
		Name:     "testGetArtistSongs0",
		Genre:    "testGetArtistSongsGenre0",
		Time:     "testGetArtistSongsTime0",
		Price:    "testGetArtistSongsPrice0",
		AlbumId:  "testGetArtistSongsAlbumId",
		ArtistId: "testGetArtistSongsArtistId",
	}
	song1 := Song{
		Id:       "testGetArtistSongsId1",
		Name:     "testGetArtistSongs1",
		Genre:    "testGetArtistSongsGenre1",
		Time:     "testGetArtistSongsTime1",
		Price:    "testGetArtistSongsPrice1",
		AlbumId:  "testGetArtistSongsAlbumId",
		ArtistId: "testGetArtistSongsArtistId",
	}

	// Add the songs first.
	err := addSong(&song0)
	if err != nil {
		test.Errorf("Unable to add song %#v: %s", song0, err)
		test.FailNow()
	}
	err = addSong(&song1)
	if err != nil {
		test.Errorf("Unable to add song %#v: %s", song1, err)
		test.FailNow()
	}

	// Now list the songs.
	songs, err := getArtistSongs(song0.ArtistId)
	if err != nil {
		test.Errorf("Unable to get artist songs %s", song0.ArtistId)
		test.FailNow()
	}

	expectedSongs := []string{song0.Id, song1.Id}

	sort.Strings(expectedSongs)
	sort.Strings(songs)

	eq := true
	for i := range songs {
		if expectedSongs[i] != songs[i] {
			eq = false
			break
		}
	}

	if !eq {
		test.Errorf("Songs did not match what was stored. %#v != %#v", songs, expectedSongs)
		test.FailNow()
	}
}

func TestUpdateSong(test *testing.T) {
	song := Song{
		Id:       "testUpdateId",
		Name:     "testUpdate",
		Genre:    "testUpdateGenre",
		Time:     "testUpdateSongTime",
		Price:    "testUpdateSongPrice",
		AlbumId:  "testUpdateSongAlbumId",
		ArtistId: "testUpdateSongArtistId",
	}

	// First, add the song.
	err := addSong(&song)
	if err != nil {
		test.Errorf("Unable to add song %#v: %s", song, err)
		test.FailNow()
	}

	// Update the song now.
	song.Name = "testUpdate_updated"
	err = updateSong(&song)
	if err != nil {
		test.Errorf("Unable to update song %#v: %s", song, err)
		test.FailNow()
	}

	songF, err := getSong(song.Id)
	if err != nil {
		test.Errorf("Unable to get song %s: %s", song.Id, err)
		test.FailNow()
	}

	// Check to make sure the update occurred.
	if song.Name != songF.Name {
		test.Errorf("Update failed, field 'name' was not updated: %s != %s", song.Name, songF.Name)
		test.FailNow()
	}
}
