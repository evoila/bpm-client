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
