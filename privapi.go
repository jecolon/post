package post

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

