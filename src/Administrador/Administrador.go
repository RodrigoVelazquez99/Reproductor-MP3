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
    rolas = append(rolas, title, performer, Album, genre, path)
  }
  tabla.Close()
  renglones = rolas
  return rolas, nil
}

func VaciaRenglones()  {
  nuevosRenglones := make([]string, 0)
  renglones = nuevosRenglones
}

func creaRenglon(titulo string, interprete string, album string, genero string, ruta string )  {
  renglones = append(renglones, titulo, interprete, album, genero, ruta)
}

/**
* Modifica el renglon de la interfaz
*/
func cambiaRenglon(nuevoTitulo string, nuevoInterprete string, nuevoAlbum string, nuevoGenero string, ruta string) error {
  for i := 0; i < len(renglones) ; i++ {
    if renglones[i] == ruta {
      renglones[i-1] = nuevoGenero
      renglones[i-2] = nuevoAlbum
      renglones[i-3] = nuevoInterprete
      renglones[i-4] = nuevoTitulo
      return nil
    }
  }
  return errors.New("Ocurrio un error")
}

  /* Inserta una persona a la base de datos */
  func InsertaInterpretePersona(person_na string, person_rn string, person_bd string, person_dd string){
    stm, err := base.Prepare("INSERT INTO persons (stage_name, real_name, birth_date, death_date) VALUES (?, ?, ?, ?)")
    if err != nil { panic(err) }
    stm.Exec(person_na, person_rn, person_bd, person_dd)
    stm.Close()
    stm1, err1 := base.Prepare("INSERT INTO performers (id_type, name) VALUES (?, ?)")
    if err1 != nil { panic(err1) }
    stm1.Exec(1, person_na)
    stm1.Close()
  }

  /* Inserta un grupo a la base de datos */
  func InsetaInterpreteGrupo(group_n string, group_sd string, group_ed string){
    stm, err := base.Prepare("INSERT INTO groups (name, start_date, end_date) VALUES (?,?,?)")
    if err != nil { panic(err) }
    stm.Exec(group_n, group_sd, group_ed)
    stm.Close()
    stm1, err1 := base.Prepare("INSERT INTO performers (id_type, name) VALUES (?, ?)")
    if err1 != nil { panic(err1) }
    stm1.Exec(2, group_n)
    stm1.Close()
  }

  /* Inserta la ruta de una imagen a la base de datos */
  //func InsertaImagen(rutaCancion string, rutaImagen string ){}

/**
* Busca las canciones que coinciden con la entrada y las agrega
* a los renglones que se muestran en la interfaz
*/
func buscaCancion(banderas []string, coincidencias []string) error {
  VaciaRenglones()
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
    tabla, err := base.Query("SELECT id_album, id_performer, path, title, genre FROM rolas WHERE title LIKE '" + string(coincidencias[0]) + "%'")
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
    tabla, err := base.Query("SELECT id_album, id_performer, path, title, genre FROM rolas WHERE title LIKE '" + string(coincidencias[0]) + "%'")
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
    tabla, err := base.Query("SELECT id_album, id_performer, path, title, genre FROM rolas WHERE title LIKE '" + string(coincidencias[0]) + "%'")
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
    tabla, err := base.Query("SELECT id_album, id_performer, path, title, genre FROM rolas WHERE title LIKE '" + string(coincidencias[0]) + "%'")
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

func BuscaPorRuta(ruta string) []string {
  etiquetas := make([]string, 0)
  tabla,err := base.Query("SELECT title, id_performer, id_album, genre, image FROM rolas WHERE path=?",ruta)
  if err != nil { panic(err) }
  var titulo, genero, image string
  var idInterprete, idAlbum int
  if tabla.Next() {
    tabla.Scan(&titulo, &idInterprete, &idAlbum, &genero, &image)
    album := Buscador.BuscaAlbum(idAlbum, base)
    interprete := Buscador.BuscaInterprete(idInterprete, base)
    etiquetas = append(etiquetas, titulo, interprete, album, genero, image)
  }
  tabla.Close()
  return etiquetas
}

func CambiaEtiquetas(ruta string, nuevoTitulo string, nuevoInterprete string, nuevoAlbum string, nuevoGenero string, nuevaImagen string) {
  idAlbum := Buscador.ObtenerIdAlbum(base, nuevoAlbum)
  idInterprete := Buscador.ObtenerIdInterprete(base, nuevoInterprete)
  tabla, err := base.Prepare("UPDATE rolas SET id_album=?, id_performer=?, path=?, title=?, genre=?, image=? WHERE path=?")
  if err != nil { panic(err) }
  tabla.Exec(idAlbum, idInterprete, ruta, nuevoTitulo, nuevoGenero, nuevaImagen, ruta)
  cambiaRenglon(nuevoTitulo, nuevoInterprete, nuevoAlbum, nuevoGenero, ruta)
  tabla.Close()
}
