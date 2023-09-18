package models

type Worker struct {
	ID                string
	KeepAliveDeadline uint64
	Online            bool
	Error             string
}
