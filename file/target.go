package file

import "time"

const (
	CloudFormation = "cfm"
	CloudWatchLogs = "cwl"
)

type Target struct {
	ID      *string
	Name    *string
	Region  string
	Service string
	Ref     *time.Time
}

func (t *Target) Parts() []string {
	return []string{
		t.Service,
		t.Region,
		*t.Name,
	}
}
