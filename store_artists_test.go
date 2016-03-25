package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"
)

func AddArtist(artist *Artist) error {
	buffer, err := json.Marshal(artist)
	if err != nil {
		return err
	}

	resp, err := http.Post(
		TEST_SERVER_END_POINT+"addArtist",
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

func deleteArtist(id string) error {
	buffer, err := json.Marshal(id)
	if err != nil {
		return err
	}

	resp, err := http.Post(
		TEST_SERVER_END_POINT+"deleteArtist",
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

func getArtist(id string) (*Artist, error) {
	buffer, err := json.Marshal(id)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(
		TEST_SERVER_END_POINT+"getArtist",
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

	artist := new(Artist)
	err = json.Unmarshal(body, artist)
	if err != nil {
		return nil, err
	}

	return artist, nil
}

func getAllArtists() ([]string, error) {
	resp, err := http.Post(
		TEST_SERVER_END_POINT+"getAllArtists",
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

	artists := make([]string, 0)
	err = json.Unmarshal(body, &artists)
	if err != nil {
		return nil, err
	}

	return artists, nil
}

func updateArtist(artist *Artist) error {
	buffer, err := json.Marshal(artist)
	if err != nil {
		return err
	}

	resp, err := http.Post(
		TEST_SERVER_END_POINT+"updateArtist",
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

func TestAddArtist(test *testing.T) {
	artist := Artist{
		Id:        "testAddId",
		Name:      "testAdd",
		Birthdate: "1234",
	}

	test.Log("Adding artist")
	err := AddArtist(&artist)
	if err != nil {
		test.Errorf("Unable to add artist %#v: %s", artist, err)
		test.FailNow()
	}
}

func TestDeleteArtist(test *testing.T) {
	// First, add the artist.
	artist := Artist{
		Id:        "testDeleteArtistId",
		Name:      "testDelete",
		Birthdate: "1234",
	}

	err := AddArtist(&artist)
	if err != nil {
		test.Errorf("Unable to add artist %#v: %s", artist, err)
		test.FailNow()
	}

	// Now, delete the artist.
	err = deleteArtist(artist.Id)
	if err != nil {
		test.Errorf("Unable to delete artist %#v: %s", artist, err)
		test.FailNow()
	}

	// Try to get the artist, should fail.
	_, err = getArtist(artist.Id)
	if err == nil {
		test.Errorf("Artist %s should have been deleted, it was not.", artist.Id)
		test.FailNow()
	}
}

func TestGetArtist(test *testing.T) {
	artistI := Artist{
		Id:        "testGetId",
		Name:      "testGet",
		Birthdate: "1234",
	}

	err := AddArtist(&artistI)
	if err != nil {
		test.Errorf("Unable to add artist %#v: %s", artistI, err)
		test.FailNow()
	}

	artistF, err := getArtist(artistI.Id)
	if err != nil {
		test.Errorf("Unable to get artist %s", artistI.Id)
		test.FailNow()
	}

	matches := true

	if artistF.Id != artistI.Id {
		test.Errorf("Artist id's did not match: %s != %s", artistI.Id, artistF.Id)
		matches = false
	}
	if artistF.Name != artistI.Name {
		test.Errorf("Artist id's did not match: %s != %s", artistI.Name, artistF.Name)
		matches = false
	}
	if artistF.Birthdate != artistI.Birthdate {
		test.Errorf("Artist id's did not match: %s != %s", artistI.Birthdate, artistF.Birthdate)
		matches = false
	}

	if !matches {
		test.FailNow()
	}
}

func TestGetAllArtists(test *testing.T) {
	artistI := Artist{
		Id:        "testGetAllArtistsId",
		Name:      "testGetAllArtists",
		Birthdate: "testGetAllArtistsBirthdate",
	}

	err := AddArtist(&artistI)
	if err != nil {
		test.Errorf("Unable to add artist %#v: %s", artistI, err)
		test.FailNow()
	}

	artists, err := getAllArtists()
	if err != nil {
		test.Errorf("Unable to get artist %s: %s", artistI.Id, err)
		test.FailNow()
	}

	if len(artists) <= 0 {
		test.Errorf("GetAllArtists did not return at least one artist.")
		test.FailNow()
	}
}

func TestUpdateArtist(test *testing.T) {
	artist := Artist{
		Id:        "testUpdateId",
		Name:      "testUpdate",
		Birthdate: "1234",
	}

	// First, add the artist.
	err := AddArtist(&artist)
	if err != nil {
		test.Errorf("Unable to add artist %#v: %s", artist, err)
		test.FailNow()
	}

	// Update the artist now.
	artist.Name = "testUpdate_updated"
	err = updateArtist(&artist)
	if err != nil {
		test.Errorf("Unable to update artist %#v: %s", artist, err)
		test.FailNow()
	}

	artistF, err := getArtist(artist.Id)
	if err != nil {
		test.Errorf("Unable to get artist %s: %s", artist.Id, err)
		test.FailNow()
	}

	// Check to make sure the update occurred.
	if artist.Name != artistF.Name {
		test.Errorf("Update failed, field 'name' was not updated: %s != %s", artist.Name, artistF.Name)
		test.FailNow()
	}
}
