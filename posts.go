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
	req := request{"GET", p, make(chan response)}
	requests <- req
	resp := <-req.response
	if len(resp.posts) == 0 {
		return p, false
	}
	return resp.posts[0], true
}

// List devuelve un slice de todos los posts.
func List() []Post {
	req := request{"LIST", Post{}, make(chan response)}
	requests <- req
	resp := <- req.response
	return resp.posts
}

// Add guarda un post nuevo.
func Add(p Post) {
	req := request{"POST", p, nil} // no hay que esperar respuesta
	requests <- req
}

// Set guarda un post existente.
func Set(p Post) {
	req := request{"PUT", p, nil} // no hay que esperar respuesta
	requests <- req
}

// Del borra un post.
func Del(id int) {
	req := request{"DELETE", Post{Id: id}, nil} // no hay que esperar respuesta
	requests <- req
}

// Solicitud de lectura o mutación al map de posts.
type request struct {
	verb     string
	post     Post
	response chan response
}

// Respuesta a una solicitud de lectura o mutación al map de posts.
type response struct {
	posts []Post
}

// Canal para recibir los requests de lectura o mutación al map de posts.
var requests chan request

// Monitor lanza la goroutine que maneja el map de posts.
func Monitor() chan request {
	// Iniciamos el canal para los requests al monitor.
	requests = make(chan request)

	// Lanzamos la goroutine monitor.
	go func() {
		// El map de posts está confinado a esta goroutine solamente.
		posts := make(map[int]Post)

		// Procesamos los requests.
		for req := range requests {
			var resp response
			// Actuamos según el verbo del request.
			switch req.verb {
			case "LIST":
				for _, p := range posts {
					resp.posts = append(resp.posts, p)
				}
				req.response <- resp
			case "GET":
				if p, ok := posts[req.post.Id]; ok {
					resp.posts = append(resp.posts, p)
				}
				req.response <- resp
			case "POST":
				posts[req.post.Id] = req.post
			case "PUT":
				posts[req.post.Id] = req.post
			case "DELETE":
				delete(posts, req.post.Id)
			}
		}
	}()

	return requests
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

