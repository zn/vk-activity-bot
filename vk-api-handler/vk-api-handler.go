package vk_api_handler

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	urlHelper "net/url"
)

type VkApiHandler struct{
	GroupId string
	ApiToken string
	Version string
}

const baseUrl = "https://api.vk.com/method/"

// см. https://vk.com/dev/groups.getMembers "Пример запроса"
type GetMembersResult struct{
	Response struct{
		Count int `json:"count"`
		Items []int `json:"items"`
	} `json:"response"`
}

// Возвращает список участников сообщества
func (handler VkApiHandler) GetMembersList() GetMembersResult{
	if handler.ApiToken == "" ||
		handler.Version == "" ||
		handler.GroupId == ""{
		panic(errors.New("ApiToken and/or GroupId and/or Version is nil"))
	}

	params := map[string]string	{
		"group_id" : handler.GroupId,
		"access_token" : handler.ApiToken,
		"v" : handler.Version,
	}

	var result GetMembersResult
	callMethod("groups.getMembers", params, &result)
	return result
}

// Делает вызов API-метода и преобразует полученный json в структуру parseTo
func callMethod(method string, params map[string]string, parseTo interface{}){
	url := addParameters(baseUrl + method, params)
	req, err := http.Get(url)
	if err != nil{
		// TODO: handle error
		log.Fatal(err)
	}
	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil{
		// TODO: handle error
		log.Fatal(err)
	}
	jsonErr := json.Unmarshal(body, &parseTo)
	if jsonErr != nil{
		// TODO: handle error
		log.Fatal(err)
	}
}

// Добавляет параметры к url (например: group_id, offset и т.п.)
func addParameters(url string, params map[string]string) string{
	parsedUrl, _ := urlHelper.Parse(url)
	query := urlHelper.Values{}
	for key, value := range params{
		query.Add(key, value)
	}
	parsedUrl.RawQuery = query.Encode()
	return parsedUrl.String()
}