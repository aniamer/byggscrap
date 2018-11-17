package scraper

import (
	"encoding/json"
	"fmt"
	"github.com/mohamedamer/byygscrap/models"
	"io/ioutil"
	"log"
	"net/http"
)

type ScrapeExecutor struct {
	Client *http.Client
}

type Status int

const(
	SUCCESS Status = 0
	PENDING Status = 1
	CANCLED Status = 2
	READY Status = 3
	FAIL Status = -1
)
type JobResult struct {
	Log string
	Content []byte
}
type ScrapeJob struct {
	Id string
	Recurrent bool
	Request *http.Request
	Result JobResult
	Status Status
}

type Scraper interface {
	Scrape(job ScrapeJob)
}

type Parser interface {
	Parse(interface{})error
}

func (jobResult *JobResult) Parse(contentStruct interface{}) error {
	jsonErr := json.Unmarshal(jobResult.Content, contentStruct)
	return jsonErr
}

func (executor *ScrapeExecutor) Scrape (job *ScrapeJob) models.Result {
	var result models.Result
	if job.Status == READY {
		log.Printf("job %v is ready and is starting ...", job.Id)
		response, err := executor.Client.Do(job.Request)
		defer response.Body.Close()
		if err != nil {
			job.Status= FAIL
		}

		if response.StatusCode >=200 && response.StatusCode <= 299 {
			job.Status = FAIL
		}

		if err != nil {
			fmt.Errorf("error occured %v", err)
		}

		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Errorf("error occured %v", err)
		}

		job.Result.Content = contents
		job.Result.Parse(&result)
	} else {
		log.Printf("job %v is not ready and is not ready for execution...", job.Id)
	}
	return result
}
