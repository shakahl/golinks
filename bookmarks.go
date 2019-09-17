package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/prologic/bitcask"
)

// Bookmark ...
type Bookmark struct {
	name string
	url  string
}

// Name ...
func (b Bookmark) Name() string {
	return b.name
}

// URL ...
func (b Bookmark) URL() string {
	return b.url
}

// Exec ...
func (b Bookmark) Exec(w http.ResponseWriter, r *http.Request, q string) {
	url := b.url
	if q != "" {
		url = fmt.Sprintf(b.url, q)
	}
	http.Redirect(w, r, url, http.StatusFound)
}

// LookupBookmark ...
func LookupBookmark(name string) (bookmark Bookmark, ok bool) {
	key := fmt.Sprintf("bookmark_%s", strings.ToLower(name))
	val, err := db.Get([]byte(key))
	if err != nil {
		if err == bitcask.ErrKeyNotFound {
			return
		}
		log.Printf("error looking up bookmark for %s: %s", name, err)
	}

	bookmark.name = name
	bookmark.url = string(val)
	ok = true

	return
}
