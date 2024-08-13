package dto

type Status int

const (
	PENDING Status = iota
	IN_PROGRESS
	COMPLETE
	FAILED
)

func (s Status) String() string {
	return [...]string{"PENDING", "IN_PROGRESS", "SUCCESS", "FAILED"}[s]
}
