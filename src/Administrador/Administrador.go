package main

import(
   "fmt"
   "database/sql"
   _ "github.com/mattn/go-sqlite3"
)

type Columna struct {
  titulo string
  interprete string
  album string
  genero string
  ruta string
}


func main()  {
  base := IniciaBase()
  BuscaCancion(base, "s")
}

func IniciaBase() *sql.DB {
  nuevaBase, err := sql.Open("sqlite3", "../Base/base.db")
  if err != nil {
    panic(err)
  }
  return nuevaBase
}

func BuscaCancion(db *sql.DB, solicitud string) []Columna {
  columnas:= make([]Columna, 0)
  rows, err := db.Query("SELECT path, title, genre, id_performer, id_album FROM rolas WHERE title LIKE '%" + solicitud + "%'")
  if err != nil {
    panic(err)
  }
  var title, genre, path string
  var idAlbum, idPerformer int
  for rows.Next() {
    rows.Scan(&path, &title, &genre, &idPerformer ,&idAlbum)
    performer := buscaInterprete(idPerformer,db )
    Album := buscaAlbum(idAlbum,db )
    nuevaColumna := Columna{
      titulo: title,
      interprete: performer,
      album: Album,
      genero: genre,
      ruta: path,
    }
    columnas = append(columnas, nuevaColumna)
  }
  rows.Close()
  fmt.Println(columnas)
  return columnas
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
