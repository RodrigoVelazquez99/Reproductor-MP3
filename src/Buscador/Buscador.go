package Buscador

import(
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
)

/* Busca del id de un album, sino existe lo crea */
func ObtenerIdAlbum(db *sql.DB, album string) int {
  tabla, err := db.Query("SELECT id_album FROM albums WHERE name =?",album)
  if err != nil { panic(err) }
  var id int
  if tabla.Next() {
    tabla.Scan(&id)
  } else {
    /* El album no existe, entonces se crea */
    stm, _ := db.Prepare("INSERT INTO albums (path, name, year) VALUES (?,?,?)")
    stm.Exec("no_path", album, "2020")
    stm.Close()
    stm1, err1 := db.Query("SELECT id_album FROM albums WHERE name=?", album)
    if err1 != nil { panic(err1) }
    if stm1.Next() {
      stm1.Scan(&id)
    }
    stm1.Close()
  }
  tabla.Close()
  return id
}

func ObtenerIdInterprete(db *sql.DB, interprete string) int {
  tabla, err := db.Query("SELECT id_performer FROM performers WHERE name =?",interprete)
  if err != nil { panic(err) }
  var id int
  for tabla.Next() {
    tabla.Scan(&id)
  }
  tabla.Close()
  return id
}

func BuscaAlbums(db *sql.DB, coincidencia string) []string {
  tabla, err := db.Query("SELECT name FROM albums WHERE name LIKE '" + coincidencia + "%'")
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

func BuscaInterpretes(db *sql.DB, coincidencia string) []string {
  tabla, err := db.Query("SELECT name FROM performers WHERE name LIKE '" + coincidencia + "%'")
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

func BuscaInterprete(id int, db *sql.DB) string {
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

func BuscaAlbum(id int, db *sql.DB) string {
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
