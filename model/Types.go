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
	Files                           []string
	Dependencies                    []Dependency
}

type Dependency struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Vendor  string `json:"vendor"`
}

func (m MetaData) String(depth string) string {

	var depsAsStrings []string

	if len(m.Dependencies) > 0 {
		for _, d := range m.Dependencies {
			depsAsStrings = append(depsAsStrings, d.String())
		}
	} else {
		depsAsStrings = append(depsAsStrings, "none")
	}

	if len(depsAsStrings) > 1 {
		return depth + "|	Name:         " + m.Name + "\n" +
			depth + "|	Version:      " + m.Version + "\n" +
			depth + "|	Vendor:       " + m.Vendor + "\n" +
			depth + "|	Files:        " + strings.Join(m.Files, "\n              ") + "\n" +
			depth + "|	Dependencies: " + strings.Join(depsAsStrings, depth+"\n|	              ") +
			depth + "\n|"
	} else {
		return depth + "|	Name:         " + m.Name + "\n" +
			depth + "|	Version:      " + m.Version + "\n" +
			depth + "|	Vendor:       " + m.Vendor + "\n" +
			depth + "|	Files:        " + strings.Join(m.Files, "\n              ") + "\n" +
			depth + "|	Dependencies: " + strings.Join(depsAsStrings, "\n|") +
			depth + "\n|"
	}
}

func (d Dependency) String() string {
	return d.Name + ":" + d.Version + " by " + d.Vendor
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
