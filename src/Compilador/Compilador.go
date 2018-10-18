package Compilador

import (
  "strings"
  "errors"
)

func obtenerCoincidencia(palabra string, clave string) (string){
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
  return cadena
}
func BuscaCoincidencias(entrada string) ([]string, []string, error) {
  var titulo, interprete, album bool
  var coincidenciaTitulo, coincidenciaAlbum, coincidenciaInterprete string
  if !strings.Contains(entrada, "T:") && !strings.Contains(entrada, "A:") && !strings.Contains(entrada, "P:") {
    return nil,nil, errors.New("Entrada no valida:  " + entrada)
  }
  palabras := strings.Fields(entrada)
  for i := 0 ; i < len(palabras) ; i++ {
    if strings.Contains(palabras[i], "T:") {
      coincidenciaTitulo = obtenerCoincidencia(palabras[i], "T")
      titulo = true
    } else if strings.Contains(palabras[i], "P:") {
      coincidenciaInterprete = obtenerCoincidencia(palabras[i], "P")
      interprete = true
    } else if strings.Contains(palabras[i], "A:") {
      coincidenciaAlbum = obtenerCoincidencia(palabras[i], "A")
      album = true
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
