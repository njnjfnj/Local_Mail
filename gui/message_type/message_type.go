package messagetype

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

type Message_type struct {
	Text       *widget.Label
	Image      *canvas.Image
	File       *File_type
	Holdername string
}

func New_message(holdername, new_text, new_file_path, new_image_path string, a fyne.Window, startFileDownloadingChan chan string) *Message_type {
	m := &Message_type{}

	m.Holdername = holdername

	m.Text = widget.NewLabel(new_text)
	if new_text == "" {
		m.Text.Hide()
	}

	if new_file_path != "" {
		m.File = New_file(new_file_path, holdername, a, startFileDownloadingChan)
		// проверка на ошибку нахождения файла, в случае ошибки выкинуть ошибку
	} else {
		m.File = New_nill_file()
	}

	if new_image_path != "" {
		new_image := canvas.NewImageFromFile(new_image_path)
		m.Image = new_image
	} else {
		m.Image = canvas.NewImageFromImage(nil)
		m.Image.Hide()
	}

	return m
}

func New_text_message(holdername, new_text string) *Message_type {
	m := &Message_type{}

	m.Holdername = holdername

	m.Text = widget.NewLabel(new_text)
	if new_text == "" {
		m.Text.Hide()
	}

	return m
}
