## Installation

Requires Docker (tested on 1.9.1)
You may need sudo permissions if on Ubuntu.

> clone <this repo>
> ./build

All methods takes HTTP POST requests with JSON data as the arguments.
See below for specifics on the API.
Example:
  curl -X POST http://localhost:8080/addArtist -d '{"id": "1", "name":"bob","birthdate":"1234"}'
  curl -X POST http://localhost:8080/getArtist -d '"1"'

## Structs

Note: This reflects how the JSON data needs to be structured, not strictly what is stored in memory.

Album = JSON struct of {
  id:       string,
  name:     string,
  price:    string,
  artistId: string
}

Artist = JSON struct of {
  id:        string,
  name:      string,
  birthdate: string
}

Song = JSON struct of {
  id:       string,
  name:     string,
  genre:    string,
  time:     string,
  price:    string,
  albumId:  string,
  artistId: string
}

## Album HTTP API

All methods will either return 200 OK with the data, or a failure and the appropriate error code.

#### /addAlbum: Album -> unit
This method will add a new Album.
This method will also add an association to the 'albumId'.

Takes an Album

Returns no data.

#### /deleteAlbum: string -> unit
This method will delete an existing Album, as well as the link to the Artist.

Takes a string of the Album's id.

Returns no data.

#### /getAlbum: string -> unit
This method will get an existing Album.

Takes a string of the Album's id.

Returns Album.

#### /getArtistAlbums: string -> []string
This method will look up an Artist by it's 'id' and return the list of Albums associated with that Artist.

Takes a string of the Artist's id.

Returns array of Album ids as []string

#### /updateAlbum: Album -> unit
This method will update an existing Album by it's 'id'.
This method will change the Artist association if the artistId is different.

Takes an Album.

Returns no data.

## Artist HTTP API

All methods will either return 200 OK with the data, or a failure and the appropriate error code.

#### /addArtist: Artist -> unit
This method will add a new Artist.

Takes an Artist

Returns no data.

#### /deleteArtist: string -> unit
This method will delete an existing Artist.

Takes a string of the Artist's id.

Returns no data.

#### /getArtist: string -> Artist
This method will get an existing Artist.

Takes a string of the Artist's id.

Returns Artist.

#### /updateArtist: Artist -> unit
This method will update an existing Artist by it's 'id'.

Takes an Artist.

Returns no data.

## Song HTTP API

All methods will either return 200 OK with the data, or a failure and the appropriate error code.

#### /addSong: Song -> unit
This method will add a new Song.
This method will also add an association to the albumId and artistId.

Takes a Song.

Returns no data.

#### /deleteSong: string -> unit
This method will delete an existing Song.
This method will also remove the association to albumId and artistId.

Takes a string of the Song's id.

Returns no data.

#### /getSong: string -> unit
This method will get an existing Song.

Takes a string of the Song's id.

Returns Song.

#### /getAlbumSongs: string -> []string
This method will look up an Album by it's 'id' and return the list of Songs associated with that Album.

Takes a string of the Album's id.

Returns array of Song ids as []string

#### /getArtistSongs: string -> []string
This method will look up an Artist by it's 'id' and return the list of Songs associated with that Artist.

Takes a string of the Artist's id.

Returns array of Song ids as []string

#### /updateSong: Song -> unit
This method will update an existing Song by it's 'id'.
This method will change the Artists/Albums association if the artistId/albumId is different.

Takes a Song.

Returns no data.
