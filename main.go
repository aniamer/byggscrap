package main

import (
	"bytes"
	"fmt"
	"github.com/mohamedamer/byggscrap/models"
	"github.com/mohamedamer/byggscrap/scraper"
	"github.com/mailgun/mailgun-go"
	"github.com/robfig/cron"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
)

func main() {
	signalChannel := make(chan os.Signal)
	c := cron.New()
	c.AddFunc("@every 5s", byggJob)
	c.Start()
	signal.Notify(signalChannel, os.Interrupt, os.Kill)
	<-signalChannel
}

func byggJob() {
	from := "mo@gozyra.com"
	to := "m.elsayedamer@gmail.com"
	fetchAndNotify(sendEmail, from, to)
}

func fetchAndNotify(fn func(msg string, from string, to string), from string, to string) {
	result := getApartmentInfo()
	if (result.TotalCount > 0) {
		email := createEmail(*result)
		log.Print("send email")

		fn(email, from, to)
	}
}

func sendEmail(msg string, from string, to string) {
	//postmaster := "postmaster@mg.gozyra.com"
	mailgun := mailgun.NewMailgun("mg.gozyra.com", apiKey)
	message := mailgun.NewMessage(from, "New apartments offered", "", to)
	message.SetHtml(msg)
	_, _, err := mailgun.Send(message)
	if err != nil {
		log.Printf("error sending message\n",err)
	}
}

func createEmail(result models.Result) string {
	t, err := template.ParseFiles("./templates/alert.html")
	if err != nil {
		log.Printf("error occured loading email template", err)
	}
	buf := new(bytes.Buffer)

	if err = t.Execute(buf, result); err != nil {
		log.Printf("error occured while rendering email template", err)
	}

	return buf.String()
	//var sb strings.Builder
	//sb.WriteString(fmt.Sprint("the following apartments are now available:\n"))
	//
	//for _, item := range result.Result {
	//	sb.WriteString(fmt.Sprintf("Address %v\nFloor %v\nArea %v\nDescription %v\nRent per month %v\nLatest application date %v\n",item.StreetName, item.ObjectFloor, item.ObjectArea, item.ObjectTypeDescription, item.RentPerMonth, item.EndPeriodMPDateString))
	//
	//}
	//msg := sb.String()
	//return msg
}

func getApartmentInfo() *models.Result {
	client := &http.Client{}
	scrapeExecutor := scraper.ScrapeExecutor{Client: client}
	request, err := createRequest()
	job := scraper.ScrapeJob{Request: request, Status: scraper.READY}
	result := scrapeExecutor.Scrape(&job)
	if err != nil {
		log.Printf("error occured %v", err)
	}

	fmt.Printf("results %v %v\n", result.TotalCount, result.Result)
	return &result
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
