package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/prologic/bitcask"
	"github.com/stretchr/testify/assert"
)

type Explode struct{}

func (e Explode) Name() string {
	return "explode"
}

func (e Explode) Desc() string {
	return `explode

	This command just quite literally explodes and always returns an error!
	`
}

func (e Explode) Exec(w http.ResponseWriter, r *http.Request, args []string) error {
	return errors.New("kaboom")
}

func TestRender(t *testing.T) {
	assert := assert.New(t)

	s, err := NewServer(":8000", Config{})
	assert.NoError(err)

	w := httptest.NewRecorder()

	s.render("index", w, nil)

	assert.Equal(w.Code, http.StatusOK)
	assert.Contains(w.Body.String(), `name="q"`)
}

func TestRenderError(t *testing.T) {
	assert := assert.New(t)

	s, err := NewServer(":8000", Config{})
	assert.NoError(err)

	w := httptest.NewRecorder()

	s.render("asdf", w, nil)

	assert.Equal(w.Code, http.StatusInternalServerError)
}

func TestIndex(t *testing.T) {
	assert := assert.New(t)

	s, err := NewServer(":8000", Config{})
	assert.NoError(err)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/", nil)
	p := httprouter.Params{}

	s.IndexHandler()(w, r, p)
	assert.Equal(w.Code, http.StatusOK)
	assert.Contains(w.Body.String(), `name="q"`)
}

func TestOpenSearch(t *testing.T) {
	assert := assert.New(t)

	s, err := NewServer(":8000", Config{})
	assert.NoError(err)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/opensearch.xml", nil)
	p := httprouter.Params{}

	s.OpenSearchHandler()(w, r, p)
	assert.Equal(w.Code, http.StatusOK)
	assert.Contains(w.Body.String(), "<OpenSearchDescription")
}

func TestCommand(t *testing.T) {
	assert := assert.New(t)

	s, err := NewServer(":8000", Config{})
	assert.NoError(err)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/?q=ping", nil)
	p := httprouter.Params{}

	s.IndexHandler()(w, r, p)
	assert.Equal(w.Code, http.StatusOK)

	body := w.Body.String()
	tokens := strings.Split(body, " ")
	assert.Len(tokens, 2)
	assert.Equal(tokens[0], "pong")

	n, err := strconv.ParseInt(tokens[1], 10, 64)
	assert.Nil(err)
	assert.WithinDuration(time.Now(), time.Unix(n, 0), 1*time.Second)
}

func TestCaseInsensitive(t *testing.T) {
	assert := assert.New(t)

	s, err := NewServer(":8000", Config{})
	assert.NoError(err)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/?q=Ping", nil)
	p := httprouter.Params{}

	s.IndexHandler()(w, r, p)
	assert.Equal(w.Code, http.StatusOK)

	body := w.Body.String()
	tokens := strings.Split(body, " ")
	assert.Len(tokens, 2)
	assert.Equal(tokens[0], "pong")

	n, err := strconv.ParseInt(tokens[1], 10, 64)
	assert.Nil(err)
	assert.WithinDuration(time.Now(), time.Unix(n, 0), 1*time.Second)
}

func TestInvalidCommand(t *testing.T) {
	assert := assert.New(t)

	db, _ = bitcask.Open("test.db")
	defer db.Close()

	s, err := NewServer(":8000", Config{URL: ""})
	assert.NoError(err)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/?q=asdf", nil)
	p := httprouter.Params{}

	s.IndexHandler()(w, r, p)
	assert.Equal(w.Code, http.StatusBadRequest)

	body := w.Body.String()
	assert.Equal(body, "Invalid Command: asdf\n")
}

func TestInvalidCommandDefaultURL(t *testing.T) {
	assert := assert.New(t)

	db, _ = bitcask.Open("test.db")
	defer db.Close()

	s, err := NewServer(":8000", Config{URL: DefaultURL})
	assert.NoError(err)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/?q=asdf", nil)
	p := httprouter.Params{}

	s.IndexHandler()(w, r, p)
	assert.Condition(func() bool {
		return w.Code >= http.StatusMultipleChoices &&
			w.Code <= http.StatusTemporaryRedirect
	})

	assert.Equal(
		w.Header().Get("Location"),
		"https://www.google.com/search?q=asdf&btnK",
	)
}

func TestCommandError(t *testing.T) {
	assert := assert.New(t)

	db, _ = bitcask.Open("test.db")
	defer db.Close()

	RegisterCommand("explode", Explode{})

	s, err := NewServer(":8000", Config{})
	assert.NoError(err)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/?q=explode", nil)
	p := httprouter.Params{}

	s.IndexHandler()(w, r, p)
	assert.Equal(w.Code, http.StatusInternalServerError)

	body := w.Body.String()
	assert.Equal(body, "Error processing command explode: kaboom\n")
}

func TestCommandBookmark(t *testing.T) {
	assert := assert.New(t)

	db, _ = bitcask.Open("test.db")
	defer db.Close()

	err := EnsureDefaultBookmarks()
	assert.Nil(err)

	s, err := NewServer(":8000", Config{})
	assert.NoError(err)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/?q=g%20foo%20bar", nil)
	p := httprouter.Params{}

	s.IndexHandler()(w, r, p)
	assert.Condition(func() bool {
		return w.Code >= http.StatusMultipleChoices &&
			w.Code <= http.StatusTemporaryRedirect
	})

	assert.Equal(
		w.Header().Get("Location"),
		"https://www.google.com/search?q=foo bar&btnK",
	)
}
