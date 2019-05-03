// Package post provides an API for accessing and manipulating content posts.
// The API functions are safe for concurrent access from multiple goroutines.
package post

import (
	"math/rand"
	"time"
)

// Post es una entrada de contenido.
type Post struct {
	Id     int
	UserId int
	Title  string
	Body   string
}

// Get busca un post por ID. El bool es falso si no lo encontramos.
func Get(id int) (Post, bool) {
	p := Post{Id: id}
	req := request{"GET", p, make(chan []Post)}
	requests <- req
	posts := <-req.response
	if len(posts) == 0 {
		return p, false
	}
	return posts[0], true
}

// List devuelve un slice de todos los posts.
func List() []Post {
	req := request{"LIST", Post{}, make(chan []Post)}
	requests <- req
	return <-req.response
}

// Add guarda un post nuevo.
func New(p Post) {
	req := request{"POST", p, nil} // no hay que esperar respuesta
	requests <- req
}

// Set guarda un post existente.
func Put(p Post) {
	req := request{"PUT", p, nil} // no hay que esperar respuesta
	requests <- req
}

// Del borra un post.
func Del(id int) {
	req := request{"DELETE", Post{Id: id}, nil} // no hay que esperar respuesta
	requests <- req
}

// Shutdown detiene la monitor goroutine. Cualquier uso de las funciones de este
// paquete luego de llamar Shutdown resultarán en un panic.
func Shutdown() {
	// Al cerrar requests, la monitor goroutine terminará.
	close(requests)
}

// NewId genera un ID único para un post.
func NewId() int {
	rand.Seed(time.Now().UnixNano())
	id := rand.Intn(1000)
	for {
		if _, ok := Get(id); !ok {
			break
		}
		id = rand.Intn(1000)
	}
	return id
}
