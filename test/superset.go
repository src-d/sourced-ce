// +build integration

package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/publicsuffix"
)

type supersetClient struct {
	*http.Client
	csrf string
}

func newSupersetClient() (*supersetClient, error) {
	// Superset client needs a session cookie
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Jar: jar,
	}

	// To POST in /login we need to first GET /login, read the hidden CSRF value
	// from the HTML, and send it back in the POST
	res, err := client.Get("http://127.0.0.1:8088/login/")
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	csrf, ok := doc.Find("input#csrf_token").First().Attr("value")
	if !ok {
		return nil, fmt.Errorf("CSRF token was not found in the /login page")
	}

	res, err = client.PostForm("http://127.0.0.1:8088/login/", url.Values{
		"csrf_token": {csrf},
		"username":   {"admin"},
		"password":   {"admin"},
	})
	if err != nil {
		return nil, err
	}

	// After a successful POST the client is authenticated and can be used to call
	// the API
	return &supersetClient{client, csrf}, nil
}

func (c *supersetClient) dashboards() ([]string, error) {
	res, err := c.Get("http://127.0.0.1:8088/dashboard/api/read")
	if err != nil {
		return nil, err
	}

	var decoded struct {
		Result []struct {
			Link string `json:"dashboard_link"`
		}
	}

	err = json.NewDecoder(res.Body).Decode(&decoded)
	if err != nil {
		return nil, err
	}

	links := []string{}
	for _, result := range decoded.Result {
		links = append(links, result.Link)
	}

	return links, nil
}

func (c *supersetClient) sql(query, dbId, schema string) ([]map[string]interface{}, error) {
	res, err := c.PostForm("http://127.0.0.1:8088/superset/sql_json/", url.Values{
		//"client_id":      {""},
		"database_id": {dbId},
		"json":        {"true"},
		"runAsync":    {"false"},
		"schema":      {schema},
		"sql":         {query},
		//"sql_editor_id":  {"jzl0KCm5Z"},
		//"tab":            {"Untitled Query 2"},
		//"tmp_table_name": {""},
		//"select_as_cta":  {"false"},
		//"templateParams": {"{}"},
		//"queryLimit": {"1000"},

		"csrf_token": {c.csrf},
	})
	if err != nil {
		return nil, err
	}

	var decoded struct {
		Status string
		Error  string
		Data   []map[string]interface{}
	}

	err = json.NewDecoder(res.Body).Decode(&decoded)
	if err != nil {
		return nil, err
	}

	if decoded.Status != "success" {
		return nil, fmt.Errorf("/sql_json endpoint returned an error: " + decoded.Error)
	}

	return decoded.Data, nil
}

func (c *supersetClient) gitbase(query string) ([]map[string]interface{}, error) {
	return c.sql(query, "1", "gitbase")
}

func (c *supersetClient) metadata(query string) ([]map[string]interface{}, error) {
	return c.sql(query, "2", "public")
}

func (c *supersetClient) bblfsh(filename, content string) (string, error) {
	var jsonStr = []byte(fmt.Sprintf(
		`{"mode":"semantic", "filename":%q, "content":%q, "query":""}`,
		filename, content))

	req, err := http.NewRequest("POST", "http://127.0.0.1:8088/bblfsh/api/parse", bytes.NewBuffer(jsonStr))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.Do(req)
	if err != nil {
		return "", err
	}

	var decoded struct {
		Status   int
		Language string
		// Uast     interface{}
		Errors []struct{ Message string }
	}

	err = json.NewDecoder(res.Body).Decode(&decoded)
	if err != nil {
		body := ""
		bytes, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			body = string(bytes)
		}

		return "", fmt.Errorf("could not decode the response body: %s, err: %s", body, err)
	}

	if decoded.Status != 0 {
		return "", fmt.Errorf("/bblfsh/api/parse endpoint returned an error: %v", decoded.Errors)
	}

	return decoded.Language, nil
}
