package messagetype

import (
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type File_type struct {
	*widget.Button
	File_path                string
	HolderIP                 string
	StartFileDownloadingChan chan string
	a                        fyne.Window
}

func New_file(new_file_path, new_HolderIP string, a fyne.Window, startFileDownloadingChan chan string) *File_type {
	file := &File_type{
		File_path:                new_file_path,
		HolderIP:                 new_HolderIP,
		StartFileDownloadingChan: startFileDownloadingChan,
		a:                        a,
	}

	button := widget.NewButton(filepath.Base(file.File_path), func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if uri == nil || err != nil {
				return
			}

			savePath := filepath.Join(uri.Path(), "share", filepath.Base(file.File_path))

			startFileDownloadingChan <- new_HolderIP + "~" + savePath
		}, a)
	})

	file.Button = button

	return file
}

func New_nill_file() *File_type {
	file := &File_type{
		File_path:                "",
		HolderIP:                 "",
		StartFileDownloadingChan: nil,
	}

	button := widget.NewButton("", func() {})

	file.Button = button

	file.Button.Hide()

	return file
}

func CopyFileType(new_file *File_type) *File_type {
	file := &File_type{
		File_path:                new_file.File_path,
		HolderIP:                 new_file.HolderIP,
		StartFileDownloadingChan: new_file.StartFileDownloadingChan,
		a:                        new_file.a,
	}

	button := widget.NewButton(filepath.Base(file.File_path), func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if uri == nil || err != nil {
				return
			}

			savePath := filepath.Join(uri.Path(), "share", filepath.Base(file.File_path))

			file.StartFileDownloadingChan <- file.HolderIP + "~" + savePath
		}, file.a)
	})

	file.Button = button

	return file
}
