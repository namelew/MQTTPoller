package output

type Worker struct {
	Id      int
	NetId   string
	Online  bool
	History []interface{}
}