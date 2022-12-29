package timedelement

type TimedElement struct {
	Timestamp int
	Data      string
}

type TimedElementRepo interface {
	Add(element *TimedElement) error
	GetAll() ([]*TimedElement, error)
	GetAllBeforeTimestamp(timestamp int) ([]*TimedElement, error)
	Remove(element *TimedElement) error
	GetCurrentTimestamp() (*int64, error)
	UpdateTimestamp(timestamp int64) error
}

func HandleElementsBeforeTimestamp(r TimedElementRepo, timestamp int, handle func(*TimedElement)) error {
	all, err := r.GetAllBeforeTimestamp(timestamp)
	if err != nil {
		return err
	}
	for _, timedElement := range all {
		handle(timedElement)
		err := r.Remove(timedElement)
		if err != nil {
			return err
		}
	}
	return nil
}
