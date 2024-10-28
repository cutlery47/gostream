package repo

type InFile struct {
	Name string
	Size int
	// location of the file in s3
	Location string
}

type InVideo struct {
	File      InFile
	VideoName string
}

type File struct {
	Id   int
	Data InFile
}
