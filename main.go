package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	client := &http.Client{}
	values := make(url.Values)
	json := `{"CompanyNo":-1,"SyndicateNo":1,"ObjectMainGroupNo":1,"Advertisements":[{"No":-1}],"ObjectSeekAreas":[{"No":1}],"RentLimit":{"Min":0,"Max":20000},"AreaLimit":{"Min":0,"Max":150},"ApplySearchFilter":true,"Page":1,"Take":10,"SortOrder":"","ReturnParameters":["ObjectNo","FirstEstateImageUrl","Street","SeekAreaDescription","PlaceName","ObjectSubDescription","ObjectArea","RentPerMonth","MarketPlaceDescription","CountInterest","FirstInfoTextShort","FirstInfoText","EndPeriodMP","FreeFrom","SeekAreaUrl","Latitude","Longitude","BoardNo"]}`
	values.Set("Parm1", json)
	values.Set("CallbackMethod","PostObjectSearch")
	values.Set("CallbackParmCount", "1")
	values.Set("__WWEVENTCALLBACK", "")
	request, err := http.NewRequest("POST", "https://marknad.byggvesta.se/API/Service/SearchServiceHandler.ashx", strings.NewReader(values.Encode()))
	if err != nil {
		fmt.Errorf("error occured %v", err)
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("X-Requested-With", "XMLHttpRequest")
	request.Header.Set("Host", "marknad.byggvesta.se")
	request.Header.Set("Referer", "https://marknad.byggvesta.se/pgSearchResult.aspx")
	request.Header.Set("X-Momentum-API-KEY", "dIwHbOLgCS+FoZLYNYNToP9zK6VUoSgVC8BOT6cYljU=")
	request.Header.Set("Cookie","ga=GA1.2.1142674423.1542053180; _gid=GA1.2.1223576119.1542053180; ASP.NET_SessionId=b1csjxwfo03ydyfxz3yoh1qy; Language=se")
	request.Header.Set("User-Agent","Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36")
	request.Header.Set("Pragma","no-cache")
	request.Header.Set("Origin","https://marknad.byggvesta.se")
	request.Header.Set("Accept-Encoding","gzip, deflate, br")
	request.Header.Set("Accept-Language","en-US,en;q=0.9,ar;q=0.8,de;q=0.7")
	request.Header.Set("Accept","application/json,text/*")
	request.Header.Set("Cache-Control","no-cache")

	resp, err := client.Do(request)
	defer resp.Body.Close()
	if err != nil {
		fmt.Errorf("error occured %v", err)
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("error occured %v", err)
	}

	fmt.Printf("response %v", string(contents))
}
