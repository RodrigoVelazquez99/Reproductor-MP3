package Compilador

import (
  "strings"
  "errors"
)

func obtenerCoincidencia(palabras []string, clave string) string {
  	for i:= 0 ; i < len(palabras) ; i++{
  		if palabras[i] != clave && !strings.Contains(palabras[i],clave) {
  	        	continue
  		}
  		if palabras[i] == clave {
  			if i + 1 < len(palabras){
          if strings.HasSuffix(palabras[i+1],",") {
            s := palabras[i+1]
            s = s[:len(palabras[i+1])-1]
            return s
          }
  			return palabras[i+1]
  			}
  		}
  		if strings.Contains(palabras[i],clave){
  			palabra := palabras[i]
  			arr := palabra[2:]
        if strings.HasSuffix(arr,",") {
          arr = arr[:len(arr)-1]
        }
  			return arr
  		}
  	}
  	return "Ingresa bien la busqueda"
  }

func BuscaCoincidencias(entrada string) ([]string, []string, error) {
  var titulo, interprete, album bool
  var coincidenciaTitulo, coincidenciaAlbum, coincidenciaInterprete string
  if !strings.Contains(entrada, "T:") && !strings.Contains(entrada, "A:") && !strings.Contains(entrada, "P:") {
    return nil,nil, errors.New("Entrada no valida:  " + entrada)
  }
  palabras := strings.Fields(entrada)
  for i := 0 ; i < len(palabras) ; i++ {
    if !titulo && strings.Contains(palabras[i], "T:") {
      coincidenciaTitulo = obtenerCoincidencia(palabras, "T:")
      titulo = true
    } else if !interprete && strings.Contains(palabras[i], "P:") {
      coincidenciaInterprete = obtenerCoincidencia(palabras, "P:")
      interprete = true
    } else if !album && strings.Contains(palabras[i], "A:") {
      coincidenciaAlbum = obtenerCoincidencia(palabras, "A:")
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
