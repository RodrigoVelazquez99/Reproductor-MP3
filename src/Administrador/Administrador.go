package main

import(
  "errors"
   "fmt"
   "database/sql"
   _ "github.com/mattn/go-sqlite3"
   "github.com/RodrigoVelazquez99/Reproductor-MP3/src/Buscador"
   "github.com/RodrigoVelazquez99/Reproductor-MP3/src/Compilador"
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
  entrada := "T:Touch, A:Unknown"
  base := IniciaBase()
  columnas = make([]Columna, 0)
  solicitudes, coincidencias, err := Compilador.BuscaCoincidencias(entrada)
  if err != nil { panic(err) }
  err1 := buscaCancion(base, solicitudes, coincidencias)
  if err1 != nil { panic(err1) }
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
      performer := Buscador.BuscaInterprete(idPerformer, db)
      Album := Buscador.BuscaAlbum(idAlbum, db)
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
      performer := Buscador.BuscaInterprete(idPerformer, db)
      Album := Buscador.BuscaAlbum(idAlbum, db)
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
      performer := Buscador.BuscaInterprete(idPerformer, db)
      Album := Buscador.BuscaAlbum(idAlbum, db)
      if strings.Contains(performer, string(coincidencias[1])) {
        creaColumna(title, performer, Album, genre, path)
      }
    }
    tabla.Close()
    return nil
  } else if banderaAlbum && banderaPerformer {
    albums := Buscador.BuscaAlbums(db, string(coincidencias[1]))
    for _,album := range albums {
      idAlbum := Buscador.ObtenerIdAlbum(db, album)
      tabla, err := db.Query("SELECT id_performer, path, title, genre FROM rolas WHERE id_album=?",idAlbum)
      if err != nil { panic(err)}
      var title, genre, path string
      var idPerformer int
      for tabla.Next() {
        tabla.Scan(&idPerformer, &path, &title, &genre)
        performer := Buscador.BuscaInterprete(idPerformer, db)
        if strings.Contains(performer, coincidencias[0]){
          creaColumna(title, performer, album, genre, path)
        }
      }
      tabla.Close()
    }
    return nil
  } else if banderaAlbum {
    albums := Buscador.BuscaAlbums(db, coincidencias[0])
    for _,album := range albums {
      idAlbum := Buscador.ObtenerIdAlbum(db, album)
      tabla, err := db.Query("SELECT id_performer, path, title, genre FROM rolas WHERE id_album=?",idAlbum)
      if err != nil { panic(err)}
      var title, genre, path string
      var idPerformer int
      for tabla.Next() {
        tabla.Scan(&idPerformer, &path, &title, &genre)
        performer := Buscador.BuscaInterprete(idPerformer, db)
        creaColumna(title, performer, album, genre, path)
      }
      tabla.Close()
    }
    return nil
  } else if banderaPerformer {
    interpretes := Buscador.BuscaInterpretes(db, coincidencias[0])
    for _, interprete := range interpretes {
      idInterprete := Buscador.ObtenerIdInterprete(db, interprete)
      tabla, err := db.Query("SELECT id_album, path, title, genre FROM rolas WHERE id_performer=?",idInterprete)
      if err != nil { panic(err) }
      var title, genre, path string
      var idAlbum int
      for tabla.Next() {
        tabla.Scan(&idAlbum, &path, &title, &genre)
        performer := Buscador.BuscaInterprete(idInterprete, db)
        album := Buscador.BuscaAlbum(idAlbum, db)
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
      performer := Buscador.BuscaInterprete(idPerformer, db)
      Album := Buscador.BuscaAlbum(idAlbum, db)
      creaColumna(title, performer, Album, genre, path)
    }
    tabla.Close()
    return nil
  } else {
    return errors.New("fallo")
  }
}
