package model

import "time"

/*
	Analysis Result of the url
*/
type AnalysisResult struct {
	HTMLVersion    string         `json:"html_version"`
	Title          string         `json:"title"`
	Headings       map[string]int `json:"headings"`
	Links          map[string]int `json:"links"`
	LoginFormFound bool           `json:"login_form_found"`
}

/*
	Job status
*/
type Job struct {
	ID        string          `json:"job_id"`
	Status    string          `json:"status"`
	Result    *AnalysisResult `json:"result,omitempty"`
	Error     error           `json:"error,omitempty"`
	CraetedAt time.Time       `json:"created_at"`
}
