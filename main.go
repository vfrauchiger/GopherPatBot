package main

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func chooseFile(w fyne.Window, fname *widget.Entry) {
	dialog.ShowFileOpen(func(file fyne.URIReadCloser, err error) {
		if file == nil {
			return
		}

		fileP := file.URI().Path()

		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(fileP)
		fname.SetText(fileP)
		file.Close()
	}, w)

}

func main() {
	//new window
	a := app.New()
	w := a.NewWindow("Gopher PDF Bot")

	//elements

	//File handling
	labFileIn := widget.NewLabel("Input File: ")
	inpFileName := widget.NewEntry()
	butGetFileName := widget.NewButton("Choose File", func() {
		chooseFile(w, inpFileName)

	})

	//Label underneath buttons:
	labFinishedTask := widget.NewLabel("")

	//Start the action:
	butStartDownload := widget.NewButton("Let's Get the Stuff!", func() {
		fmt.Println("Yeah!")
		publList := loadExcTable(inpFileName.Text)
		labFinishedTask.SetText("Running!")
		for _, publno := range publList {
			publnoList := numberIngestion(publno)
			// for every document in the list a new go routine is started
			go getOnePublication(publnoList)
		}

		labFinishedTask.SetText("All Done!")
	})

	//Exit the program!
	butQuit := widget.NewButton("Quit", func() {
		os.Exit(1)
	})

	//content compilation
	content := container.NewVBox(
		widget.NewLabel("(c) 2021, Vinz Frauchiger"),
		labFileIn,
		inpFileName,
		container.NewHBox(
			butGetFileName,
			butStartDownload,
			butQuit,
		),

		labFinishedTask,
	)

	//content activation
	w.SetContent(content)
	w.ShowAndRun()
	w.SetOnClosed(func() {
		os.Exit(1)
	})
}
