package post

import(
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

// Monitor lanza la goroutine que maneja los posts.
func Monitor() {
	// Iniciamos el canal para los requests al monitor goroutine.
	requests = make(chan request)

	// Lanzamos la monitor goroutine .
	go func() {
		// El map de posts está confinado a esta goroutine solamente, lo cual
		// sincroniza el acceso y lo hace seguro para uso concurrente.
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
				// Enviamos la respuesta.
				req.response <- resp
			case "GET":
				if p, ok := posts[req.post.Id]; ok {
					resp.posts = append(resp.posts, p)
				}
				// Enviamos la respuesta.
				req.response <- resp
			case "POST":
				posts[req.post.Id] = req.post
			case "PUT":
				// Validación de que el post existe debe ocurrir antes si es necesario.
				posts[req.post.Id] = req.post
			case "DELETE":
				// Validación de que el post existe debe ocurrir antes si es necesario.
				delete(posts, req.post.Id)
			}
		}
	}()
}

func Stop() {
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

