package timedelement

type TimedElement struct {
	Timestamp int
	Data      string
}

type TimedElementRepo interface {
	Add(element *TimedElement) error
	GetAll() ([]*TimedElement, error)
	GetAllBeforeTimestamp(timestamp int64) ([]*TimedElement, error)
	Remove(element *TimedElement) error
	GetCurrentTimestamp() (*int64, error)
	UpdateTimestamp(timestamp int64) (*int64, error)
}

func HandleElementsBeforeTimestamp(r TimedElementRepo, timestamp int64, handle func(*TimedElement)) error {
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

func UpdateTimestampAndHandleElementsBeforeTimestamp(r TimedElementRepo, timestamp int64, handle func(*TimedElement)) error {
	latestTimestamp, err := r.UpdateTimestamp(timestamp)
	if err != nil {
		return err
	}
	return HandleElementsBeforeTimestamp(r, *latestTimestamp, handle)
}
