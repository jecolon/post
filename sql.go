package post

import (
	"database/sql"
)

// db es la base de datos global
var db *sql.DB

// Prepared statements
type stmtConfig struct {
	stmt *sql.Stmt
	q    string
}

var prepStmts = map[string]*stmtConfig{
	"get":    {q: "select * from post where id = ?;"},
	"list":   {q: "select * from post;"},
	"insert": {q: "insert into post (id, userId, title, body) values (?, ?, ?, ?);"},
	"update": {q: "update post set userId = ?, title = ?, body = ? where id = ?;"},
	"delete": {q: "delete from post where id = ?;"},
}
