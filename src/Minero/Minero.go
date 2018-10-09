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
  nuevaTabla := "CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, title TEXT, performer TEXT)"
  tabla, _ := db.Prepare(nuevaTabla)
  tabla.Exec()
  tabla, _ = db.Prepare("INSERT INTO people (title, performer) VALUES (?, ?)")
  tabla.Exec("Rod", "Vel")
  /* rows, _ := db.Query("SELECT id, title, performer FROM people")
  var id int
  var primerNombre string
  var ultimoNombre string
  for rows.Next() {
      rows.Scan(&id, &primerNombre, &ultimoNombre)
      fmt.Println(strconv.Itoa(id) + ": " + primerNombre + " " + ultimoNombre)
  }*/
  return db
}

func llenaBase(db *sql.DB)  {
  for _, cancion := range Canciones {
    tabla, _ := db.Prepare("INSERT INTO people (title, performer) VALUES (?, ?)")
    tabla.Exec(cancion.titulo, cancion.interprete)
  }
}
