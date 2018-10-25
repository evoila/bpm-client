package model

type ResponseBoshPackage struct {
	Id, Name, Spec, Packaging, Version string
	Blobs, Dependencies                []string
}

type BoshPackage struct {
	id, name, spec, packaging, version string
	blobs                              []BoshBlob
	dependencies                       []BoshPackage
}

type BoshBlob struct {
	Id, Name, Version string
}

type MetaData struct {
	Name, Version, Vendor, FilePath string
	Files                           []string
}

type SpecFile struct {
	Name                string
	Files, Dependencies []string
}

type BackupResponse struct {
	Message  string `json:"message"`
	FileName string `json:"filename"`
	Region   string `json:"region"`
	Bucket   string `json:"bucket"`
}

type ErrorResponse struct {
	Message      string `json:"message"`
	State        string `json:"state"`
	ErrorMessage string `json:"error message"`
}

type Destination struct {
	Type       string
	Bucket     string
	Region     string
	AuthKey    string
	AuthSecret string
	File       string
}

type DbInformation struct {
	Host       string
	User       string
	Password   string
	Database   string
	Parameters []map[string]interface{}
}
