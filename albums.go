package main

import (
	"errors"
	"sync"
)

type Album struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Price    string `json:"price"`
	ArtistId string `json:"albumId"`
}

func (album *Album) clone() *Album {
	return &Album{
		Id:       album.Id,
		Name:     album.Name,
		Price:    album.Price,
		ArtistId: album.ArtistId,
	}
}

type Albums struct {
	sync.RWMutex
	albums       map[string]*Album
	artistAlbums map[string][]string
}

func NewAlbums() *Albums {
	albums := &Albums{
		albums:       make(map[string]*Album),
		artistAlbums: make(map[string][]string),
	}

	return albums
}

func (state *Albums) Add(album *Album) error {
	state.Lock()
	defer state.Unlock()

	// Check if it already exists.
	if _, ok := state.albums[album.Id]; ok {
		return errors.New("Artist by 'id' already exists")
	}

	// Add this album to the artist.
	err := state.addArtistAlbum(album.ArtistId, album.Id)
	if err != nil {
		return err
	}

	// Copy the struct.
	state.albums[album.Id] = album.clone()

	return nil
}

func (state *Albums) addArtistAlbum(artistId, albumId string) error {
	albums, ok := state.artistAlbums[artistId]
	if ok {
		// Only insert if unique.
		for _, id := range albums {
			if id == albumId {
				return errors.New("Album id already exists under that artist")
			}
		}

		albums = append(albums, albumId)
		state.artistAlbums[artistId] = albums
		return nil
	}

	// Does not exist. Create it.
	albums = make([]string, 1)
	albums[0] = albumId

	state.artistAlbums[artistId] = albums
	return nil
}

func (state *Albums) Delete(id string) error {
	state.Lock()
	defer state.Unlock()

	album, ok := state.albums[id]
	if !ok {
		return errors.New("Album does not exist")
	}

	// Get the artist of this album, and make sure it's removed from the list.
	err := state.deleteArtistAlbum(album.ArtistId, id)
	if err != nil {
		return err
	}

	delete(state.albums, id)

	return nil
}

func (state *Albums) deleteArtistAlbum(artistId, albumId string) error {
	albums, ok := state.artistAlbums[artistId]
	if !ok {
		return errors.New("Artist under that id does not exist")
	}

	// Find the index of the albumId.
	index := -1
	for i, id := range albums {
		if id == albumId {
			index = i
			break
		}
	}

	if index == -1 {
		return errors.New("Album id not found under that artist.")
	}

	// Found the album id, slice the array.
	albums[index] = albums[len(albums)-1]
	albums = albums[:len(albums)-1]

	state.artistAlbums[artistId] = albums

	return nil
}

func (state *Albums) Get(id string) (*Album, error) {
	state.Lock()
	defer state.Unlock()

	album, ok := state.albums[id]
	if !ok {
		return nil, errors.New("Album does not exist")
	}

	return album, nil
}

func (state *Albums) GetArtistAlbums(artistId string) ([]string, error) {
	state.Lock()
	defer state.Unlock()

	albums, ok := state.artistAlbums[artistId]
	if !ok || len(albums) == 0 {
		return nil, errors.New("Artist under id does not contain any albums")
	}

	return albums, nil
}

func (state *Albums) Update(album *Album) error {
	state.Lock()
	defer state.Unlock()

	// Grab the existing album by it's id.
	oldAlbum, ok := state.albums[album.Id]
	if !ok {
		// Can't update an album that doesn't exist.
		return errors.New("Unable to update album, given album Id does not exist")
	}

	// Only need to update if the artists are different.
	if album.ArtistId != oldAlbum.ArtistId {
		err := state.deleteArtistAlbum(oldAlbum.ArtistId, oldAlbum.Id)
		if err != nil {
			return err
		}

		err = state.addArtistAlbum(album.ArtistId, album.Id)
		if err != nil {
			return err
		}
	}

	state.albums[album.Id] = album.clone()

	return nil
}
