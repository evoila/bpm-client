package model

import (
	"strconv"
	"strings"
)

type ResponseBody struct {
	Name       string   `json:"name"`
	Version    string   `json:"version"`
	Vendor     string   `json:"vendor"`
	S3location string   `json:"s3location"`
	Files      []string `json:"files"`
}

type MetaData struct {
	Name, Version, Vendor, FilePath, UploadDate string
	Files                                       []string
	Stemcell                                    Stemcell
	Dependencies                                []Dependency
}

type Dependency struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Vendor  string `json:"vendor"`
}

type Stemcell struct {
	Family       string `json:"family"`
	MajorVersion int    `json:"major_version"`
	MinorVersion int    `json:"minor_version"`
}

func (m MetaData) String(depth string) string {

	var dependenciesAsStrings []string

	if len(m.Dependencies) > 0 {
		for _, d := range m.Dependencies {
			dependenciesAsStrings = append(dependenciesAsStrings, d.String())
		}
	} else {
		dependenciesAsStrings = append(dependenciesAsStrings, "none")
	}

	var stemcellString string
	if (Stemcell{}) != m.Stemcell {
		stemcellString = m.Stemcell.stringFormat(depth)
	} else {
		stemcellString = depth + "│  Stemcell:       Not specified"
	}

	return depth + "│  Name:           " + m.Name + "\n" +
		depth + "│  Version:        " + m.Version + "\n" +
		depth + "│  Vendor:         " + m.Vendor + "\n" +
		depth + "│  UploadDate:     " + m.UploadDate + "\n" +
		depth + "│  Files:          " + formatStringArray(m.Files, depth) +
		depth + "│  Dependencies:   " + formatStringArray(dependenciesAsStrings, depth) +
		stemcellString +
		depth + "│"
}

func (s Stemcell) isNotEmpty() bool {
	return len(s.Family) > 0
}

func (s Stemcell) stringFormat(depth string) string {

	return depth + "│  Stemcell:       " + "\n" +
		depth + "│   Family:        " + s.Family + "\n" +
		depth + "│   Major Version: " + strconv.Itoa(s.MajorVersion) + "\n" +
		depth + "│   Minor Version: " + strconv.Itoa(s.MinorVersion) + "\n"
}

func formatStringArray(stringArray []string, depth string) string {

	if len(stringArray) > 1 {
		return strings.Join(stringArray, "\n"+depth+"│                  ") + "\n"
	} else {
		return strings.Join(stringArray, "") + "\n"
	}
}

func (d Dependency) String() string {
	return d.Name + ":" + d.Version + " by " + d.Vendor
}

type SpecFile struct {
	Name, Version       string
	Stemcell            Stemcell
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
	Bucket       string `json:"bucket"`
	Region       string `json:"region"`
	AuthKey      string `json:"auth-key"`
	AuthSecret   string `json:"auth-secret"`
	S3location   string `json:"s3location"`
	SessionToken string `json:"session-token"`
}

type Config struct {
	Url, Port, Vendor string
}

type PackageRequestBody struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Vendor  string `json:"vendor"`
}
