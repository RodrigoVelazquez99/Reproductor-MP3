package main

import(
   "fmt"
   "database/sql"
   _ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main()  {
  IniciaBase()
  BuscaCancion("touch")
}

func IniciaBase(){
  nuevaBase, err := sql.Open("sqlite3", "../Base/base.db")
  if err != nil {
    panic(err)
  }
  db = nuevaBase
}

func BuscaCancion(solicitud string) (string, string, string ) {
  rows, err := db.Query("SELECT path, title, genre FROM rolas WHERE title=?","City")
  if err != nil {
    panic(err)
  }
  var nombre, genero, ruta string
  for rows.Next() {
    rows.Scan(&ruta, &nombre, &genero)
    fmt.Println(ruta + " " + nombre + " " + " " + genero)
  }

}

func buscaTitulo(cadena string)  {

}

func buscaInterprete(cadena string)  {

}

func buscaGenero(cadena string) {

}

func buscaAlbum(cadena string)  {

}
