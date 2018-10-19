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
  path = "../Interfaz/Interfaz.glade"
)

func main()  {
  Administrador.IniciaBase()

  gtk.Init(nil)
  builder,err := build(path)
  if err != nil {
    panic(err)
  }
  ventana, err := window(builder)
  if err != nil {
    panic(err)
  }
  boton, err := button(builder, "button1")
  if err != nil {
    panic(err)
  }
  botonMinero, err1 := button(builder,"buttonMiner")
  if err1 != nil {
    panic(err)
  }
  grid, err := grid(builder)
  grid.SetOrientation(gtk.ORIENTATION_VERTICAL)
  if err != nil  {
    panic(err)
  }
  entry, err := entry(builder)
  treeView, listStore := creaTreeView()
  treeView.SetSearchEntry(entry)
  treeView.SetHoverSelection(true)
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
        if ok {
          treeView.ExpandAll()
          listStore.Clear()
          renglones := Administrador.ObtenerRenglones()
          actualizaLista(listStore,renglones)
          Administrador.VaciaRenglones()
        } else {
          listStore.Clear()

        }
      }
  })
  ventana.SetTitle("Reproductor-MP3")
	ventana.SetDefaultSize(800, 500)
	ventana.Connect("destroy", func ()  {
    gtk.MainQuit()
  })
  grid.Add(treeView)
	ventana.ShowAll()
	gtk.Main()
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


func window(builder *gtk.Builder) (*gtk.Window ,error) {
  object, err := builder.GetObject("window1")
	if err != nil {
		return nil, err
	}
	ventana, ok := object.(*gtk.Window)
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
	if err != nil {
		panic(err)
	}
	treeView.AppendColumn(creaColumna("Title", COLUMN_TITLE))
	treeView.AppendColumn(creaColumna("Performer", COLUMN_PERFORMER))
  treeView.AppendColumn(creaColumna("Album", COLUMN_ALBUM))
	treeView.AppendColumn(creaColumna("Genre", COLUMN_GENRE))
	listStore, err := gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING)
	if err != nil {
		panic(err)
	}
	treeView.SetModel(listStore)
	return treeView, listStore
}

func actualizaLista(listStore *gtk.ListStore, canciones []string)  {
  i := 0
  for i < len(canciones) {
    nuevoRenglon(listStore,canciones[i],canciones[i+1],canciones[i+2],canciones[i+3])
    i += 4
  }
}

func nuevoRenglon(listStore *gtk.ListStore, titulo string , interprete string, album string, genero string) {
	iter := listStore.Append()
	err := listStore.Set(iter,
		[]int{COLUMN_TITLE, COLUMN_PERFORMER,COLUMN_ALBUM, COLUMN_GENRE},
		[]interface{}{titulo, interprete, album, genero})
	if err != nil {
		panic(err)
	}
}
