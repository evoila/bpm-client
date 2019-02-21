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
	Id           string
	Name         string       `json:"name"`
	Version      string       `json:"version"`
	Vendor       string       `json:"vendor"`
	FilePath     string       `json:"file_path"`
	UploadDate   string       `json:"upload_date"`
	Description  string       `json:"description"`
	Files        []string     `json:"files"`
	Stemcell     Stemcell     `json:"stemcell"`
	Dependencies []Dependency `json:"dependencies"`
	Size         int64        `json:"size"`
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

func (m MetaData) String2() string {

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
		stemcellString = m.Stemcell.stringFormat("")
	} else {
		stemcellString = "│  Stemcell:       Not specified"
	}

	return "│  Name:           " + m.Name + "\n" +
		"│  Version:        " + m.Version + "\n" +
		"│  Vendor:         " + m.Vendor + "\n" +
		"│  UploadDate:     " + m.UploadDate + "\n" +
		"│  Files:          " + formatStringArray(m.Files, "") +
		"│  Dependencies:   " + formatStringArray(dependenciesAsStrings, "") +
		stemcellString
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
		stemcellString = depth + "│  Stemcell:       Not specified \n"
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
	Name, Version, Vendor, Description string
	Stemcell                           Stemcell
	Files, Dependencies                []string
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
	Url            string         `yaml:"url"`
	Port           string         `yaml:"port"`
	Vendor         string         `yaml:"vendor"`
	Username       string         `yaml:"username"`
	Password       string         `yaml:"password"`
	KeycloakConfig KeycloakConfig `yaml:"keycloakConfig"`
}

type KeycloakConfig struct {
	Url      string `yaml:"url"`
	Realm    string `yaml:"realm"`
	ClientID string `yaml:"clientId"`
}

type PackageRequestBody struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Vendor  string `json:"vendor"`
}

type OpenId struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	IdToken          string `json:"id_token"`
	NotBeforePolicy  string `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}

type Page struct {
	Size          int `json:"size"`
	TotalElements int `json:"totalElements"`
	TotalPages    int `json:"totalPages"`
	Number        int `json:"number"`
}

type Embedded struct {
	Packages []MetaData `json:"packages"`
}

type PaginatedMetaData struct {
	Embedded Embedded `json:"_embedded"`
	Page     Page     `json:"page"`
}
