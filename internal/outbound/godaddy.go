package outbound

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kecci/goscription/models"
	"github.com/spf13/viper"
)

type godaddyOutbound struct {
}

// NewGodaddyOutbound will create new an articleUsecase object representation of usecase.ArticleUsecase interface
func NewGodaddyOutbound() GodaddyOutbound {
	return &godaddyOutbound{}
}

// GodaddyOutbound represent the usecase of the domain
type GodaddyOutbound interface {
	GetDomainAvailable(domain string) (res models.DomainAvailableResponse, err error)
}

func (g *godaddyOutbound) GetDomainAvailable(domain string) (res models.DomainAvailableResponse, err error) {
	var domainAvailableResponse = models.DomainAvailableResponse{}

	method := "GET"
	endpoint := "/v1/domains/available"
	request, _ := goddadyBaseRequest(method, endpoint)

	q := request.URL.Query()
	q.Add("domain", domain)
	request.URL.RawQuery = q.Encode()

	client := &http.Client{}
	response, err := client.Do(request)

	var data []byte

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ = ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
	}

	err = json.Unmarshal(data, &domainAvailableResponse)

	if err != nil {
		return domainAvailableResponse, err
	}

	return domainAvailableResponse, nil
}

func goddadyBaseRequest(method string, endpoint string) (*http.Request, error) {

	host := viper.GetString(`godaddy.host`)
	authorization := viper.GetString(`godaddy.authorization`)
	url := host + endpoint

	request, err := http.NewRequest(method, url, nil)

	request.Header.Set("Authorization", authorization)
	request.Header.Set("Content-Type", "application/json")

	return request, err
}
