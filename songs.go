package main

import (
	"errors"
	"sync"
)

type Song struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Genre    string `json:"genre"`
	Time     string `json:"time"`
	Price    string `json:"price"`
	AlbumId  string `json:"albumId"`
	ArtistId string `json:"artistId"`
}

func (song *Song) clone() *Song {
	return &Song{
		Id:       song.Id,
		Name:     song.Name,
		Genre:    song.Genre,
		Time:     song.Time,
		Price:    song.Price,
		AlbumId:  song.AlbumId,
		ArtistId: song.ArtistId,
	}
}

type Songs struct {
	sync.RWMutex
	songs       map[string]*Song
	albumSongs  map[string][]string
	artistSongs map[string][]string
}

func NewSongs() *Songs {
	songs := &Songs{
		songs:       make(map[string]*Song),
		albumSongs:  make(map[string][]string),
		artistSongs: make(map[string][]string),
	}

	return songs
}

func (state *Songs) Add(song *Song) error {
	state.Lock()
	defer state.Unlock()

	// Check if it already exists.
	if _, ok := state.songs[song.Id]; ok {
		return errors.New("Song by 'id' already exists")
	}

	// Add this song to the artist.
	err := state.addArtistSong(song.ArtistId, song.Id)
	if err != nil {
		return err
	}

	// Add this song to the artist.
	err = state.addAlbumSong(song.AlbumId, song.Id)
	if err != nil {
		return err
	}

	// Store the song and release the lock.
	// Copy the struct.
	state.songs[song.Id] = song.clone()

	return nil
}

func (state *Songs) addAlbumSong(albumId, songId string) error {
	songs, ok := state.albumSongs[albumId]
	if ok {
		// Only insert if unique.
		for _, id := range songs {
			if id == songId {
				return errors.New("Song id already exists under that album")
			}
		}

		songs = append(songs, songId)
		state.albumSongs[albumId] = songs
		return nil
	}

	// Does not exist. Create it.
	songs = make([]string, 1)
	songs[0] = songId

	state.albumSongs[albumId] = songs
	return nil
}

func (state *Songs) addArtistSong(artistId, songId string) error {
	songs, ok := state.artistSongs[artistId]
	if ok {
		// Only insert if unique.
		for _, id := range songs {
			if id == songId {
				return errors.New("Song id already exists under that artist")
			}
		}

		songs = append(songs, songId)
		state.artistSongs[artistId] = songs
		return nil
	}

	// Does not exist. Create it.
	songs = make([]string, 1)
	songs[0] = songId

	state.artistSongs[artistId] = songs
	return nil
}

func (state *Songs) Delete(id string) error {
	state.Lock()
	defer state.Unlock()

	if _, ok := state.songs[id]; !ok {
		return errors.New("Song does not exist")
	}

	delete(state.songs, id)

	return nil
}

func (state *Songs) deleteAlbumSong(albumId, songId string) error {
	songs, ok := state.albumSongs[albumId]
	if !ok {
		return errors.New("under that id does not exist")
	}

	// Find the index of the songId.
	index := -1
	for i, id := range songs {
		if id == songId {
			index = i
			break
		}
	}

	if index == -1 {
		return errors.New("Song id not found under that album.")
	}

	// Found the song id, slice the array.
	songs[index] = songs[len(songs)-1]
	songs = songs[:len(songs)-1]

	state.albumSongs[albumId] = songs

	return nil
}

func (state *Songs) deleteArtistSong(artistId, songId string) error {
	songs, ok := state.artistSongs[artistId]
	if !ok {
		return errors.New("under that id does not exist")
	}

	// Find the index of the songId.
	index := -1
	for i, id := range songs {
		if id == songId {
			index = i
			break
		}
	}

	if index == -1 {
		return errors.New("Song id not found under that artist.")
	}

	// Found the song id, slice the array.
	songs[index] = songs[len(songs)-1]
	songs = songs[:len(songs)-1]

	state.artistSongs[artistId] = songs

	return nil
}

func (state *Songs) Get(id string) (*Song, error) {
	state.Lock()
	defer state.Unlock()

	song, ok := state.songs[id]
	if !ok {
		return nil, errors.New("Song does not exist")
	}

	return song, nil
}

func (state *Songs) GetAlbumSongs(albumId string) ([]string, error) {
	state.Lock()
	defer state.Unlock()

	songs, ok := state.albumSongs[albumId]
	if !ok || len(songs) == 0 {
		return nil, errors.New("Album under id does not contain any songs")
	}

	return songs, nil
}

func (state *Songs) GetArtistSongs(artistId string) ([]string, error) {
	state.Lock()
	defer state.Unlock()

	songs, ok := state.artistSongs[artistId]
	if !ok || len(songs) == 0 {
		return nil, errors.New("Artist under id does not contain any songs")
	}

	return songs, nil
}

func (state *Songs) GetAll() ([]string, error) {
	songIds := make([]string, len(state.songs))

	i := -1
	for _, song := range state.songs {
		i++
		songIds[i] = song.Id
	}

	return songIds, nil
}

func (state *Songs) Update(song *Song) error {
	state.Lock()
	defer state.Unlock()

	// Grab the existing song by it's id.
	oldSong, ok := state.songs[song.Id]
	if !ok {
		// Can't update an song that doesn't exist.
		return errors.New("Unable to update song, given song Id does not exist")
	}

	// Only need to update if the albums are different.
	if song.AlbumId != oldSong.AlbumId {
		err := state.deleteAlbumSong(oldSong.AlbumId, oldSong.Id)
		if err != nil {
			return err
		}

		err = state.addAlbumSong(song.AlbumId, song.Id)
		if err != nil {
			return err
		}
	}
	// Only need to update if the artists are different.
	if song.ArtistId != oldSong.ArtistId {
		err := state.deleteArtistSong(oldSong.ArtistId, oldSong.Id)
		if err != nil {
			return err
		}

		err = state.addArtistSong(song.ArtistId, song.Id)
		if err != nil {
			return err
		}
	}

	state.songs[song.Id] = song.clone()

	return nil
}
