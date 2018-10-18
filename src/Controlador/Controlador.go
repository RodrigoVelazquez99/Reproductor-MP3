package main

import(
  "github.com/gotk3/gotk3/gtk"
  "errors"
  "fmt"
  "github.com/gotk3/gotk3/glib"
)

const path = "../Interfaz/interfaz.glade"

func main()  {
  gtk.Init(nil)
  builder,err := build(path)
  if err != nil {
    panic(err)
  }
  ventana, err := window(builder)
  if err != nil {
    panic(err)
  }
  boton, err := button(builder)
  if err != nil {

    panic(err)
  }
  grid, err := grid(builder)
  grid.SetOrientation(gtk.ORIENTATION_VERTICAL)
  if err != nil  {
    panic(err)
  }
  //entry, err := entry(builder)
  treeView, _ := creaTreeView()
  boton.Connect("clicked", func ()  {
      fmt.Println("test")
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

func button(builder *gtk.Builder) (*gtk.Button, error)  {
  object, err := builder.GetObject("button1")
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

const (
	COLUMN_VERSION = iota
	COLUMN_FEATURE
)

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
	treeView.AppendColumn(creaColumna("Title", COLUMN_VERSION))
	treeView.AppendColumn(creaColumna("Performer", COLUMN_FEATURE))
  treeView.AppendColumn(creaColumna("Album", COLUMN_VERSION))
	treeView.AppendColumn(creaColumna("Genre", COLUMN_FEATURE))
	listStore, err := gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_STRING)
	if err != nil {
		panic(err)
	}
	treeView.SetModel(listStore)
	return treeView, listStore
}

func nuevoRenglon(listStore *gtk.ListStore, version, feature string) {
	iter := listStore.Append()
	err := listStore.Set(iter,
		[]int{COLUMN_VERSION, COLUMN_FEATURE},
		[]interface{}{version, feature})

	if err != nil {
		panic(err)
	}
}
