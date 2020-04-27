package vk

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	urlHelper "net/url"
	"strconv"
	"time"
)

type ApiHandler struct {
	GroupId           string
	Version           string
	ApiToken          string
	RequestsPerSecond int // 3 or 20
}

const baseUrl = "https://api.vk.com/method/"

type GetMembersResult struct {
	Response struct {
		Done   bool  `json:"done"`
		Offset int   `json:"offset"`
		Items  []int `json:"items"`
	} `json:"response"`
}

// Возвращает список участников сообщества
func (handler ApiHandler) GetMembersList() ([]int, error) {
	if handler.ApiToken == "" ||
		handler.Version == "" ||
		handler.GroupId == "" {
		return nil, errors.New("ApiToken and/or GroupId and/or Version is nil")
	}

	params := map[string]string{
		"code":         GetSubscribersScript,
		"group_id":     handler.GroupId,
		"offset":       "0",
		"access_token": handler.ApiToken,
		"v":            handler.Version,
	}
	var items []int
	var counter int
	fmt.Println("Started", time.Now().Format("15:04:05"))
	for {
		var result GetMembersResult
		err := callMethod("execute", params, &result)
		if err != nil {
			return nil, err
		}
		params["offset"] = strconv.Itoa(result.Response.Offset)
		items = append(items, result.Response.Items...)
		if result.Response.Done {
			break
		}
		counter++
		if counter == handler.RequestsPerSecond {
			fmt.Println("Obtained:", params["offset"])
			time.Sleep(time.Second)
			counter = 0
		}
	}
	fmt.Println("Finished", time.Now().Format("15:04:05"))
	return items, nil
}

// Делает вызов API-метода и преобразует полученный json в структуру parseTo
func callMethod(method string, params map[string]string, parseTo interface{}) error {
	url := addParameters(baseUrl+method, params)
	req, err := http.Get(url)
	if err != nil {
		return err
	}
	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &parseTo)
	if err != nil {
		// скорее всего rate limit reached или unauthorized access
		return err
	}
	return nil
}

// Добавляет параметры к url (например: group_id, offset и т.п.)
func addParameters(url string, params map[string]string) string {
	parsedUrl, _ := urlHelper.Parse(url)
	query := urlHelper.Values{}
	for key, value := range params {
		query.Add(key, value)
	}
	parsedUrl.RawQuery = query.Encode()
	return parsedUrl.String()
}
