package main

import(
  "io/ioutil"
  "github.com/bogem/id3v2"
  "github.com/RodrigoVelazquez99/Reproductor-MP3/src/Etiquetas"
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
  //"fmt"
  "strings"
  //"strconv"
  "os/user"
  "log"
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
}

var Canciones []Cancion

func main() {
  direccion, err := user.Current()
    if err != nil {
        log.Fatal( err )
    }
  //fmt.Println( direccion.HomeDir )
  Canciones = make([]Cancion, 0)
  buscaArchivos(direccion.HomeDir + "/Music")
  base := creaBase()
  defer base.Close()
  llenaBase(base)
}

func buscaArchivos(directorio string) {
  archivos, error := ioutil.ReadDir(directorio)
  if error != nil {
    panic(error)
  }
  for _, archivo := range archivos {
    nombre := archivo.Name()
    ruta := directorio + "/" + nombre
    if strings.Contains(nombre, ".mp3") {
        creaCancion(ruta)
    }
    if archivo.IsDir() {
      buscaArchivos(ruta)
    }
  }
}

func creaCancion(direccion string) {
  pista, error := id3v2.Open(direccion, id3v2.Options{Parse: true})
  if error != nil {
    panic(error)
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
  }
  Canciones = append(Canciones, nuevaCancion)
}

func creaBase() *sql.DB {
  db,_ := sql.Open("sqlite3", "../Base/base.db")
  TYPES_TABLE := "CREATE TABLE IF NOT EXISTS types (id_type INTEGER PRIMARY KEY, description TEXT)"
  typesTable,_ := db.Prepare(TYPES_TABLE)
  typesTable.Exec()
  i,_ := db.Query("SELECT description FROM types WHERE description=?","Person")
  if !i.Next() {
    tipos,_ := db.Prepare("INSERT INTO types (description) VALUES (?)")
    tipos.Exec("Person")
    tipos.Exec("Group")
    tipos.Exec("Unknown")
  }
  PERFORMERS_TABLE := "CREATE TABLE IF NOT EXISTS performers (id_performer INTEGER PRIMARY KEY, id_type INTEGER, name TEXT, FOREIGN KEY (id_type) REFERENCES types(id_type))"
  performersTable, _ := db.Prepare(PERFORMERS_TABLE)
  performersTable.Exec()
  PERSONS_TABLE := "CREATE TABLE IF NOT EXISTS persons (id_person INTEGER PRIMARY KEY, stage_name TEXT, real_name TEXT, birth_date TEXT, death_date TEXT)"
  personsTable, _ := db.Prepare(PERSONS_TABLE)
  personsTable.Exec()
  GROUPS_TABLE := "CREATE TABLE IF NOT EXISTS groups (id_group INTEGER PRIMARY KEY, name TEXT, start_date TEXT, end_date TEXT)"
  groupsTable,_ := db.Prepare(GROUPS_TABLE)
  groupsTable.Exec()
  ALBUMS_TABLE := "CREATE TABLE IF NOT EXISTS albums (id_album INTEGER PRIMARY KEY, path TEXT, name TEXT, year INTEGER)"
  albumsTable,_ :=  db.Prepare(ALBUMS_TABLE)
  albumsTable.Exec()
  ROLAS_TABLE := "CREATE TABLE IF NOT EXISTS rolas (id_rola INTEGER PRIMARY KEY, id_performer INTEGER, id_album INTEGER, path TEXT, title TEXT, track INTEGER, year INTEGER, genre TEXT, FOREIGN KEY (id_performer) REFERENCES performers(id_performer), FOREIGN KEY (id_album) REFERENCES albums(id_album))"
  rolasTable,_ := db.Prepare(ROLAS_TABLE)
  rolasTable.Exec()
  IN_GROUP_TABLE := "CREATE TABLE IF NOT EXISTS in_group (id_person INTEGER, id_group INTEGER, PRIMARY KEY (id_person, id_group), FOREIGN KEY (id_person) REFERENCES persons(id_person), FOREIGN KEY (id_group) REFERENCES groups(id_group))"
  inGroupTable,_ := db.Prepare(IN_GROUP_TABLE)
  inGroupTable.Exec()
  return db
}

func llenaBase(db *sql.DB)  {
  llenaRolas(db)
}

func llenaInterpretes(db *sql.DB)  {

}

func llenaPersonas(db *sql.DB)  {

}

func llenaGrupos(db *sql.DB)  {

}

func llenaAlbumes(db *sql.DB)  {

}

func llenaRolas(db *sql.DB)  {
  for _, cancion := range Canciones {
    registrado,_ := db.Query("SELECT path FROM rolas WHERE path=?",cancion.ruta)
    if registrado.Next() {
      continue
    }
    tabla, _ := db.Prepare("INSERT INTO rolas (path, title, track, year, genre) VALUES (?, ?, ?, ?, ?)")
    tabla.Exec(cancion.ruta, cancion.titulo,cancion.track, cancion.año, cancion.genero)
  }
  /*rows,_ := db.Query("SELECT id_rola, title, path, track, year, genre FROM rolas")
  var titulo string
  var dir string
  var genre string
  var track, year int
  var id int
  for rows.Next() {
  rows.Scan(&id, &titulo, &dir, &track, &year, &genre)
  fmt.Println(strconv.Itoa(id) + ": " + titulo + " " + dir)
  }*/
}

func ingresaGrupo()  {

}
