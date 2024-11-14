package models

type Event interface {
	GetType() string
	GetObject() interface{}
}
