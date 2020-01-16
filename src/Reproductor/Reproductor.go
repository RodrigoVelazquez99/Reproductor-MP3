package main

import(
  "github.com/gotk3/gotk3/gtk"
  "errors"
  "github.com/gotk3/gotk3/glib"
  "github.com/RodrigoVelazquez99/Reproductor-MP3/src/Administrador"
  "github.com/RodrigoVelazquez99/Reproductor-MP3/src/Minero"
)

const (
  COLUMN_TITLE = iota
	COLUMN_PERFORMER
  COLUMN_ALBUM
  COLUMN_GENRE
  COLUMN_PATH
  path = "../Interfaz/Interfaz.glade"
)

func main()  {
  Administrador.IniciaBase()

  gtk.Init(nil)
  builder,err := build(path)
  if err != nil { panic(err) }
  ventana, err := window(builder,"window1")
  if err != nil { panic(err) }
  ventanaEditar,err := window(builder,"windowEdit")
  if err != nil { panic(err) }
  ventana2, err := scrolledWindow(builder)
  if err != nil { panic(err) }
  boton, err := button(builder, "button1")
  if err != nil { panic(err) }
  botonMinero, err1 := button(builder,"buttonMiner")
  if err1 != nil { panic(err) }
  botonEditar, err2 := button(builder, "editTags")
  if err2 != nil { panic(err) }
  botonGuardar, err3 := button(builder, "editTag")
  if err3 != nil { panic(err3) }
  grid, err := grid(builder)
  grid.SetOrientation(gtk.ORIENTATION_VERTICAL)
  if err != nil { panic(err) }
  entry, err := entry(builder)
  titleEntry, performerEntry, albumEntry, genreEntry, err := entryEdit(builder)
  if err != nil{ panic(err) }
  var rutaCancionSeleccionada  string
  treeView, listStore := creaTreeView()
  treeView.SetSearchEntry(entry)


  botonMinero.Connect("clicked", func ()  {
    Minero.Mina()
    rolas, err := Administrador.ObtenerBase()
    if err == nil {
      listStore.Clear()
      actualizaLista(listStore, rolas)
    }
  })

  boton.Connect("clicked", func ()  {
      entrada, err := entry.GetText()
      if entrada != "" && err == nil {
        entry.SetText("")
        ok := Administrador.LeeEntrada(entrada)
        listStore.Clear()
        if ok {
          treeView.ExpandAll()
          renglones := Administrador.ObtenerRenglones()
          actualizaLista(listStore,renglones)
        }
      }
  })

  /**
  * Cuando se modifican las etiquetas de una cancion se actualiza la interfaz y
  * se oculta la ventana de editar
  */
  botonGuardar.Connect("clicked", func ()  {
    nuevoTitulo, err := titleEntry.GetText()
    nuevoInterprete, err := performerEntry.GetText()
    nuevoAlbum, err := albumEntry.GetText()
    nuevoGenero, err := genreEntry.GetText()
    if err != nil { panic(err) }
    Administrador.CambiaEtiquetas(rutaCancionSeleccionada, nuevoTitulo, nuevoInterprete, nuevoAlbum, nuevoGenero)
    listStore.Clear()
    renglones := Administrador.ObtenerRenglones()
    actualizaLista(listStore, renglones)
    ventanaEditar.Hide()
  })

  botonEditar.Connect("clicked", func () {
    if rutaCancionSeleccionada != "" {
      etiquetas := Administrador.BuscaPorRuta(rutaCancionSeleccionada)
      titleEntry.SetText(etiquetas[0])
      performerEntry.SetText(etiquetas[1])
      albumEntry.SetText(etiquetas[2])
      genreEntry.SetText(etiquetas[3])
      ventanaEditar.ShowAll()
    }
  })

  seleccion, err := treeView.GetSelection()
	if err != nil { panic(err) }
	seleccion.SetMode(gtk.SELECTION_SINGLE)
  seleccion.Connect("changed", func ()  {
    model, iter, ok := seleccion.GetSelected()
    if ok {
      columna4,_ := model.(*gtk.TreeModel).GetValue(iter,4)
      ruta,_ := columna4.GetString()
      rutaCancionSeleccionada = ruta
      }
  })
  ventana2.Add(treeView)
  ventana2.ShowAll()
  ventanaEditar.Connect("destroy", func () {
    ventanaEditar.Hide()
  })
  ventana.SetTitle("Reproductor-MP3")
	ventana.SetDefaultSize(800, 500)
	ventana.Connect("destroy", func ()  {
    gtk.MainQuit()
  })
	ventana.ShowAll()
	gtk.Main()
}

func entryEdit(builder *gtk.Builder) (*gtk.Entry, *gtk.Entry, *gtk.Entry, *gtk.Entry, error)  {
  object, err := builder.GetObject("setTitle")
  if err != nil { panic(err) }
  entryTitle, ok := object.(*gtk.Entry)
  if !ok { return nil, nil, nil, nil, err }
  object1, err := builder.GetObject("setPerformer")
  if err != nil { panic(err) }
  entryPerformer,ok := object1.(*gtk.Entry)
  if !ok { return nil, nil, nil, nil, err }
  object2, err := builder.GetObject("setAlbum")
  if err != nil { panic(err) }
  entryAlbum,ok := object2.(*gtk.Entry)
  if !ok { return nil, nil, nil, nil, err }
  object3, err := builder.GetObject("setGenre")
  if err != nil { panic(err) }
  entryGenre,ok := object3.(*gtk.Entry)
  if !ok { return nil, nil, nil, nil, err }
  return entryTitle, entryPerformer, entryAlbum, entryGenre, nil
}


func build(ruta string) (*gtk.Builder, error)  {
  	builder, err := gtk.BuilderNew()
  	if err != nil {
  		return nil, err
  	}
  	if ruta != "" {
  		err = builder.AddFromFile(ruta)
  		if err != nil {
  			return nil, errors.New("Error")
  		}
  	}
  	return builder, nil
}


func window(builder *gtk.Builder, tipo string) (*gtk.Window ,error) {
  object, err := builder.GetObject(tipo)
	if err != nil {
		return nil, err
	}
	ventana, ok := object.(*gtk.Window)
	if !ok {
		return nil, err
	}
	return ventana, nil
}

func scrolledWindow(builder *gtk.Builder) (*gtk.ScrolledWindow, error) {
  object, err := builder.GetObject("window2")
	if err != nil {
		return nil, err
	}
	ventana, ok := object.(*gtk.ScrolledWindow)
	if !ok {
		return nil, err
	}
	return ventana, nil
}

func dialog(builder *gtk.Builder) (*gtk.Dialog, error)  {
  object, err := builder.GetObject("windowEdit")
	if err != nil {
		return nil, err
	}
	ventana, ok := object.(*gtk.Dialog)
	if !ok {
		return nil, err
	}
	return ventana, nil
}


func button(builder *gtk.Builder, tipo string) (*gtk.Button, error)  {
  object, err := builder.GetObject(tipo)
  if err != nil {
    panic(err)
  }
  boton, ok := object.(*gtk.Button)
  if !ok {
    return nil, err
  }
  return boton, nil
}

func grid(builder *gtk.Builder) (*gtk.Grid, error)  {
  object, err := builder.GetObject("grid1")
  if err != nil {
    panic(err)
  }
  grid, ok := object.(*gtk.Grid)
  if !ok {
    return nil, err
  }
  return grid, nil
}

func entry(builder *gtk.Builder) (*gtk.Entry, error) {
  object, err := builder.GetObject("entry1")
  if err != nil {
    panic(err)
  }
  entry, ok := object.(*gtk.Entry)
  if !ok {
    return nil, err
  }
  return entry,nil
}

func creaColumna(nombre string, id int) *gtk.TreeViewColumn {
	cellRenderer, err := gtk.CellRendererTextNew()
	if err != nil {
		panic(err)
	}
	columna, err := gtk.TreeViewColumnNewWithAttribute(nombre, cellRenderer, "text", id)
	if err != nil {
		panic(err)
	}
	return columna
}

func creaTreeView() (*gtk.TreeView, *gtk.ListStore) {
	treeView, err := gtk.TreeViewNew()
	if err != nil { panic(err) }
	treeView.AppendColumn(creaColumna("Title", COLUMN_TITLE))
	treeView.AppendColumn(creaColumna("Performer", COLUMN_PERFORMER))
  treeView.AppendColumn(creaColumna("Album", COLUMN_ALBUM))
	treeView.AppendColumn(creaColumna("Genre", COLUMN_GENRE))
  Paths := creaColumna("Path", COLUMN_PATH)
  Paths.SetVisible(false)
  treeView.AppendColumn(Paths)
	listStore, err := gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING)
	if err != nil { panic(err) }
	treeView.SetModel(listStore)
	return treeView, listStore
}

func actualizaLista(listStore *gtk.ListStore, canciones []string)  {
  i := 0
  for i < len(canciones) {
    nuevoRenglon(listStore,canciones[i],canciones[i+1],canciones[i+2],canciones[i+3], canciones[i+4])
    i += 5
  }
}

func nuevoRenglon(listStore *gtk.ListStore, titulo string , interprete string, album string, genero string, ruta string) {
	iter := listStore.Append()
	err := listStore.Set(iter,
		[]int{COLUMN_TITLE, COLUMN_PERFORMER, COLUMN_ALBUM, COLUMN_GENRE, COLUMN_PATH},
		[]interface{}{titulo, interprete, album, genero, ruta})
	if err != nil { panic(err) }
}
