package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

type State struct {
	albums  *Albums
	artists *Artists
	songs   *Songs
}

func NewState() (*State, error) {
	state := &State{
		albums:  NewAlbums(),
		artists: NewArtists(),
		songs:   NewSongs(),
	}

	return state, nil
}

func writeRespError(resp http.ResponseWriter, errResp string) {
	// Set the header.
	resp.Header().Set(
		"Content-Type",
		"application/json;charset=UTF-8",
	)

	resp.WriteHeader(422)
	err := json.NewEncoder(resp).Encode(errResp)
	if err != nil {
		log.Printf("Error writing error response %s: %s", errResp, err)
	}
}

func (state *State) addAlbumHandle(resp http.ResponseWriter, req *http.Request) {
	log.Printf("Got request for addAlbum")

	var album Album
	var err error

	body, err := ioutil.ReadAll(io.LimitReader(req.Body, 1<<10))
	if err != nil {
		log.Printf("Error reading body from %s: %s", req.RemoteAddr, err)
		writeRespError(resp, "Cannot read body from request")
		return
	}

	err = json.Unmarshal(body, &album)
	if err != nil {
		log.Printf("Error deserializing json from %s: %s", req.RemoteAddr, err)
		writeRespError(resp, "Invalid JSON")
		return
	}

	// Try to create the album.
	err = state.albums.Add(&album)
	if err != nil {
		log.Printf("Error adding album %#v for %s: %s", album, req.RemoteAddr, err)
		writeRespError(resp, "Cannot add new album")
		return
	}

	log.Printf("Added album %#v", album)

	resp.WriteHeader(http.StatusOK)
}

func (state *State) addArtistHandle(resp http.ResponseWriter, req *http.Request) {
	log.Printf("Got request for addArtist")

	var artist Artist
	var err error

	body, err := ioutil.ReadAll(io.LimitReader(req.Body, 1<<20))
	if err != nil {
		log.Printf("Error reading body from %s", req.RemoteAddr)
		writeRespError(resp, "Cannot read body from request")
		return
	}

	err = json.Unmarshal(body, &artist)
	if err != nil {
		log.Printf("Error deserializing json from %s: %s", req.RemoteAddr, err)
		writeRespError(resp, "Invalid JSON")
		return
	}

	// Try to create the artist.
	err = state.artists.Add(&artist)
	if err != nil {
		log.Printf("Error storing artist %#v for %s: %s", artist, req.RemoteAddr, err)
		writeRespError(resp, "Unable to store artist")
		return
	}

	log.Printf("Added artist %#v", artist)

	resp.WriteHeader(http.StatusOK)
}

func (state *State) addSongHandle(resp http.ResponseWriter, req *http.Request) {
	log.Printf("Got request for addSong")

	var song Song
	var err error

	body, err := ioutil.ReadAll(io.LimitReader(req.Body, 1<<10))
	if err != nil {
		log.Printf("Error reading body from %s: %s", req.RemoteAddr, err)
		writeRespError(resp, "Cannot read body from request")
		return
	}

	err = json.Unmarshal(body, &song)
	if err != nil {
		log.Printf("Error deserializing json from %s: %s", req.RemoteAddr, err)
		writeRespError(resp, "Invalid JSON")
		return
	}

	// Try to create the song.
	err = state.songs.Add(&song)
	if err != nil {
		log.Printf("Error adding song %#v for %s: %s", song, req.RemoteAddr, err)
		writeRespError(resp, "Cannot add new song")
		return
	}

	log.Printf("Added song %#v", song)

	resp.WriteHeader(http.StatusOK)
}

func (state *State) deleteAlbumHandle(resp http.ResponseWriter, req *http.Request) {
	log.Printf("Got request for deleteAlbum")

	var id string
	var err error

	body, err := ioutil.ReadAll(io.LimitReader(req.Body, 1<<10))
	if err != nil {
		log.Printf("Error reading body from %s: %s", req.RemoteAddr, err)
		writeRespError(resp, "Cannot read body from request")
		return
	}

	err = json.Unmarshal(body, &id)
	if err != nil {
		log.Printf("Error deserializing json from %s: %s", req.RemoteAddr, err)
		writeRespError(resp, "Invalid JSON")
		return
	}

	// Try to delete the album.
	err = state.albums.Delete(id)
	if err != nil {
		log.Printf("Error deleting album %s for %s: %s", id, req.RemoteAddr, err)
		writeRespError(resp, "Unable to delete album")
		return
	}

	resp.WriteHeader(http.StatusOK)
}

func (state *State) deleteArtistHandle(resp http.ResponseWriter, req *http.Request) {
	log.Printf("Got request for deleteArtist")

	var id string
	var err error

	body, err := ioutil.ReadAll(io.LimitReader(req.Body, 1<<10))
	if err != nil {
		log.Printf("Error reading body from %s: %s", req.RemoteAddr, err)
		writeRespError(resp, "Cannot read body from request")
		return
	}

	err = json.Unmarshal(body, &id)
	if err != nil {
		log.Printf("Error deserializing json from %s: %s", req.RemoteAddr, err)
		writeRespError(resp, "Invalid JSON")
		return
	}

	// Try to delete the artist.
	err = state.artists.Delete(id)
	if err != nil {
		log.Printf("Error deleting artist %s for %s: %s", id, req.RemoteAddr, err)
		writeRespError(resp, "Unable to delete artist")
		return
	}

	resp.WriteHeader(http.StatusOK)
}

func (state *State) deleteSongHandle(resp http.ResponseWriter, req *http.Request) {
	log.Printf("Got request for deleteSong")

	var id string
	var err error

	body, err := ioutil.ReadAll(io.LimitReader(req.Body, 1<<10))
	if err != nil {
		log.Printf("Error reading body from %s: %s", req.RemoteAddr, err)
		writeRespError(resp, "Cannot read body from request")
		return
	}

	err = json.Unmarshal(body, &id)
	if err != nil {
		log.Printf("Error deserializing json from %s: %s", req.RemoteAddr, err)
		writeRespError(resp, "Invalid JSON")
		return
	}

	// Try to delete the song.
	err = state.songs.Delete(id)
	if err != nil {
		log.Printf("Error deleting song %s for %s: %s", id, req.RemoteAddr, err)
		writeRespError(resp, "Unable to delete song")
		return
	}

	resp.WriteHeader(http.StatusOK)
}

func (state *State) getAlbumHandle(resp http.ResponseWriter, req *http.Request) {
	log.Printf("Got request for getAlbum")

	var id string
	var err error

	body, err := ioutil.ReadAll(io.LimitReader(req.Body, 1<<10))
	if err != nil {
		log.Printf("Error reading body from %s: %s", req.RemoteAddr, err)
		writeRespError(resp, "Cannot read body from request")
		return
	}

	err = json.Unmarshal(body, &id)
	if err != nil {
		log.Printf("Error deserializing json from %s: %s", req.RemoteAddr, err)
		writeRespError(resp, "Invalid JSON")
		return
	}

	artist, err := state.albums.Get(id)
	if err != nil {
		log.Printf("Error getting artist %s from %s: %s", id, req.RemoteAddr, err)
		writeRespError(resp, "Album does not exist")
		return
	}

	resp.WriteHeader(http.StatusOK)
	err = json.NewEncoder(resp).Encode(*artist)
	if err != nil {
		log.Printf(
			"Error writeing getAlbum response %#v to %s: %s",
			*artist,
			req.RemoteAddr,
			err,
		)
		return
	}
}

func (state *State) getArtistHandle(resp http.ResponseWriter, req *http.Request) {
	log.Printf("Got request for getArtist")

	var id string
	var err error

	body, err := ioutil.ReadAll(io.LimitReader(req.Body, 1<<10))
	if err != nil {
		log.Printf("Error reading body from %s: %s", req.RemoteAddr, err)
		writeRespError(resp, "Cannot read body from request")
		return
	}

	err = json.Unmarshal(body, &id)
	if err != nil {
		log.Printf("Error deserializing json from %s: %s", req.RemoteAddr, err)
		writeRespError(resp, "Invalid JSON")
		return
	}

	artist, err := state.artists.Get(id)
	if err != nil {
		log.Printf("Error getting artist %s from %s: %s", id, req.RemoteAddr, err)
		writeRespError(resp, "Artist does not exist")
		return
	}

	resp.WriteHeader(http.StatusOK)
	err = json.NewEncoder(resp).Encode(*artist)
	if err != nil {
		log.Printf("Error writeing getArtist response %#v to %s: %s", *artist, req.RemoteAddr, err)
	}
}

func (state *State) getSongHandle(resp http.ResponseWriter, req *http.Request) {
	log.Printf("Got request for getSong")

	var id string
	var err error

	body, err := ioutil.ReadAll(io.LimitReader(req.Body, 1<<10))
	if err != nil {
		log.Printf("Error reading body from %s: %s", req.RemoteAddr, err)
		writeRespError(resp, "Cannot read body from request")
		return
	}

	err = json.Unmarshal(body, &id)
	if err != nil {
		log.Printf("Error deserializing json from %s: %s", req.RemoteAddr, err)
		writeRespError(resp, "Invalid JSON")
		return
	}

	song, err := state.songs.Get(id)
	if err != nil {
		log.Printf("Error getting song %s from %s: %s", id, req.RemoteAddr, err)
		writeRespError(resp, "Song does not exist")
		return
	}

	resp.WriteHeader(http.StatusOK)
	err = json.NewEncoder(resp).Encode(*song)
	if err != nil {
		log.Printf("Error writeing getSong response %#v to %s: %s", *song, req.RemoteAddr, err)
	}
}

func (state *State) getArtistAlbumsHandle(resp http.ResponseWriter, req *http.Request) {
	log.Printf("Got request for getArtistAlbums")

	var artistId string
	var err error

	body, err := ioutil.ReadAll(io.LimitReader(req.Body, 1<<10))
	if err != nil {
		log.Printf("Error reading body from %s: %s", req.RemoteAddr, err)
		writeRespError(resp, "Cannot read body from request")
		return
	}

	err = json.Unmarshal(body, &artistId)
	if err != nil {
		log.Printf("Error deserializing json from %s: %s", req.RemoteAddr, err)
		writeRespError(resp, "Invalid JSON")
		return
	}

	albums, err := state.albums.GetArtistAlbums(artistId)
	if err != nil {
		log.Printf("Error retrieving artist %s albums for %s: %s", artistId, req.RemoteAddr, err)
		writeRespError(resp, "Error retrieving artist's albums")
		return
	}

	resp.WriteHeader(http.StatusOK)
	err = json.NewEncoder(resp).Encode(albums)
	if err != nil {
		log.Printf("Error writing getArtistAlbums of artist %s for %s: %s", artistId, req.RemoteAddr, err)
	}
}

func (state *State) getAlbumSongsHandle(resp http.ResponseWriter, req *http.Request) {
	log.Printf("Got request for getAlbumSongs")

	var albumId string
	var err error

	body, err := ioutil.ReadAll(io.LimitReader(req.Body, 1<<10))
	if err != nil {
		log.Printf("Error reading body from %s: %s", req.RemoteAddr, err)
		writeRespError(resp, "Cannot read body from request")
		return
	}

	err = json.Unmarshal(body, &albumId)
	if err != nil {
		log.Printf("Error deserializing json from %s: %s", req.RemoteAddr, err)
		writeRespError(resp, "Invalid JSON")
		return
	}

	songs, err := state.songs.GetAlbumSongs(albumId)
	if err != nil {
		log.Printf("Error retrieving album %s songs for %s: %s", albumId, req.RemoteAddr, err)
		writeRespError(resp, "Error retrieving album's songs")
		return
	}

	resp.WriteHeader(http.StatusOK)
	err = json.NewEncoder(resp).Encode(songs)
	if err != nil {
		log.Printf("Error writing getAlbumSongs of album %s for %s: %s", albumId, req.RemoteAddr, err)
	}
}

func (state *State) getArtistSongsHandle(resp http.ResponseWriter, req *http.Request) {
	log.Printf("Got request for getArtistSongs")

	var artistId string
	var err error

	body, err := ioutil.ReadAll(io.LimitReader(req.Body, 1<<10))
	if err != nil {
		log.Printf("Error reading body from %s: %s", req.RemoteAddr, err)
		writeRespError(resp, "Cannot read body from request")
		return
	}

	err = json.Unmarshal(body, &artistId)
	if err != nil {
		log.Printf("Error deserializing json from %s: %s", req.RemoteAddr, err)
		writeRespError(resp, "Invalid JSON")
		return
	}

	songs, err := state.songs.GetArtistSongs(artistId)
	if err != nil {
		log.Printf("Error retrieving artist %s songs for %s: %s", artistId, req.RemoteAddr, err)
		writeRespError(resp, "Error retrieving artist's songs")
		return
	}

	resp.WriteHeader(http.StatusOK)
	err = json.NewEncoder(resp).Encode(songs)
	if err != nil {
		log.Printf("Error writing getArtistSongs of artist %s for %s: %s", artistId, req.RemoteAddr, err)
	}
}

func (state *State) updateAlbumHandle(resp http.ResponseWriter, req *http.Request) {
	log.Printf("Got request for updateAlbum")

	var album Album
	var err error

	body, err := ioutil.ReadAll(io.LimitReader(req.Body, 1<<10))
	if err != nil {
		log.Printf("Error reading body from %s: %s", req.RemoteAddr, err)
		writeRespError(resp, "Cannot read body from request")
		return
	}

	err = json.Unmarshal(body, &album)
	if err != nil {
		log.Printf("Error deserializing json from %s: %s", req.RemoteAddr, err)
		writeRespError(resp, "Invalid JSON")
		return
	}

	err = state.albums.Update(&album)
	if err != nil {
		log.Printf("Error updating album %#v for %s: %s", album, req.RemoteAddr, err)
		writeRespError(resp, "Unable to update album")
		return
	}

	resp.WriteHeader(http.StatusOK)
}

func (state *State) updateArtistHandle(resp http.ResponseWriter, req *http.Request) {
	log.Printf("Got request for updateArtist")

	var artist Artist
	var err error

	body, err := ioutil.ReadAll(io.LimitReader(req.Body, 1<<10))
	if err != nil {
		log.Printf("Error reading body from %s: %s", req.RemoteAddr, err)
		writeRespError(resp, "Cannot read body from request")
		return
	}

	err = json.Unmarshal(body, &artist)
	if err != nil {
		log.Printf("Error deserializing json from %s: %s", req.RemoteAddr, err)
		writeRespError(resp, "Invalid JSON")
		return
	}

	err = state.artists.Update(&artist)
	if err != nil {
		log.Printf("Error updating artist %#v for %s: %s", artist, req.RemoteAddr, err)
		writeRespError(resp, "Unable to update artist")
		return
	}

	resp.WriteHeader(http.StatusOK)
}

func (state *State) updateSongHandle(resp http.ResponseWriter, req *http.Request) {
	log.Printf("Got request for updateSong")

	var song Song
	var err error

	body, err := ioutil.ReadAll(io.LimitReader(req.Body, 1<<10))
	if err != nil {
		log.Printf("Error reading body from %s: %s", req.RemoteAddr, err)
		writeRespError(resp, "Cannot read body from request")
		return
	}

	err = json.Unmarshal(body, &song)
	if err != nil {
		log.Printf("Error deserializing json from %s: %s", req.RemoteAddr, err)
		writeRespError(resp, "Invalid JSON")
		return
	}

	err = state.songs.Update(&song)
	if err != nil {
		log.Printf("Error updating song %#v for %s: %s", song, req.RemoteAddr, err)
		writeRespError(resp, "Unable to update song")
		return
	}

	resp.WriteHeader(http.StatusOK)
}

func (state *State) notFoundHandle(resp http.ResponseWriter, req *http.Request) {
	log.Printf("Got invalid request url of %s", req.RequestURI)
	resp.WriteHeader(http.StatusNotFound)
}

type tcpListener struct {
	*net.TCPListener
}

func (ln tcpListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}

	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

func init() {
	state, err := NewState()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/addAlbum", state.addAlbumHandle)
	http.HandleFunc("/addArtist", state.addArtistHandle)
	http.HandleFunc("/addSong", state.addSongHandle)

	http.HandleFunc("/deleteAlbum", state.deleteAlbumHandle)
	http.HandleFunc("/deleteArtist", state.deleteArtistHandle)
	http.HandleFunc("/deleteSong", state.deleteSongHandle)

	http.HandleFunc("/getAlbum", state.getAlbumHandle)
	http.HandleFunc("/getArtist", state.getArtistHandle)
	http.HandleFunc("/getSong", state.getSongHandle)

	http.HandleFunc("/getAlbumSongs", state.getAlbumSongsHandle)
	http.HandleFunc("/getArtistAlbums", state.getArtistAlbumsHandle)
	http.HandleFunc("/getArtistSongs", state.getArtistSongsHandle)

	http.HandleFunc("/updateAlbum", state.updateAlbumHandle)
	http.HandleFunc("/updateArtist", state.updateArtistHandle)
	http.HandleFunc("/updateSong", state.updateSongHandle)

	http.HandleFunc("/", state.notFoundHandle)

	log.Printf("Starting http server")

	server := &http.Server{}

	listener, err := net.Listen("tcp", "localhost:3410")
	if err != nil {
		panic(err)
	}

	log.Printf("Started listening.")

	go func() {
		err := server.Serve(tcpListener{listener.(*net.TCPListener)})
		if err != nil {
			panic(err)
		}
	}()
}

func main() {
}
