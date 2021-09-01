package objects

import (
	"time"
)

type PipelineStruct struct {
	Pipeline map[string]interface{} `json:""`
}

type PipelineList []struct {
	ID        int       `json:"id"`
	ProjectID int       `json:"project_id"`
	Sha       string    `json:"sha"`
	Ref       string    `json:"ref"`
	Status    string    `json:"status"`
	Source    string    `json:"source"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	WebURL    string    `json:"web_url"`
}

func (pl *PipelineStruct) ToTEXT(noHeaders bool) error {
	return nil
}
