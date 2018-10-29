package model

type ResponseBody struct {
	Name       string   `json:"name"`
	Version    string   `json:"version"`
	Vendor     string   `json:"vendor"`
	S3location string   `json:"s3location"`
	Files      []string `json:"files"`
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

type UploadPermission struct {
	Bucket     string `json:"bucket"`
	Region     string `json:"region"`
	AuthKey    string `json:"auth-key"`
	AuthSecret string `json:"auth-secret"`
	S3location string `json:"s3location"`
}