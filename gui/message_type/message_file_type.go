package messagetype

import (
	"fmt"
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

			savePath := filepath.Join(uri.Path(), filepath.Base(file.File_path))
			fmt.Println(file.HolderIP, "~", file.File_path, "~", savePath)

			startFileDownloadingChan <- fmt.Sprint(file.HolderIP, "~", file.File_path, "~", savePath)
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

func (f *File_type) CopyFileType(new_file *File_type) {
	f.a = new_file.a
	f.File_path = new_file.File_path
	f.HolderIP = new_file.HolderIP
	f.StartFileDownloadingChan = new_file.StartFileDownloadingChan

	f.SetText(filepath.Base(f.File_path))

	f.OnTapped = func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if uri == nil || err != nil {
				return
			}

			savePath := filepath.Join(uri.Path(), filepath.Base(f.File_path))

			f.StartFileDownloadingChan <- fmt.Sprint(f.HolderIP, "~", f.File_path, "~", savePath)
		}, f.a)
	}

	f.Refresh()
}
