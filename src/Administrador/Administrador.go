package Administrador

import(
  "errors"
   "database/sql"
   _ "github.com/mattn/go-sqlite3"
   "github.com/RodrigoVelazquez99/Reproductor-MP3/src/Buscador"
   "github.com/RodrigoVelazquez99/Reproductor-MP3/src/Compilador"
   "strings"
)

var renglones []string
var base *sql.DB

func LeeEntrada(entrada string) bool {
  solicitudes, coincidencias, err := Compilador.BuscaCoincidencias(entrada)
  if err != nil { return false }
  err1 := buscaCancion(solicitudes, coincidencias)
  if err1 != nil { return false }
  return true
}

func IniciaBase() {
  renglones = make([]string, 0)
  nuevaBase, err := sql.Open("sqlite3", "../Base/base.db")
  if err != nil { panic(err) }
  base = nuevaBase
}

func ObtenerRenglones() []string {
  return renglones
}

func ObtenerBase() ([]string, error) {
  rolas := make([]string, 0)
  tabla, err := base.Query("SELECT id_album, id_performer, path, title, genre FROM rolas")
  if err != nil { return nil, err }
  var title, genre, path string
  var idAlbum, idPerformer int
  for tabla.Next() {
    tabla.Scan(&idAlbum, &idPerformer, &path, &title, &genre)
    performer := Buscador.BuscaInterprete(idPerformer, base)
    Album := Buscador.BuscaAlbum(idAlbum, base)
    rolas = append(rolas, title, performer, Album, genre)
  }
  tabla.Close()
  return rolas, nil
}

func VaciaRenglones()  {
  nuevosRenglones := make([]string, 0)
  renglones = nuevosRenglones
}

func creaRenglon(titulo string, interprete string, album string, genero string, ruta string )  {
  renglones = append(renglones, titulo, interprete, album, genero)
}

func buscaCancion(banderas []string, coincidencias []string) error {
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
    tabla, err := base.Query("SELECT id_album, id_performer, path, title, genre FROM rolas WHERE title LIKE '%" + string(coincidencias[0]) + "%'")
    if err != nil { return err }
    var title, genre, path string
    var idAlbum, idPerformer int
    for tabla.Next() {
      tabla.Scan(&idAlbum, &idPerformer, &path, &title, &genre)
      performer := Buscador.BuscaInterprete(idPerformer, base)
      Album := Buscador.BuscaAlbum(idAlbum, base)
      if strings.Contains(performer, string(coincidencias[1])) && strings.Contains(Album, string(coincidencias[2])) {
        creaRenglon(title, performer, Album, genre, path)
      }
    }
    tabla.Close()
    return nil
  } else if banderaTitle && banderaAlbum {
    tabla, err := base.Query("SELECT id_album, id_performer, path, title, genre FROM rolas WHERE title LIKE '%" + string(coincidencias[0]) + "%'")
    if err != nil { return err }
    var title, genre, path string
    var idAlbum, idPerformer int
    for tabla.Next() {
      tabla.Scan(&idAlbum, &idPerformer, &path, &title, &genre)
      performer := Buscador.BuscaInterprete(idPerformer, base)
      Album := Buscador.BuscaAlbum(idAlbum, base)
      if strings.Contains(Album, string(coincidencias[1])) {
        creaRenglon(title, performer, Album, genre, path)
      }
    }
    tabla.Close()
    return nil
  } else if banderaTitle && banderaPerformer{
    tabla, err := base.Query("SELECT id_album, id_performer, path, title, genre FROM rolas WHERE title LIKE '%" + string(coincidencias[0]) + "%'")
    if err != nil { return err }
    var title, genre, path string
    var idAlbum, idPerformer int
    for tabla.Next() {
      tabla.Scan(&idAlbum, &idPerformer, &path, &title, &genre)
      performer := Buscador.BuscaInterprete(idPerformer, base)
      Album := Buscador.BuscaAlbum(idAlbum, base)
      if strings.Contains(performer, string(coincidencias[1])) {
        creaRenglon(title, performer, Album, genre, path)
      }
    }
    tabla.Close()
    return nil
  } else if banderaAlbum && banderaPerformer {
    albums := Buscador.BuscaAlbums(base, string(coincidencias[1]))
    for _,album := range albums {
      idAlbum := Buscador.ObtenerIdAlbum(base, album)
      tabla, err := base.Query("SELECT id_performer, path, title, genre FROM rolas WHERE id_album=?",idAlbum)
      if err != nil { return err }
      var title, genre, path string
      var idPerformer int
      for tabla.Next() {
        tabla.Scan(&idPerformer, &path, &title, &genre)
        performer := Buscador.BuscaInterprete(idPerformer, base)
        if strings.Contains(performer, coincidencias[0]){
          creaRenglon(title, performer, album, genre, path)
        }
      }
      tabla.Close()
    }
    return nil
  } else if banderaAlbum {
    albums := Buscador.BuscaAlbums(base, coincidencias[0])
    for _,album := range albums {
      idAlbum := Buscador.ObtenerIdAlbum(base, album)
      tabla, err := base.Query("SELECT id_performer, path, title, genre FROM rolas WHERE id_album=?",idAlbum)
      if err != nil { return err }
      var title, genre, path string
      var idPerformer int
      for tabla.Next() {
        tabla.Scan(&idPerformer, &path, &title, &genre)
        performer := Buscador.BuscaInterprete(idPerformer, base)
        creaRenglon(title, performer, album, genre, path)
      }
      tabla.Close()
    }
    return nil
  } else if banderaPerformer {
    interpretes := Buscador.BuscaInterpretes(base, coincidencias[0])
    for _, interprete := range interpretes {
      idInterprete := Buscador.ObtenerIdInterprete(base, interprete)
      tabla, err := base.Query("SELECT id_album, path, title, genre FROM rolas WHERE id_performer=?",idInterprete)
      if err != nil { return err }
      var title, genre, path string
      var idAlbum int
      for tabla.Next() {
        tabla.Scan(&idAlbum, &path, &title, &genre)
        performer := Buscador.BuscaInterprete(idInterprete, base)
        album := Buscador.BuscaAlbum(idAlbum, base)
        creaRenglon(title, performer, album, genre, path)
      }
      tabla.Close()
    }
    return nil
  } else if banderaTitle {
    tabla, err := base.Query("SELECT id_album, id_performer, path, title, genre FROM rolas WHERE title LIKE '%" + string(coincidencias[0]) + "%'")
    if err != nil { return err }
    var title, genre, path string
    var idAlbum, idPerformer int
    for tabla.Next() {
      tabla.Scan(&idAlbum, &idPerformer, &path, &title, &genre)
      performer := Buscador.BuscaInterprete(idPerformer, base)
      Album := Buscador.BuscaAlbum(idAlbum, base)
      creaRenglon(title, performer, Album, genre, path)
    }
    tabla.Close()
    return nil
  } else {
    return errors.New("fallo")
  }
}
