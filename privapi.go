package post

import (
	"database/sql"
	"log"
	"math/rand"
)

// getPosts busca un post con id o listado de todos si id es -1.
func getPosts(id int) []Post {
	res := []Post{}
	if id != -1 {
		var p Post
		// Obtenemos y ejecutamos el get prepared statement.
		get := prepStmts["get"].stmt
		err := get.QueryRow(id).Scan(&p.Id, &p.UserId, &p.Title, &p.Body)
		if err != nil {
			if err != sql.ErrNoRows {
				log.Printf("post: error getting post. Id: %d, err: %v\n", id, err)
			}
		} else {
			res = append(res, p)
		}
		return res
	}

	// Obtenemos y ejecutamos el list prepared statement.
	list := prepStmts["list"].stmt
	rows, err := list.Query()
	if err != nil {
		log.Printf("post: error getting posts. err: %v\n", err)
	}
	defer rows.Close()

	// Procesamos los rows.
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.Id, &p.UserId, &p.Title, &p.Body); err != nil {
			log.Printf("post: error scanning row: %v\n", err)
			continue
		}
		res = append(res, p)
	}
	// Verificamos si hubo error procesando los rows.
	if err := rows.Err(); err != nil {
		log.Printf("post: error reading rows: %v\n", err)
	}

	return res
}

// newPost inserta un post en la DB.
func newPost(p Post) []Post {
	// Generamos ID único para el nuevo post.
	p.Id = rand.Intn(1000)
	for {
		l := getPosts(p.Id)
		if len(l) == 0 {
			break
		}
		p.Id = rand.Intn(1000)
	}

	// Obtenemos y ejecutamos insert prepared statement.
	insert := prepStmts["insert"].stmt
	_, err := insert.Exec(p.Id, p.UserId, p.Title, p.Body)
	if err != nil {
		log.Printf("post: error inserting post %d into DB: %v\n", p.Id, err)
	}
	return []Post{p}
}

// putPost actualiza un post en la DB.
func putPost(p Post) {
	// Obtenemos y ejecutamos update prepared statement.
	update := prepStmts["update"].stmt
	_, err := update.Exec(p.UserId, p.Title, p.Body, p.Id)
	if err != nil {
		log.Printf("post: error updating post %d into DB: %v\n", p.Id, err)
	}
}

// delPost borra un post de la DB.
func delPost(id int) {
	// Obtenemos y ejecutamos delete prepared statement.
	del := prepStmts["delete"].stmt
	_, err := del.Exec(id)
	if err != nil {
		log.Printf("post: error deleting post %d into DB: %v\n", id, err)
	}
}
