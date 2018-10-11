package Etiquetas

import(
  "github.com/bogem/id3v2"
  "strconv"
  "os"
)

/* Maneja las etiquetas de las canciones */

 func ObtenerEtiquetas(pista *id3v2.Tag, direccion string) ([]string) {
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
