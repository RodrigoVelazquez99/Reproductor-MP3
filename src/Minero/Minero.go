package Minero

import(
  "io/ioutil"
  "github.com/bogem/id3v2"
  "github.com/RodrigoVelazquez99/Reproductor-MP3/src/Etiquetas"
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
  "strings"
  "os/user"
  "log"
  "os"
)

/* Lee las etiquetas de cada archivo y crea la base */

type Cancion struct {
  interprete string
  titulo string
  album string
  año int
  genero string
  track int
  ruta string
  carpeta string
}

var Canciones []Cancion

func Mina() {
  direccion, err := user.Current()
    if err != nil {
        log.Fatal( err )
    }
  Canciones = make([]Cancion, 0)
  buscaArchivos(direccion.HomeDir + "/Música")
  err1 := os.MkdirAll("../Base",0755)
  if err1 != nil { panic(err1) }
  base := creaBase()
  defer base.Close()
  llenaBase(base)
}

func buscaArchivos(carpeta string) {
  archivos, error := ioutil.ReadDir(carpeta)
  if error != nil {
    panic(error)
  }
  for _, archivo := range archivos {
    nombre := archivo.Name()
    ruta := carpeta + "/" + nombre
    if strings.Contains(nombre, ".mp3") {
        err := creaCancion(ruta,carpeta)
        if err != nil {
          continue
        }
    }
    if archivo.IsDir() {
      buscaArchivos(ruta)
    }
  }
}

func creaCancion(direccion string, directorio string) error {
  pista, err := id3v2.Open(direccion, id3v2.Options{Parse: true})
  if err != nil {
    return err
  }
  defer pista.Close()
  etiquetas, year := Etiquetas.ObtenerEtiquetas(pista, direccion)
  nuevaCancion := Cancion {
    interprete: etiquetas[0],
    titulo: etiquetas[1],
    album: etiquetas[2],
    año: year,
    genero: etiquetas[3],
    track: 0,
    ruta: direccion,
    carpeta: directorio + "/",
  }
  Canciones = append(Canciones, nuevaCancion)
  return nil
}

func creaBase() *sql.DB {
  db,err := sql.Open("sqlite3", "../Base/base.db")
  manejaError(err)
  TYPES_TABLE := "CREATE TABLE IF NOT EXISTS types (id_type INTEGER PRIMARY KEY, description TEXT)"
  typesTable,err1 := db.Prepare(TYPES_TABLE)
  manejaError(err1)
  typesTable.Exec()
  i,err2 := db.Query("SELECT description FROM types WHERE description=?","Person")
  manejaError(err2)
  if !i.Next() {
    tipos,_ := db.Prepare("INSERT INTO types (description) VALUES (?)")
    tipos.Exec("Person")
    tipos.Exec("Group")
    tipos.Exec("Unknown")
  }
  i.Close()
  PERFORMERS_TABLE := "CREATE TABLE IF NOT EXISTS performers (id_performer INTEGER PRIMARY KEY, id_type INTEGER, name TEXT, FOREIGN KEY (id_type) REFERENCES types(id_type))"
  performersTable, err3 := db.Prepare(PERFORMERS_TABLE)
  manejaError(err3)
  performersTable.Exec()
  PERSONS_TABLE := "CREATE TABLE IF NOT EXISTS persons (id_person INTEGER PRIMARY KEY, stage_name TEXT, real_name TEXT, birth_date TEXT, death_date TEXT)"
  personsTable, err4 := db.Prepare(PERSONS_TABLE)
  manejaError(err4)
  personsTable.Exec()
  GROUPS_TABLE := "CREATE TABLE IF NOT EXISTS groups (id_group INTEGER PRIMARY KEY, name TEXT, start_date TEXT, end_date TEXT)"
  groupsTable,err5 := db.Prepare(GROUPS_TABLE)
  manejaError(err5)
  groupsTable.Exec()
  ALBUMS_TABLE := "CREATE TABLE IF NOT EXISTS albums (id_album INTEGER PRIMARY KEY, path TEXT, name TEXT, year INTEGER)"
  albumsTable, err6 :=  db.Prepare(ALBUMS_TABLE)
  manejaError(err6)
  albumsTable.Exec()
  ROLAS_TABLE := "CREATE TABLE IF NOT EXISTS rolas (id_rola INTEGER PRIMARY KEY, id_performer INTEGER, id_album INTEGER, path TEXT, title TEXT, track INTEGER, year INTEGER, genre TEXT, FOREIGN KEY (id_performer) REFERENCES performers(id_performer), FOREIGN KEY (id_album) REFERENCES albums(id_album))"
  rolasTable, err7 := db.Prepare(ROLAS_TABLE)
  manejaError(err7)
  rolasTable.Exec()
  IN_GROUP_TABLE := "CREATE TABLE IF NOT EXISTS in_group (id_person INTEGER, id_group INTEGER, PRIMARY KEY (id_person, id_group), FOREIGN KEY (id_person) REFERENCES persons(id_person), FOREIGN KEY (id_group) REFERENCES groups(id_group))"
  inGroupTable,err8 := db.Prepare(IN_GROUP_TABLE)
  manejaError(err8)
  inGroupTable.Exec()
  return db
}

func llenaBase(db *sql.DB)  {
  for _, cancion := range Canciones {
    rolaRegistrada, albumRegistrado, interpreteRegistrado := registrado(cancion.ruta, cancion.album, cancion.interprete, db)
    if rolaRegistrada {
      continue
    } else {
      if !albumRegistrado {
        registraAlbum(cancion.carpeta, cancion.album, cancion.año, db)
      }
      if !interpreteRegistrado {
        registraInterprete(db, cancion.interprete, "Unknown")
      }
      idAlbum := obtenerIdAlbum(db, cancion.album)
      idInterprete := obtenerIdInterprete(db, cancion.interprete)
      registraRola(cancion.ruta, cancion.titulo, cancion.track, cancion.año, cancion.genero, idAlbum, idInterprete, db)
    }
  }
}

func registrado(ruta string, album string, interprete string, db *sql.DB) (bool,bool,bool) {
  rolaRegistrada := false
  albumRegistrado := false
  interpreteRegistrado := false
  registrada,err := db.Query("SELECT path FROM rolas WHERE path=?",ruta)
  manejaError(err)
  if registrada.Next() {
    registrada.Close()
    return true, true, true
  }
  registrada1, err1 := db.Query("SELECT name FROM albums WHERE name=?",album)
  manejaError(err1)
  if registrada1.Next() {
    albumRegistrado = true
  }
  registrada2, err2 := db.Query("SELECT name FROM performers WHERE name=?",interprete)
  manejaError(err2)
  if registrada2.Next() {
    interpreteRegistrado = true
  }
  registrada.Close()
  registrada1.Close()
  registrada2.Close()
  return rolaRegistrada, albumRegistrado, interpreteRegistrado
}

func registraInterprete(db *sql.DB, interprete string, descripcion string)  {
  idType := obtenerIdType(descripcion)
  tabla, err := db.Prepare("INSERT INTO performers (id_type, name) VALUES (?, ?)")
  manejaError(err)
  tabla.Exec(idType, interprete)
  tabla.Close()
}

func registraAlbum(carpeta string, album string, año int, db *sql.DB)  {
  tabla, err := db.Prepare("INSERT INTO albums (path, name, year ) VALUES (?, ?, ?)")
  manejaError(err)
  tabla.Exec(carpeta, album,año)
  tabla.Close()
}

func registraRola(ruta string, titulo string, track int, año int, genero string, idAlbum int, idPerformer int, db *sql.DB)   {
  tabla1, err := db.Prepare("INSERT INTO rolas (path, title, track, year, genre, id_album, id_performer) VALUES (?, ?, ?, ?, ?, ?, ?)")
  manejaError(err)
  tabla1.Exec(ruta,titulo, track, año, genero, idAlbum, idPerformer)
  tabla1.Close()
}

func obtenerIdInterprete(db *sql.DB, interpreteRequerido string) int {
  tabla, err := db.Query("SELECT id_performer FROM performers WHERE name=?",interpreteRequerido)
  manejaError(err)
  var aux int
  if tabla.Next() {
    tabla.Scan(&aux)
  }
  tabla.Close()
  return aux
}

func obtenerIdAlbum(db *sql.DB, albumRequerido string) int {
  tabla, err := db.Query("SELECT id_album FROM albums WHERE name =?",albumRequerido)
  manejaError(err)
  var aux int
  if tabla.Next() {
    tabla.Scan(&aux)
  }
  tabla.Close()
  return aux
}

func obtenerIdType(descripcion string) int {
  switch descripcion {
  case "Person":
    return 0
  case "Group":
    return 1
  default:
    return 2
  }
}

func manejaError(err error)  {
  if err != nil {
    panic(err)
  }
}
