// Package post provides an API for accessing and manipulating content posts.
// The API functions are safe for concurrent access from multiple goroutines.
// Users of this package must call Init() to setup the data store properly.
package post

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	_ "github.com/mattn/go-sqlite3"
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
	posts := getPosts(id)
	if len(posts) == 0 {
		// Slice vacío; no se encontró el post.
		return Post{}, false
	}
	return posts[0], true
}

// List devuelve un slice de todos los posts.
func List() []Post {
	return getPosts(-1)
}

// New guarda un post nuevo.
func New(p Post) []Post {
	return newPost(p)
}

// Put guarda un post existente.
func Put(p Post) {
	putPost(p)
}

// Del borra un post.
func Del(id int) {
	delPost(id)
}

// Init prepara el generador de números aleatorios, la DB y los "prepared statements".
// Si Init devuelve un error, el código cliente no debe usar este paquete.
func Init() error {
	rand.Seed(time.Now().UnixNano())

	// Info para la DB.
	const (
		driver       = "sqlite3"
		dsn          = "posts.db"
		postTableSQL = `create table if not exists post(
			id int primary key not null,
			userId int not null,
			title text not null,
			body text not null
		);`
	)

	// Abrimos la base de datos
	var err error
	db, err = sql.Open(driver, dsn)
	if err != nil {
		return fmt.Errorf("post: error opening DB: %v", err)
	}

	// Creamos la tabla para los post, si no existe
	_, err = db.Exec(postTableSQL)
	if err != nil {
		return fmt.Errorf("post: error creating post table: %v", err)
	}

	// Preparamos los "prepared statements" para get, list, new, put y del.
	for verb, sc := range prepStmts {
		sc.stmt, err = db.Prepare(sc.q)
		if err != nil {
			return fmt.Errorf("post: error preparing %s statement: %v", verb, err)
		}
	}

	return nil
}
