package post

// init lanza la monitor goroutine.
func init() {
	start()
}

// request es enviado al monitor goroutine para interactuar con los posts.
// Una vez enviado, el canal request.response recibirá la respuesta.
type request struct {
	verb     string // LIST, GET, POST, PUT, DELETE
	post     Post
	response chan response
}

// response es la respuesta a un request. Solo los verbos LIST y GET lo usan.
type response struct {
	posts []Post
}

// requests es el canal para recibir los requests para interactuar con los posts.
// Monitor() lo inicia y lanza la monitor goroutine que recibe de él para
// procesar requests. Stop() lo cierra para que la monitor goroutine termine.
var requests chan request

// start lanza la monitor goroutine que maneja los posts.
func start() {
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
