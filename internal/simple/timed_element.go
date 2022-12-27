package simple

type TimedElement struct {
	Timestamp int
	Data      string
}

type TimedElementRepo interface {
	Add(key string, element TimedElement)
	GetAll()
}
