package main

import (
	"fmt"
	"github.com/mohamedamer/byggscrap/models"
	"github.com/mohamedamer/byggscrap/scraper"
	"github.com/robfig/cron"
	"log"
	"net/http"
	"net/smtp"
	"net/url"
	"os"
	"os/signal"
	"strings"
)

func main() {
	c := cron.New()
	c.AddFunc("@every 10s", getApartmentInfo)
	c.Start()
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
	log.Print("send email")
		§//sendEmail("done")
}

func sendEmail(msg string) {
	from := ""
	pass := ""
	to := ""
	err := smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))
	if err != nil {
		log.Fatal(err)
	}
}
func createEmail(result models.Result) {
	if result.TotalCount > 0 {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("the following apartments are now available %v !"))
		for _, item := range result.Result {
			sb.WriteString(fmt.Sprintf("Address %v\nFloor %v\nArea %v\nDescription %v\nRent per month %v\nLatest application date %v",item.StreetName, item.ObjectFloor, item.ObjectArea, item.ObjectTypeDescription, item.RentPerMonth, item.EndPeriodMPDateString))
		}
		//msg := sb.String()
	}
}

func getApartmentInfo() {
	client := &http.Client{}
	scrapeExecutor := scraper.ScrapeExecutor{Client: client}
	request, err := createRequest()
	job := scraper.ScrapeJob{Request: request, Status: scraper.READY}
	result := scrapeExecutor.Scrape(&job)
	if err != nil {
		log.Fatalf("error occured %v", err)
	}
	fmt.Printf("results %v %v", result.TotalCount, result.Result)
}

func createBody() url.Values {
	values := make(url.Values)
	data := `{"CompanyNo":-1,"SyndicateNo":1,"ObjectMainGroupNo":1,"HouseForms":[{"No":5}],"Advertisements":[{"No":-1}],"ObjectSeekAreas":[{"No":1}],"RentLimit":{"Min":0,"Max":20000},"AreaLimit":{"Min":0,"Max":150},"ApplySearchFilter":true,"Page":1,"Take":10,"SortOrder":"","ReturnParameters":["ObjectNo","FirstEstateImageUrl","Street","SeekAreaDescription","PlaceName","ObjectSubDescription","ObjectArea","RentPerMonth","MarketPlaceDescription","CountInterest","FirstInfoTextShort","FirstInfoText","EndPeriodMP","FreeFrom","SeekAreaUrl","Latitude","Longitude","BoardNo"]}`
	values.Set("Parm1", data)
	values.Set("CallbackMethod", "PostObjectSearch")
	values.Set("CallbackParmCount", "1")
	values.Set("__WWEVENTCALLBACK", "")
	return values
}

func createRequest() (*http.Request, error) {
	values := createBody()
	request, err := http.NewRequest("POST", "https://marknad.byggvesta.se/API/Service/SearchServiceHandler.ashx", strings.NewReader(values.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("X-Requested-With", "XMLHttpRequest")
	request.Header.Set("Host", "marknad.byggvesta.se")
	request.Header.Set("Referer", "https://marknad.byggvesta.se/pgSearchResult.aspx")
	request.Header.Set("X-Momentum-API-KEY", "dIwHbOLgCS+FoZLYNYNToP9zK6VUoSgVC8BOT6cYljU=")
	request.Header.Set("Cookie", "ga=GA1.2.1142674423.1542053180; _gid=GA1.2.1223576119.1542053180; ASP.NET_SessionId=b1csjxwfo03ydyfxz3yoh1qy; Language=se")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36")
	request.Header.Set("Pragma", "no-cache")
	request.Header.Set("Origin", "https://marknad.byggvesta.se")
	request.Header.Set("Accept-Encoding", "gzip, deflate, br")
	request.Header.Set("Accept-Language", "en-US,en;q=0.9,ar;q=0.8,de;q=0.7")
	request.Header.Set("Accept", "application/data,text/*")
	request.Header.Set("Cache-Control", "no-cache")
	return request, err
}
