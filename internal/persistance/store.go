package persistance

import "web-page-analyzer/model"

/*
	interface to underline store
	for demo purpose this is in-memory, but ideally this should be reliable,distributed storage like Redis,DB
*/
type Store interface {
	CreateJob(url string) model.Job
	CompleteJob(id string, result model.AnalysisResult, err error)
	GetJob(id string) (*model.Job, bool)
	ListJobs() []*model.Job
}
