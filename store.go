package main

import (
  "encoding/json"
  "errors"
  "fmt"
  "io"
  "io/ioutil"
  "log"
	"net/http"
  "sync"
)

type Artist struct {
	Id        string
	Name      string `json:"name"`
  Birthdate string `json:"birthdate"`
}

type Artists struct {
  sync.RWMutex
  m map[string] *Artist
}

type State struct {
  artists *Artists
}

func NewState () (*State, error) {
  artists := &Artists{
    m: make(map[string] *Artist),
  }

  state := &State{
    artists: artists,
  }

  return state, nil
}

func (state *State) addArtist(artist *Artist) error {
  state.artists.Lock()

  if _, ok := state.artists.m[artist.Name]; ok {
    return errors.New("Artist already exists")
  }

  log.Printf("Adding %v", artist)
  state.artists.m[artist.Name] = artist
  state.artists.Unlock()

  return nil
}

func (state *State) checkError(err error, format string, args ...interface{}) bool {
	// Check the error.
	if err == nil {
		return false
	} else {
		// Error is defined, report the error.
    log.Printf(
      format,
      fmt.Sprintf(format, args...),
      err,
    )

    return true
	}
}

func (state *State) addArtistHandle (resp http.ResponseWriter, req *http.Request) {
  var artist Artist

  log.Printf("Reading data from %s", req.RemoteAddr)
  body, err := ioutil.ReadAll(io.LimitReader(req.Body, 1 << 20))

  if state.checkError(err, "Error reading body from %s", req.RemoteAddr) {
    log.Println(err)

    // Set the header.
    resp.Header().Set(
      "Content-Type",
      "application/json;charset=UTF-8",
    )

    // Respond with error.
    resp.WriteHeader(422)
    err = json.NewEncoder(resp).Encode(err)
    state.checkError(err, "Error responding with error to %s", req.RemoteAddr)

    return
  }

  err = json.Unmarshal(body, &artist)
  if state.checkError(err, "Error deserializing json from %s", req.RemoteAddr) {
    // Set the header.
    resp.Header().Set(
      "Content-Type",
      "application/json;charset=UTF-8",
    )

    // Respond with error.
    resp.WriteHeader(422)
    err = json.NewEncoder(resp).Encode(err)
    state.checkError(err, "Error responding with error to %s", req.RemoteAddr)
  } else {
    // Try to create the artist.
    // Do not report errors to the client.

    err = state.addArtist(&artist)
    if state.checkError(err, "creating artist %v", artist) {
      resp.Header().Set(
        "Content-Type",
        "application/json;charset=UTF-8",
      )

      resp.WriteHeader(400)
      //resp.Write(err)
      err = json.NewEncoder(resp).Encode(err.Error())
      state.checkError(err, "responding with error to %s", req.RemoteAddr)
    } else {
      // Respond to the client with on OK status.
      resp.WriteHeader(http.StatusOK)
    }
  }
}

func main () {
  state, err := NewState()
  if err != nil {
    panic(err)
  }

  http.HandleFunc("/addArtist", state.addArtistHandle)

  err = http.ListenAndServe(":8080", nil)
  if err != nil {
    panic(err)
  }
}


