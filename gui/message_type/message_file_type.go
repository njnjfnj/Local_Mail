package messagetype

import "fyne.io/fyne/v2/widget"

type File_type struct {
	*widget.Button
	File_path  string
	File_name  string
	Holdername string
}

func New_file(new_file_path, new_file_name, new_holdername string) *File_type {
	// TODO: сделать нахождения файла, если не находиться выкидать ошибку
	// сохранить путь к файлу в локальной базе данных.
	file := &File_type{
		File_path:  new_file_path,
		File_name:  new_file_name,
		Holdername: new_holdername,
	}

	button := widget.NewButton(new_file_name, func() {
		// указать логику нахождения айпи владельца за юзернеймом
		// в сети, выбор где сохранить файл и собственно зашифрованая
		// передача файла
	})

	file.Button = button

	return file
}

func New_nill_file() *File_type {
	file := &File_type{
		File_path:  "",
		File_name:  "",
		Holdername: "",
	}

	button := widget.NewButton("", func() {})

	file.Button = button

	file.Button.Hide()

	return file
}
