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

  ventana, _ := window(builder,"window1")
  ventanaEditar, _ := window(builder,"windowEdit")
  ventana2, _ := scrolledWindow(builder)
  ventanaAdvertencia, _ := window(builder, "windowAdvertencia")

  boton, _ := button(builder, "button1")

  botonEditar, _ := button(builder, "editTags")
  botonGuardar, _ := button(builder, "editTag")
  botonCancelarEditar, _ := button(builder, "cancelTag")

  botonMinero, _ := button(builder,"buttonMiner")
  botonMinar, _ := button(builder, "buttonMiner1")
  botonCancelarMinar, _:= button(builder, "cancelMin")


  grid, err := grid(builder)
  grid.SetOrientation(gtk.ORIENTATION_VERTICAL)
  if err != nil { panic(err) }
  busqueda, _ := entry(builder, "entry1")
  /* FACTORIZAR */
  titleEntry, performerEntry, albumEntry, genreEntry, err := entryEdit(builder)
  //if err != nil{ panic(err) }
  var rutaCancionSeleccionada  string
  treeView, listStore := creaTreeView()
  treeView.SetSearchEntry(busqueda)

  /* Cuando se quiere buscar una cancion */
  boton.Connect("clicked", func ()  {
    entrada, err := busqueda.GetText()
    if entrada != "" && err == nil {
      busqueda.SetText("")
      ok := Administrador.LeeEntrada(entrada)
      listStore.Clear()
      if ok {
        treeView.ExpandAll()
        renglones := Administrador.ObtenerRenglones()
        actualizaLista(listStore,renglones)
      }
    }
  })

  /* Cuando se presiona el boton de minar */
  botonMinar.Connect("clicked", func () {
    ventanaAdvertencia.ShowAll()
  })

  /* Cuando se cancela la opcion de minar en la advertencia */
  botonCancelarMinar.Connect("clicked", func () {
    ventanaAdvertencia.Hide()
  })

  /* Cuando se acepta la opcion de minar en el cuadro de advertencia */
  botonMinero.Connect("clicked", func ()  {
    Minero.Mina()
    rolas, err := Administrador.ObtenerBase()
    if err == nil {
      listStore.Clear()
      actualizaLista(listStore, rolas)
    }
    ventanaAdvertencia.Hide()
  })

  /* Ventana para requerir que se llenen todos los campos */
  windowEntryReq, _ := window(builder, "windowEntry")
  botonEntryReq, _ := button(builder, "buttonAceptEntry")

  /* Ventana para agregar una persona */
  windowPerson, _ := window(builder, "windowPerson")
  botonPerson, _ := button(builder, "buttonPerson1")
  botonAddPerson, _ := button(builder, "buttonAddPerson")
  botonCancelarPerson, _ := button(builder, "cancelPerson")

  /* Entradas de la ventana para agregar una persona */
  entryPersonNa, _ := entry(builder, "entryPersonStageName")
  entryPersonRn, _ := entry(builder, "entryPersonRealName")
  entryPersonBd, _ := entry(builder, "entryPersonBirthDate")
  entryPersonDd, _ := entry(builder, "entryPersonDeathDate")

  /* Se presiona el boton de "Agregar Interprete (Persona)" */
  botonPerson.Connect("clicked", func () {
    windowPerson.ShowAll()
  })
  /* Se presiona el boton de "Cancelar" en la opcion de agregar una Persona */
  botonCancelarPerson.Connect("clicked", func () {
    entryPersonNa.SetText("")
    entryPersonRn.SetText("")
    entryPersonBd.SetText("")
    entryPersonDd.SetText("")
    windowPerson.Hide()
  })
  /* Se guarda la nueva Persona creada */
  botonAddPerson.Connect("clicked", func (){
    person_na, _ := entryPersonNa.GetText()
    person_rn, _ := entryPersonRn.GetText()
    person_bd, _ := entryPersonBd.GetText()
    person_dd, _ := entryPersonDd.GetText()
    if (person_na != "" &&
        person_rn != "" &&
        person_bd != "" &&
        person_dd != "") {
      entryPersonNa.SetText("")
      entryPersonRn.SetText("")
      entryPersonBd.SetText("")
      entryPersonDd.SetText("")
      Administrador.InsertaInterpretePersona(person_na, person_rn, person_bd, person_dd);
      windowPerson.Hide()
    } else {
      //Lanzar ventana de autocompletar busqueda...
      windowEntryReq.ShowAll()
    }
  })

  /* Ventana para agregar un grupo */
  windowGroup, _ := window(builder, "windowGroup")
  /* Botones de la entrada para agregar un grupo */
  botonGroup, _ := button(builder, "buttonGroup1")
  botonAddGroup, _ := button(builder, "addGroup")
  botonCancelGroup, _ := button(builder, "cancelGroup")
  /* Las entradas de la ventana para agregar un grupo */
  entryGroupN, _ := entry(builder, "entryGroupName")
  entryGroupSd, _ := entry(builder, "entryGroupStartDate")
  entryGroupEd, _ := entry(builder, "entryGroupEndDate")

  /* Se presiona el boton de  "Agregar Interprete (Grupo)" */
  botonGroup.Connect("clicked", func () {
    windowGroup.ShowAll()
  })
  /* Se agrega el nuevo grupo */
  botonAddGroup.Connect("clicked", func () {
    group_n, _ := entryGroupN.GetText()
    group_sd, _ := entryGroupSd.GetText()
    group_ed, _ := entryGroupEd.GetText()
    if(group_n != "" &&
       group_sd != "" &&
       group_ed != ""){
      entryGroupN.SetText("")
      entryGroupSd.SetText("")
      entryGroupEd.SetText("")
      Administrador.InsetaInterpreteGrupo(group_n, group_sd, group_ed)
      windowGroup.Hide()
    } else {
      // Lanzar ventana para terminar de insertar un grupo
      windowEntryReq.ShowAll()
      }
  })

  /* Se cancela el agregar un nuevo grupo */
  botonCancelGroup.Connect("clicked", func (){
    entryGroupN.SetText("")
    entryGroupSd.SetText("")
    entryGroupEd.SetText("")
    windowGroup.Hide()
  })

  /* Se cierra la ventana de advertencia para llenar las entradas */
  botonEntryReq.Connect("clicked", func (){
    windowEntryReq.Hide()
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

  /* Cuando de cancela la accion de editar etiquetas */
  botonCancelarEditar.Connect("clicked", func () {
    ventanaEditar.Hide()
  })

  /* Cuando se presiona la opcion de editar una cancion */
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

func entry(builder *gtk.Builder, id string) (*gtk.Entry, error) {
  object, err := builder.GetObject(id)
  if err != nil {
    panic(err)
  }
  entry, ok := object.(*gtk.Entry)
  if !ok {
    return nil, err
  }
  return entry, nil
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
