package main

import(
  "io/ioutil"
  "github.com/bogem/id3v2"
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
  "fmt"
  "os"
  "strings"
  "strconv"
  "os/user"
  "log"
)

/* Lee las etiquetas de cada archivo y crea la base */

type Cancion struct {
  interprete string
  titulo string
  album string
  año string
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
  fmt.Println( direccion.HomeDir )
  Canciones := make([]Cancion, 0)
  buscaArchivos(direccion.HomeDir + "/Music", Canciones)
  base := creaBase()
  llenaBase(base)
}

func obtenerCanciones() []Cancion {
  return Canciones
}

func buscaArchivos(directorio string, canciones []Cancion) {
  archivos, error := ioutil.ReadDir(directorio)
  if error != nil {
    panic(error)
  }
  for _, archivo := range archivos {
    nombre := archivo.Name()
    if strings.Contains(nombre, ".mp3") {
      creaCancion(directorio + "/" + nombre)
    }
    if archivo.IsDir() {
      buscaArchivos(directorio + "/" + nombre, canciones)
    }
  }
}

func creaCancion(direccion string) {
  pista, error := id3v2.Open(direccion, id3v2.Options{Parse: true})
  if error != nil {
    panic(error)
  }
  defer pista.Close()
  etiquetas := obtenerEtiquetas(pista, direccion)
  nuevaCancion := Cancion {
    interprete: etiquetas[0],
    titulo: etiquetas[1],
    album: etiquetas[2],
    año: etiquetas[3],
    genero: etiquetas[4],
    track: 0,
    ruta: direccion,
  }
  agregaCancion(nuevaCancion)
}

func agregaCancion(cancion Cancion)  {
  Canciones = append(Canciones, cancion)
}

func obtenerEtiquetas(pista *id3v2.Tag, direccion string) ([]string) {
  etiquetas := make([]string, 0)
  if pista.Artist() == "" {
    etiquetas = append(etiquetas, "Unknown")
  } else {
    etiquetas = append(etiquetas, pista.Artist())
  }
  if pista.Title() == "" {
    etiquetas = append(etiquetas, "Unknown")
  } else {
    etiquetas = append(etiquetas, pista.Title())
  }
  if pista.Album() == "" {
    etiquetas = append(etiquetas, "Unknown")
  } else {
    etiquetas = append(etiquetas, pista.Album())
  }
  if pista.Year() == "" {
    f,_ := os.Stat(direccion)
    t := f.ModTime()
    año := t.Year()
    etiquetas = append(etiquetas, strconv.Itoa(año))
  } else {
    etiquetas = append(etiquetas, pista.Year())
  }
  if pista.Genre() == "" {
    etiquetas = append(etiquetas, "Unknow")
  } else {
    etiquetas = append(etiquetas, pista.Genre())
  }
  return etiquetas
}

func creaBase() *sql.DB {
  db,_ := sql.Open("sqlite3", "../Base/base.db")
  PERFORMERS_TABLE := "CREATE TABLE IF NOT EXISTS performers (id_performer INTEGER PRIMARY KEY, id_type INTEGER, name TEXT)"
  performersTable, _ := db.Prepare(PERFORMERS_TABLE)
  performersTable.Exec()
  PERSONS_TABLE := "CREATE TABLE IF NOT EXISTS persons (id_person INTEGER PRIMARY KEY, stage_name TEXT, real_name TEXT, birth_date TEXT, death_date TEXT)"
  personsTable, _ := db.Prepare(PERSONS_TABLE)
  personsTable.Exec()
  TYPES_TABLE := "CREATE TABLE IF NOT EXISTS types (id_type INTEGER PRIMARY KEY, description TEXT)"
  typesTable,_ := db.Prepare(TYPES_TABLE)
  typesTable.Exec()
  GROUPS_TABLE := "CREATE TABLE IF NOT EXISTS groups (id_group INTEGER PRIMARY KEY, name TEXT, start_date TEXT, end_date TEXT)"
  groupsTable,_ := db.Prepare(GROUPS_TABLE)
  groupsTable.Exec()
  ALBUMS_TABLE := "CREATE TABLE IF NOT EXISTS albums (id_album INTEGER PRIMARY KEY, path TEXT, name TEXT, year INTEGER)"
  albumsTable,_ :=  db.Prepare(ALBUMS_TABLE)
  albumsTable.Exec()
  ROLAS_TABLE := "CREATE TABLE IF NOT EXISTS rolas (id_rola INTEGER PRIMARY KEY, id_performer INTEGER, id_album INTEGER, path TEXT, title TEXT, track INTEGER, year INTEGER, genre TEXT)"
  rolasTable,_ := db.Prepare(ROLAS_TABLE)
  rolasTable.Exec()
  IN_GROUP_TABLE := "CREATE TABLE IF NOT EXISTS in_group (id_person INTEGER, id_group INTEGER)"
  inGroupTable,_ := db.Prepare(IN_GROUP_TABLE)
  inGroupTable.Exec()
  return db
}

func llenaBase(db *sql.DB)  {
  for _, cancion := range Canciones {
    tabla, _ := db.Prepare("INSERT INTO rolas (path, title, track, year, genre) VALUES (?, ?, ?, ?, ?)")
    tabla.Exec(cancion.ruta, cancion.titulo,cancion.track, cancion.año, cancion.genero)
  }
  rows,_ := db.Query("SELECT id_rola, title, path FROM rolas")
  var  titulo string
  var dir string
  var id int
  for rows.Next() {
    rows.Scan(&id, &titulo, &dir)
    fmt.Println(strconv.Itoa(id) + ": " + titulo + " " + dir)
  }
}
