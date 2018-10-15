package main

import(
  "github.com/gotk3/gotk3/gtk"
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

  ventana.SetTitle("Reproductor-MP3")
	ventana.SetDefaultSize(500, 500)
	ventana.Connect("destroy", func ()  {
    gtk.MainQuit()
  })
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
  			return nil, error
  		}
  	}

  	return builder, nil
}


func window(builder *gtk.Builder) (*gtk.Window ,error) {

  object, err := b.GetObject("window1")
	if err != nil {
		return nil, err
	}

	ventana, ok := object.(*gtk.Window)
	if !ok {
		return nil, err
	}

	return ventana, nil
}

func grid()  {

}

func button()  {

}

func statusBar()  {

}
