package main

import(
   "fmt"
   "database/sql"
   _ "github.com/mattn/go-sqlite3"
)

var columnas []Columna

type Columna struct {
  titulo string
  interprete string
  album string
  genero string
  ruta string
}

func main()  {
  base := IniciaBase()
  columnas = make([]Columna, 0)
  BuscaCancion(base, "s")
}

func IniciaBase() *sql.DB {
  nuevaBase, err := sql.Open("sqlite3", "../Base/base.db")
  if err != nil {
    panic(err)
  }
  return nuevaBase
}

func BuscaCancion(db *sql.DB, solicitud string) {
  buscaPorTitulo(db, solicitud)
  buscaPorInterprete(db, solicitud)
  buscaPorAlbum(db,solicitud)
  buscaPorGenero(db, solicitud)
}

func creaColumna(title string, performer string, albums string, genre string, path string )  {
  nuevaColumna := Columna{
    titulo: title,
    interprete: performer,
    album: albums,
    genero: genre,
    ruta: path,
  }
  columnas = append(columnas, nuevaColumna)
  fmt.Println(nuevaColumna)
}

func buscaPorTitulo(db *sql.DB, solicitud string){
  tabla, err := db.Query("SELECT path, title, genre, id_performer, id_album FROM rolas WHERE title LIKE '%" + solicitud + "%'")
  if err != nil {
    panic(err)
  }
  var title, genre, path string
  var idAlbum, idPerformer int
  for tabla.Next() {
    tabla.Scan(&path, &title, &genre, &idPerformer ,&idAlbum)
    performer := buscaInterprete(idPerformer,db )
    Album := buscaAlbum(idAlbum,db )
    creaColumna(title, performer, Album, genre, path)
  }
  tabla.Close()
}
func buscaPorAlbum(db *sql.DB, solicitud string){
  tabla, err := db.Query("SELECT path, name FROM albums WHERE name LIKE '%" + solicitud + "%'")
  if err != nil {
    panic(err)
  }
  var title, genre, path string
  var idAlbum, idPerformer int
  for tabla.Next() {
    tabla.Scan(&path, &title, &genre, &idPerformer ,&idAlbum)
    performer := buscaInterprete(idPerformer,db )
    Album := buscaAlbum(idAlbum,db )
    creaColumna(title, performer, Album, genre, path)
  }
  tabla.Close()
}

func buscaPorGenero(db *sql.DB, solicitud string) {

}
func buscaPorInterprete(db *sql.DB, solicitud string){

}
func buscaColumna(columna Columna){

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
