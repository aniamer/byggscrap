package scraper

import "net/http"

type ScrapeExecutor struct {
	Client http.Client
	Url string
}

type Status int

const(
	SUCCESS Status = 0
	PENDING Status = 1
	CANCLED Status = 2
	FAIL Status = -1
)

type ScrapeJob struct {
	id string
	Recurrent bool
	Request http.Request
	Status Status
}

type Scraper interface {
	Scrape(job ScrapeJob)
}
