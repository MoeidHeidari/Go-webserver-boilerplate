package currencies

import (
	"encoding/json"
	"io/ioutil"
	"main/lib"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

type Request struct {
	logger lib.Logger
}

func NewRequest(logger lib.Logger) Request {
	return Request{
		logger: logger,
	}
}

// @Summary Gets currencies
// @Tags get tests
// @Description Get one test by id
// @Security ApiKeyAuth
// @Router /api/currency [get]
func (u Request) MakeRequest(c *gin.Context) {
	client := http.DefaultClient
	url_test := "https://www.rbc.ru/crypto/currency/btcusd"
	req, err := http.NewRequest("GET", url_test, nil)
	if err != nil {
		u.logger.Error(err)
	}

	resp, err := client.Do(req)

	if err != nil {
		u.logger.Error(err)
	}

	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		u.logger.Error(err)
	}

	doc.Find("div").Each(func(i int, s *goquery.Selection) {
		if s.AttrOr("class", "") == "chart__subtitle js-chart-value" {
			currency_rate := strings.TrimSpace(strings.Split(s.Text(), "\n")[1])
			c.JSON(200, gin.H{
				"Bitcoin price now (USD)": currency_rate,
			})

		}
	})

}

// @Summary Gets post responce
// @Tags get tests
// @Accept json
// @Produce json
// @Param input body string true "Post form"
// @Description Post request
// @Security ApiKeyAuth
// @Router /api/currency [post]
func (u Request) MakePostRequest(c *gin.Context) {
	url_test := "https://httpbin.org/post"
	i, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		u.logger.Error(err)
	}
	url_form := url.Values{
		"name": {"GodFather"},
		"data": {string(i)},
	}

	url_form.Add("value", "1234")
	resp, err := http.PostForm(url_test, url_form)

	if err != nil {
		u.logger.Error(err)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	defer resp.Body.Close()
	c.JSON(200, gin.H{
		"responce": result,
	})

}
