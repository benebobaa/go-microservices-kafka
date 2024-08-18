package dto

type Status int

const (
	PENDING Status = iota
	IN_PROGRESS
	COMPLETE
	ERROR
	FAILED
)

func (s Status) String() string {
	return [...]string{"pending", "in_progress", "success", "error", "failed"}[s]
}
