package Compilador

import (
  "strings"
  "errors"
)

func ObtenerCoincidencia(palabra string, clave string) string {
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
