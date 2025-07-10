package persistance

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"web-page-analyzer/model"
)

func TestInMemoryStore_CreateAndGetJob(t *testing.T) {
	store := NewMemoryStore()
	url := "https://test.com"

	job := store.CreateJob(url)

	assert.NotEmpty(t, job.ID)
	assert.Equal(t, "PROCESSING", job.Status)

	storedJob, _ := store.GetJob(job.ID)

	assert.Equal(t, job.ID, storedJob.ID)
	assert.Equal(t, "PROCESSING", storedJob.Status)
}

func TestInMemoryStore_CompleteJob_Success(t *testing.T) {
	store := NewMemoryStore()
	url := "https://test.com"

	job := store.CreateJob(url)

	result := model.AnalysisResult{
		HTMLVersion: "HTML5",
	}

	store.CompleteJob(job.ID, result, nil)

	updatedJob, ok := store.GetJob(job.ID)

	assert.True(t, ok)
	assert.Equal(t, "COMPLETED", updatedJob.Status)
	assert.NotNil(t, updatedJob.Result)
	assert.Equal(t, "HTML5", updatedJob.Result.HTMLVersion)
}

func TestInMemoryStore_CompleteJob_Failure(t *testing.T) {
	store := NewMemoryStore()
	url := "https://test.com"

	job := store.CreateJob(url)

	simErr := errors.New("something went wrong")
	store.CompleteJob(job.ID, model.AnalysisResult{}, simErr)

	updatedJob, ok := store.GetJob(job.ID)

	assert.True(t, ok)
	assert.Equal(t, "FAILED", updatedJob.Status)
	assert.Equal(t, simErr.Error(), updatedJob.Error)
}

func TestInMemoryStore_GetJob_NotFound(t *testing.T) {
	store := NewMemoryStore()

	job, ok := store.GetJob("non-existent-id")
	assert.False(t, ok)
	assert.Nil(t, job)
}

func TestInMemoryStore_ListJobs(t *testing.T) {
	store := NewMemoryStore()

	store.CreateJob("https://test1.com")
	store.CreateJob("https://test2.com")

	jobs := store.ListJobs()

	assert.Len(t, jobs, 2)
}
