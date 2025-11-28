package tls_communication

type mail_data struct {
	Package_type int
	Username     string
	FullAddress  string
	Message      string
	FilePath     string
}

const (
	PackageTypeHandshake    = 0
	PackageTypeMessage      = 1
	PackageTypeSendFileInfo = 2
	PackageTypeFileReq      = 3
)
