package main

import (
	"errors"
	"sync"
)

type Artist struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Birthdate string `json:"birthdate"`
}

func (artist *Artist) clone() *Artist {
	return &Artist{
		Id:        artist.Id,
		Name:      artist.Name,
		Birthdate: artist.Birthdate,
	}
}

type Artists struct {
	sync.RWMutex
	artists map[string]*Artist
}

func NewArtists() *Artists {
	artists := &Artists{
		artists: make(map[string]*Artist),
	}

	return artists
}

func (state *Artists) Add(artist *Artist) error {
	state.Lock()
	defer state.Unlock()

	// Check if it already exists.
	if _, ok := state.artists[artist.Id]; ok {
		return errors.New("Artist by 'id' already exists")
	}

	// Store the artist and release the lock.
	// Copy the struct.
	state.artists[artist.Id] = artist.clone()

	return nil
}

func (state *Artists) Delete(id string) error {
	state.Lock()
	defer state.Unlock()

	if _, ok := state.artists[id]; !ok {
		return errors.New("Artist does not exist")
	}

	delete(state.artists, id)

	return nil
}

func (state *Artists) Update(artist *Artist) error {
	state.Lock()
	defer state.Unlock()

	// Grab the existing artist by it's id.
	_, ok := state.artists[artist.Id]
	if !ok {
		// Can't update an artist that doesn't exist.
		return errors.New("Unable to update artist, given artist Id does not exist")
	}

	state.artists[artist.Id] = artist.clone()

	return nil
}

func (state *Artists) Get(id string) (*Artist, error) {
	state.Lock()
	defer state.Unlock()

	artist, ok := state.artists[id]
	if !ok {
		return nil, errors.New("Artist does not exist")
	}

	return artist, nil
}
