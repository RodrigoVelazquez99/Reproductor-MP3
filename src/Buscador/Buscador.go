package Buscador

import(
  "errors"
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
  "strings"
)


func obtenerIdAlbum(db *sql.DB, album string) int {
  tabla, err := db.Query("SELECT id_album FROM albums WHERE name =?",album)
  if err != nil { panic(err) }
  var id int
  if tabla.Next() {
    tabla.Scan(&id)
  }
  tabla.Close()
  return id
}

func obtenerIdInterprete(db *sql.DB, interprete string) int {
  tabla, err := db.Query("SELECT id_performer FROM performers WHERE name =?",interprete)
  if err != nil { panic(err) }
  var id int
  for tabla.Next() {
    tabla.Scan(&id)
  }
  tabla.Close()
  return id
}

func buscaAlbums(db *sql.DB, coincidencia string) []string {
  tabla, err := db.Query("SELECT name FROM albums WHERE name LIKE '%" + coincidencia + "%'")
  albums := make([]string, 0)
  if err != nil { panic(err) }
  var name string
  for tabla.Next() {
    tabla.Scan(&name)
    albums = append(albums, name)
  }
  tabla.Close()
  return albums
}

func buscaInterpretes(db *sql.DB, coincidencia string) []string {
  tabla, err := db.Query("SELECT name FROM performers WHERE name LIKE '%" + coincidencia + "%'")
  interpretes := make([]string, 0)
  if err != nil { panic(err) }
  var name string
  for tabla.Next() {
    tabla.Scan(&name)
    interpretes = append(interpretes, name)
  }
  tabla.Close()
  return interpretes
}

func buscaInterprete(id int, db *sql.DB) string {
  tabla,err := db.Query("SELECT name FROM performers WHERE id_performer=?",id)
  if err != nil {
    panic(err)
  }
  var interprete string
  if tabla.Next() {
    tabla.Scan(&interprete)
  }
  tabla.Close()
  return interprete
}

func buscaAlbum(id int, db *sql.DB) string {
  tabla,err := db.Query("SELECT name FROM albums WHERE id_album=?",id)
  if err != nil {
    panic(err)
  }
  var album string
  if tabla.Next() {
    tabla.Scan(&album)
  }
  tabla.Close()
  return album
}
