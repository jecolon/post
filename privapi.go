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
	response chan []Post
}

// requests es el canal para recibir los requests para interactuar con los posts.
// start() lo inicia y lanza la monitor goroutine que recibe de él para
// procesar requests. Shutdown() lo cierra para que la monitor goroutine termine.
var requests chan request

// start lanza la monitor goroutine que maneja los posts.
func start() {
	// Iniciamos el canal para los requests a la monitor goroutine.
	requests = make(chan request)
	// Lanzamos la monitor goroutine .
	go monitor()
}

// monitor es la monitor goroutine que procesa cada request recibido por el canal
// requests. Shutdown() la detiene.
func monitor() {
	// El map de posts está confinado a esta goroutine solamente, lo cual
	// sincroniza el acceso y lo hace seguro para uso concurrente.
	posts := make(map[int]Post)

	// Procesamos los requests.
	for req := range requests {
		// Actuamos según el verbo del request.
		switch req.verb {
		case "LIST":
			var list []Post
			for _, p := range posts {
				list = append(list, p)
			}
			// Enviamos la respuesta.
			req.response <- list
		case "GET":
			var list []Post
			if p, ok := posts[req.post.Id]; ok {
				list = append(list, p)
			}
			// Enviamos la respuesta.
			req.response <- list
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
}
