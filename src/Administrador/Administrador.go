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
  //solicitudes, coincidencias, err1 := buscaCoincidencias(entrada)
  //if err1 != nil { }
  err := buscaCancion(base, "oa", "INTERPRETE")
  if err != nil { panic(err) }
}

func IniciaBase() *sql.DB {
  nuevaBase, err := sql.Open("sqlite3", "../Base/base.db")
  if err != nil { panic(err) }
  return nuevaBase
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

func buscaCancion(db *sql.DB, banderas []string, coincidencias []string) error {
  var banderaTitle, banderaPerformer, banderaAlbum bool
  for i := 0; i < len(banderas) ; i++ {
    if string(banderas[i]) == "title" {
        banderaTitle = true
        continue
    }
    if string(banderas[i]) == "album" {
        banderaAlbum = true
        continue
    }
    if string(banderas[i]) == "performer" {
        banderaPerformer = true
        continue
    }
  }
  if banderaTitle && banderaAlbum && banderaPerformer {
    tabla, err := db.Query("SELECT id_album, id_performer, path, title, genre FROM rolas WHERE title LIKE '%" + string(coincidencias[0]) + "%'")
    if err != nil { panic(err) }
    var title, genre, path string
    var idAlbum, idPerformer int
    for tabla.Next() {
      tabla.Scan(&idAlbum, &idPerformer, &path, &title, &genre)
      performer := buscaInterprete(idPerformer, db)
      Album := buscaAlbum(idAlbum, db)
      if strings.Contains(performer, string(coincidencias[1])) && strings.Contains(Album, string(coincidencias[2])) {
        creaColumna(title, performer, Album, genre, path)
      }
    }
    tabla.Close()
    return nil
  } else if banderaTitle && banderaAlbum {
    tabla, err := db.Query("SELECT id_album, id_performer, path, title, genre FROM rolas WHERE title LIKE '%" + string(coincidencias[0]) + "%'")
    if err != nil { panic(err) }
    var title, genre, path string
    var idAlbum, idPerformer int
    for tabla.Next() {
      tabla.Scan(&idAlbum, &idPerformer, &path, &title, &genre)
      performer := buscaInterprete(idPerformer, db)
      Album := buscaAlbum(idAlbum, db)
      if strings.Contains(Album, string(coincidencias[1])) {
        creaColumna(title, performer, Album, genre, path)
      }
    }
    tabla.Close()
    return nil
  } else if banderaTitle && banderaPerformer{
    tabla, err := db.Query("SELECT id_album, id_performer, path, title, genre FROM rolas WHERE title LIKE '%" + string(coincidencias[0]) + "%'")
    if err != nil { panic(err) }
    var title, genre, path string
    var idAlbum, idPerformer int
    for tabla.Next() {
      tabla.Scan(&idAlbum, &idPerformer, &path, &title, &genre)
      performer := buscaInterprete(idPerformer, db)
      Album := buscaAlbum(idAlbum, db)
      if strings.Contains(performer, string(coincidencias[1])) {
        creaColumna(title, performer, Album, genre, path)
      }
    }
    tabla.Close()
    return nil
  } else if banderaAlbum && banderaPerformer {
    albums := buscaAlbums(db, strings(coincidencias[1]))
    for _,album := range albums {
      idAlbum := obtenerIdAlbum(db, album)
      tabla, err := db.Query("SELECT id_performer, path, title, genre FROM rolas WHERE id_album=?",idAlbum)
      if err != nil { panic(err)}
      var title, genre, path string
      var idPerformer int
      for tabla.Next() {
        tabla.Scan(&idPerformer, &path, &title, &genre)
        performer := buscaInterprete(idPerformer, db)
        if strins.Contains(performer, coincidencias[0]){
          creaColumna(title, performer, album, genre, path)
        }
      }

      tabla.Close()
    }
    return nil
  } else if banderaAlbum {
    albums := buscaAlbums(db, coincidencias[0])
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
  } else if banderaPerformer {
    interpretes := buscaInterpretes(db, coincidencias[0])
    for _, interprete := range interpretes {
      idInterprete := obtenerIdInterprete(interprete)
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
    }
    return nil
  } else if banderaTitle {
    tabla, err := db.Query("SELECT id_album, id_performer, path, title, genre FROM rolas WHERE title LIKE '%" + string(coincidencias[0]) + "%'")
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
  }
  return errors.New("fallo")
}

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
