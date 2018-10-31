package model

import "strings"

type ResponseBody struct {
	Name       string   `json:"name"`
	Version    string   `json:"version"`
	Vendor     string   `json:"vendor"`
	S3location string   `json:"s3location"`
	Files      []string `json:"files"`
}

type MetaData struct {
	Name, Version, Vendor, FilePath string
	Files, Dependencies             []string
}

func (m MetaData) String() string {

	return "Name:         " + m.Name + "\n" +
		"Version:      " + m.Version + "\n" +
		"Vendor:       " + m.Vendor + "\n" +
		"Files:        " + strings.Join(m.Files, "\n              ") + "\n" +
		"Dependencies: " + strings.Join(m.Dependencies, "\n              ") + "\n"
}

type SpecFile struct {
	Name, Version       string
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

type S3Permission struct {
	Bucket     string `json:"bucket"`
	Region     string `json:"region"`
	AuthKey    string `json:"auth-key"`
	AuthSecret string `json:"auth-secret"`
	S3location string `json:"s3location"`
}

type Config struct {
	Url, Port, Vendor string
}

type PackageRequestBody struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Vendor  string `json:"vendor"`
}
