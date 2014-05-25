package lds

type BEvent struct {
	Fname string
	Etype int
	Etypestring string
}

type NBEvent struct {
	Fname string
	Etype int
	Etypestring string
	Remote string
}

type FileRequest struct{
	Fname string
}

type BFile struct {
	Name string
	Data []byte
}
