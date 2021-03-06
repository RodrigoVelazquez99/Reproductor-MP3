package main

import(
  "fmt"
  "log"
	"os"
	"time"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
  "strings"
  "github.com/gotk3/gotk3/gtk"
  "errors"
  "github.com/gotk3/gotk3/glib"
  "github.com/gotk3/gotk3/gdk"
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

/* Variables para controlar la reproduccion de musica */
var streamer beep.StreamSeekCloser
var format beep.Format
var ctrl *beep.Ctrl
/* La barra de musica */
var progressBar *gtk.ProgressBar

func main()  {
  Administrador.IniciaBase()
  gtk.Init(nil)
  builder := build(path)

  ventana := window(builder,"window1")
  ventanaEditar := window(builder,"windowEdit")
  ventana2 := scrolledWindow(builder)
  ventanaAdvertencia := window(builder, "windowAdvertencia")

  /* La portada de la cancion actual */
  imagenCancion := imagen(builder, "songImage")
  imagenP, err := gdk.PixbufNewFromFile("/home/rodrigofvc/go/src/github.com/RodrigoVelazquez99/Reproductor-MP3/src/Imagenes/default_image.jpg")
  if err != nil {
    /* La imagen por default no existe*/
    errors.New("La imagen por default no existe")
  }
  imagenCancion.SetFromPixbuf(imagenP)
  imagenCancion.SetPixelSize(200)

  boton := button(builder, "button1")

  botonEditar := button(builder, "editTags")
  botonGuardar := button(builder, "editTag1")
  botonCancelarEditar := button(builder, "cancelTag1")

  botonMinero := button(builder,"buttonMiner")
  botonMinar := button(builder, "buttonMiner1")
  botonCancelarMinar:= button(builder, "cancelMin")

  grid := grid(builder)
  grid.SetOrientation(gtk.ORIENTATION_VERTICAL)

  busqueda := entry(builder, "entry1")

  /* Botones para reproducir canciones */
  botonPlay := button(builder, "buttonPlay")
  botonPause := button(builder, "buttonPause")
  botonNext := button(builder, "buttonRight")
  botonPrev := button(builder, "buttonLeft")
  //pause := false
  progressBar = progress(builder, "progressBar")

  //botonLeft := button(builder, "buttonLeft")
  //botonRight := button(builder, "buttonRight")

  /* Entradas para la opcion de editar las etiquetas de una cancion */
  titleEntry := entry (builder, "setTitle")
  performerEntry := entry(builder, "setPerformer")
  albumEntry := entry(builder, "setAlbum")
  genreEntry := entry(builder, "setGenre")

  /* Para poder abrir una ventana y escoger el archivo para imagen de la cancion */
  botonChooser := buttonChoose(builder, "buttonChooser")

  /* La ruta de la cancion que ha sido seleccionada */
  var rutaCancionSeleccionada  string

  treeView, listStore := creaTreeView()
  treeView.SetSearchEntry(busqueda)

  /* Cuando se quiere buscar una cancion */
  boton.Connect("clicked", func () {
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
  windowEntryReq := window(builder, "windowEntry")
  botonEntryReq := button(builder, "buttonAceptEntry")

  /* Ventana para agregar una persona */
  windowPerson := window(builder, "windowPerson")
  botonPerson := button(builder, "buttonPerson1")
  botonAddPerson := button(builder, "buttonAddPerson")
  botonCancelarPerson := button(builder, "cancelPerson")

  /* Entradas de la ventana para agregar una persona */
  entryPersonNa := entry(builder, "entryPersonStageName")
  entryPersonRn := entry(builder, "entryPersonRealName")
  entryPersonBd := entry(builder, "entryPersonBirthDate")
  entryPersonDd := entry(builder, "entryPersonDeathDate")

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
  windowGroup := window(builder, "windowGroup")
  /* Botones de la entrada para agregar un grupo */
  botonGroup := button(builder, "buttonGroup1")
  botonAddGroup := button(builder, "addGroup")
  botonCancelGroup := button(builder, "cancelGroup")
  /* Las entradas de la ventana para agregar un grupo */
  entryGroupN := entry(builder, "entryGroupName")
  entryGroupSd := entry(builder, "entryGroupStartDate")
  entryGroupEd := entry(builder, "entryGroupEndDate")

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
    nuevoTitulo, _ := titleEntry.GetText()
    nuevoInterprete, _ := performerEntry.GetText()
    nuevoAlbum, _ := albumEntry.GetText()
    nuevoGenero, _ := genreEntry.GetText()
    nuevaImagen := botonChooser.GetFilename()
    if nuevaImagen == "" {
      // No se escogio una imagen
      nuevaImagen = "/home/rodrigofvc/go/src/github.com/RodrigoVelazquez99/Reproductor-MP3/src/Imagenes/default_image.jpg"
    }
    /* Revisamos si el archivo es una imagen */
    if !(strings.HasSuffix(nuevaImagen, ".jpg") ||
        strings.HasSuffix(nuevaImagen, ".png")) {
      windowEntryReq.ShowAll()
    } else {
      Administrador.CambiaEtiquetas(rutaCancionSeleccionada, nuevoTitulo, nuevoInterprete, nuevoAlbum, nuevoGenero, nuevaImagen)
      listStore.Clear()
      renglones := Administrador.ObtenerRenglones()
      actualizaLista(listStore, renglones)
      ventanaEditar.Hide()
    }
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
  /* Se selecciona una cancion del TreeView */
  seleccion.Connect("changed", func ()  {
    model, iter, ok := seleccion.GetSelected()
    if ok {
      columna4,_ := model.(*gtk.TreeModel).GetValue(iter,4)
      ruta,_ := columna4.GetString()
      rutaCancionSeleccionada = ruta
      etiquetasCancionSeleccionada := Administrador.BuscaPorRuta(ruta)
      imagenP, err = gdk.PixbufNewFromFile(etiquetasCancionSeleccionada[4])
      if err != nil {
        /* La cancion no tiene imagen */
        imagenP, _ = gdk.PixbufNewFromFile("/home/rodrigofvc/go/src/github.com/RodrigoVelazquez99/Reproductor-MP3/src/Imagenes/default_image.jpg")
      }
      imagenCancion.SetFromPixbuf(imagenP)
      imagenCancion.SetPixelSize(200)
      }
  })

  /* Se presiona el boton para reproducir una cancion */
  botonPlay.Connect("clicked", func () {
    progressBar.SetFraction(0.00)
    go status()
    go play(rutaCancionSeleccionada)
  })

  // Se presiona el boton de pausa
  botonPause.Connect("clicked", func () {
      speaker.Lock()
      ctrl.Paused = !ctrl.Paused
      speaker.Unlock()
  })

  // Se presiona el boton de siguiente
  botonNext.Connect("clicked", func () {
    renglones := Administrador.ObtenerRenglones()
    i := 4
    for i < len(renglones){
      if (renglones[i] == rutaCancionSeleccionada) {
        if (len(renglones) > i + 5){
          rutaCancionSeleccionada = renglones[i+5]
          go play(rutaCancionSeleccionada)
        }
        break
      }
      i += 5
    }
  })


  // Se presiona el boton de anterior
  botonPrev.Connect("clicked", func () {
    renglones := Administrador.ObtenerRenglones()
    i := 4
    for i < len(renglones){
      if (renglones[i] == rutaCancionSeleccionada) {
        if (i - 5 >= 4){
          rutaCancionSeleccionada = renglones[i-5]
          go play(rutaCancionSeleccionada)
        }
        break
      }
      i += 5
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

/* Controla el movimiento de la barra de progreso */
func status() {
  i := 0.0
  for i <= 1.0 {
    time.Sleep(1 * time.Second)
    i += .01
    progressBar.SetFraction(i)
  }
}

func play(rutaCancionSeleccionada string) {
    f, err := os.Open(rutaCancionSeleccionada)
    if err != nil { log.Fatal(err) }
    streamer, format, err = mp3.Decode(f)
    if err != nil { log.Fatal(err) }
    defer streamer.Close()

    speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
    ctrl = &beep.Ctrl{Streamer: beep.Loop(-1, streamer), Paused: false}
    speaker.Play(ctrl)

    	for {
    		fmt.Println("Reproduciendo...")
    		fmt.Scanln()
    	}
}



func build(ruta string) *gtk.Builder  {
  	builder, err := gtk.BuilderNew()
  	if err != nil {
  		panic (err)
  	}
  	if ruta != "" {
  		err = builder.AddFromFile(ruta)
  		if err != nil {
		      errors.New("Ocurrio un error")
  		}
  	}
  	return builder
}

/* Crea una ventana que cumple con el identificador*/
func window(builder *gtk.Builder, tipo string) *gtk.Window {
  object, err := builder.GetObject(tipo)
	if err != nil {
		panic(err)
	}
	ventana, ok := object.(*gtk.Window)
	if !ok {
		errors.New("Ocurrio un error")
	}
	return ventana
}

func scrolledWindow(builder *gtk.Builder) *gtk.ScrolledWindow {
  object, err := builder.GetObject("window2")
	if err != nil {
		panic(err)
	}
	ventana, ok := object.(*gtk.ScrolledWindow)
	if !ok {
		errors.New("Ocurrio un error")
	}
	return ventana
}

func dialog(builder *gtk.Builder) *gtk.Dialog  {
  object, err := builder.GetObject("windowEdit")
	if err != nil {
		panic(err)
	}
	ventana, ok := object.(*gtk.Dialog)
	if !ok {
    errors.New("Ocurrio un error")
	}
	return ventana
}


func button(builder *gtk.Builder, tipo string) *gtk.Button  {
  object, err := builder.GetObject(tipo)
  if err != nil {
    panic(err)
  }
  boton, ok := object.(*gtk.Button)
  if !ok {
    errors.New("Ocurrio un error")
  }
  return boton
}

func progress(builder *gtk.Builder, tipo string) *gtk.ProgressBar {
  object, err := builder.GetObject(tipo)
  if err != nil {
    panic(err)
  }
  p, ok := object.(*gtk.ProgressBar)
  if !ok {
    errors.New("Ocurrio un error")
  }
  return p
}



func grid(builder *gtk.Builder) *gtk.Grid  {
  object, err := builder.GetObject("grid1")
  if err != nil {
    panic(err)
  }
  grid, ok := object.(*gtk.Grid)
  if !ok {
    errors.New("Ocurrio un error")
  }
  return grid
}

func entry(builder *gtk.Builder, id string) *gtk.Entry {
  object, err := builder.GetObject(id)
  if err != nil {
    panic(err)
  }
  entry, ok := object.(*gtk.Entry)
  if !ok {
		errors.New("Ocurrio un error")
  }
  return entry
}

func imagen(builder *gtk.Builder, id string) *gtk.Image{
  object, err := builder.GetObject(id)
  if err != nil { panic (err) }
  image, ok := object.(*gtk.Image)
  if !ok {
    errors.New("Ocurrio un error")
  }
  return image
}

func buttonChoose(builder *gtk.Builder, id string) *gtk.FileChooserButton {
  object, err := builder.GetObject(id)
  if err != nil { panic (err) }
  fileChooser, ok := object.(*gtk.FileChooserButton)
  if !ok {
    errors.New("Ocurrio un error")
  }
  return fileChooser
}

func creaColumna(nombre string, id int) *gtk.TreeViewColumn {
	cellRenderer, err := gtk.CellRendererTextNew()
	if err != nil {
		panic(err)
	}
	columna, err := gtk.TreeViewColumnNewWithAttribute(nombre, cellRenderer, "text", id)
	if err != nil {
		errors.New("Ocurrio un error")
	}
	return columna
}

func creaTreeView() (*gtk.TreeView, *gtk.ListStore) {
	treeView, err := gtk.TreeViewNew()
	if err != nil { panic(err) }
	treeView.AppendColumn(creaColumna("Titulo", COLUMN_TITLE))
	treeView.AppendColumn(creaColumna("Artista", COLUMN_PERFORMER))
  treeView.AppendColumn(creaColumna("Album", COLUMN_ALBUM))
	treeView.AppendColumn(creaColumna("Genero", COLUMN_GENRE))
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
