package main

import(
  "errors"
   "fmt"
   "database/sql"
   _ "github.com/mattn/go-sqlite3"
   "strings"
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
  //solicitud, coincidencia, err1 := buscaCoincidencia(entrada)
  if err1 != nil { }
  err := buscaCancion(base, "oa", "INTERPRETE")
  if err != nil { panic(err) }
}

func IniciaBase() *sql.DB {
  nuevaBase, err := sql.Open("sqlite3", "../Base/base.db")
  if err != nil { panic(err) }
  return nuevaBase
}

func obtenerCoincidencia(palabra string, clave string){
  cadena := ""
  j := 0
  for j < len(tmp) {
    if palabra[j] == clave && j + 1 < len(palabra) {
      if palabra[j+1] == ":" {
        for k := j+1; k < len(palabra); k++ {
          if k == " " || k == "," {
            break
          }
          cadena += palabra[k]
        }
      }
    }
    j++
  }
  return cadena
}

func buscaCoincidencia(entrada string) (string, string, error) {
  var solicitud, coincidencia string
  palabras := strings.Fields(entrada)
  if !strings.Contains(entrada, "T:") || !strings.Contains(entrada, "A:") || !strings.Contains(entrada, "P:") {
    return "","", errors.New("Entrada no valida")
  }
  for int i := 0 ; i < len(palabras) ; i++ {
    if strings.Contains(palabras[i], "T:") {
      obtenerCoincidencia(palabras[i], "T")
    } else if strings.Contains(palabras[i], "A:") {
      obtenerCoincidencia(palabras[i], "A")
    } else if strings.Contains(palabras[i], "P:") {
      obtenerCoincidencia(palabras[i], "P")
    }
  }
  return "","", nil
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
  fmt.Println(title + " " + performer + " " + albums + " " + genre)
}

func buscaCancion(db *sql.DB, coincidencia string, solicitud string) error {
  switch solicitud {
  case "TITULO":
    tabla, err := db.Query("SELECT id_album, id_performer, path, title, genre FROM rolas WHERE title LIKE '%" + coincidencia + "%'")
    if err != nil { panic(err) }
    var title, genre, path string
    var idAlbum, idPerformer int
    for tabla.Next() {
      tabla.Scan(&idAlbum, &idPerformer, &path, &title, &genre)
      performer := buscaInterprete(idPerformer, db)
      Album := buscaAlbum(idAlbum, db)
      creaColumna(title, performer, Album, genre, path)
    }
    tabla.Close()
    return nil
  case "ALBUM" :
    albums := buscaAlbums(db, coincidencia)
    for _,album := range albums {
      idAlbum := obtenerIdAlbum(db, album)
      tabla, err := db.Query("SELECT id_performer, path, title, genre FROM rolas WHERE id_album=?",idAlbum)
      if err != nil { panic(err)}
      var title, genre, path string
      var idPerformer int
      for tabla.Next() {
        tabla.Scan(&idPerformer, &path, &title, &genre)
        performer := buscaInterprete(idPerformer, db)
        creaColumna(title, performer, album, genre, path)
      }
      tabla.Close()
    }
    return nil
  case "INTERPRETE" :
    idInterprete := obtenerIdInterprete(db, coincidencia)
    tabla, err := db.Query("SELECT id_album, path, title, genre FROM rolas WHERE id_performer=?",idInterprete)
    if err != nil { panic(err) }
    var title, genre, path string
    var idAlbum int
    for tabla.Next() {
      tabla.Scan(&idAlbum, &path, &title, &genre)
      performer := buscaInterprete(idInterprete, db)
      album := buscaAlbum(idAlbum, db)
      creaColumna(title, performer, album, genre, path)
    }
    tabla.Close()
    return nil
  }
  return errors.New("fallo")
}

func obtenerIdAlbum(db *sql.DB, album string) int {
  tabla, err := db.Query("SELECT id_album FROM albums WHERE name=?",album)
  if err != nil { panic(err) }
  var id int
  if tabla.Next() {
    tabla.Scan(&id)
  }
  return id
}

func obtenerIdInterprete(db *sql.DB, interprete string) int {
  tabla, err := db.Query("SELECT id_performer FROM performers WHERE name LIKE '%" + interprete + "%'")
  if err != nil { panic(err) }
  var id int
  if tabla.Next() {
    tabla.Scan(&id)
  }
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
  return albums
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
