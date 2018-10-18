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

func obtenerCoincidencia(palabra string, clave string) string {
  cadena := ""
  j := 0
  for j < len(palabra) {
    if string(palabra[j]) == clave && j + 1 < len(palabra) {
      if string(palabra[j+1]) == ":" {
        for k := j+2; k < len(palabra); k++ {
          if string(palabra[k]) == "," {
            break
          }
          cadena += string(palabra[k])
        }
      }
    }
    j++
  }
  if cadena == "" { return "", errors.New("Sin cadena") }
  return cadena, nil
}

func buscaCoincidencias(entrada string) ([]string, []string, error) {
  var titulo, interprete, album bool
  if !strings.Contains(entrada, "T:") || !strings.Contains(entrada, "A:") || !strings.Contains(entrada, "P:") {
    return "","", errors.New("Entrada no valida")
  }
  palabras := strings.Fields(entrada)
  for int i := 0 ; i < len(palabras) ; i++ {
    if strings.Contains(palabras[i], "T:") {
      coincidenciaTitulo, err := obtenerCoincidencia(palabras[i], "T")
      if err != nil { titulo = true }
    } else if strings.Contains(palabras[i], "P:") {
      coincidenciaInterprete, err:= obtenerCoincidencia(palabras[i], "P")
      if err != nil { interprete = true }
    } else if strings.Contains(palabras[i], "A:") {
      coincidenciaAlbum, err := obtenerCoincidencia(palabras[i], "A")
      if err != nil { album = true }
    }
  }
  banderas := make([]string, 0)
  coincidencias := make([]string, 0)
  if titulo {
    banderas = append(banderas,"title")
    coincidencias = append(coincidencias, coincidenciaTitulo)
  }
  if interprete {
    banderas = append(banderas,"performer")
    coincidencias = append(coincidencias, coincidenciaInterprete)
  }
  if album {
    banderas = append(banderas,"album")
    coincidencias = append(coincidencias,coincidenciaAlbum)
  }
  return banderas, coincidencias, nil
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
    idInterprete := obtenerIdInterprete(db, coincidencias[0])
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
  } else if banderaTitle {
    tabla, err := db.Query("SELECT id_album, id_performer, path, title, genre FROM rolas WHERE title LIKE '%" +  + "%'")
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

func obtenerIdAlbum(db *sql.DB, album string) []int {
  ids := make([]int, 0)
  tabla, err := db.Query("SELECT id_album FROM albums WHERE name LIKE '%" + album + "%'")
  if err != nil { panic(err) }
  for tabla.Next() {
    var id int
    tabla.Scan(&id)
    ids = append(ids, id)
  }
  return id
}

func obtenerIdInterprete(db *sql.DB, interprete string) []int {
  ids := make([]int, 0)
  tabla, err := db.Query("SELECT id_performer FROM performers WHERE name LIKE '%" + interprete + "%'")
  if err != nil { panic(err) }
  for tabla.Next() {
    var id int
    tabla.Scan(&id)
    ids = append(ids, id)
  }
  return ids
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
