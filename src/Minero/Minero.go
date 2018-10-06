package main

import(
  "io/ioutil"
  "github.com/bogem/id3v2"
  "fmt"
  "os"
  "strings"
  "strconv"
)

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
  Canciones := make([]Cancion, 0)
  buscaArchivos("/home/rodrigofvc/Music", Canciones)
  fmt.Println(obtenerCanciones())
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
