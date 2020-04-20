package vk_api_handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	urlHelper "net/url"
	"strconv"
	"time"
)

type VkApiHandler struct{
	GroupId string
	Version string
	ApiToken string
	RequestsPerSecond int // 3 or 20
}

const baseUrl = "https://api.vk.com/method/"

type GetMembersResult struct{
	Response struct{
		Done bool `json:"done"`
		Offset int `json:"offset"`
		Items []int `json:"items"`
	} `json:"response"`
}

// Возвращает список участников сообщества
func (handler VkApiHandler) GetMembersList(){
	if handler.ApiToken == "" ||
		handler.Version == "" ||
		handler.GroupId == ""{
		panic(errors.New("ApiToken and/or GroupId and/or Version is nil"))
	}

	var script = `
		var offset = parseInt(Args.offset);
		if(offset == null) {
			offset = 0;
		}
		var group_id = Args.group_id;
		var calls = 0;
		var users = [];
		
		while(calls < 25) {
			var received = [API.groups.getMembers({"group_id":group_id,"offset":offset})]@.items[0];
			if(received.length == 0){
				return {"done":true, "offset":offset, "items":users};
			}
			offset = offset + received.length;
			users = users + received;
		
			calls = calls + 1;
		}
		
		return {"done":false,"offset":offset, "items":users};`

	params := map[string]string	{
		"code":script,
		"group_id" : handler.GroupId,
		"offset": "0",
		"access_token" : handler.ApiToken,
		"v" : handler.Version,
	}

	fmt.Println("Started!", time.Now().Format("15:04:05"))
	var sum int
	var counter int
	for {
		var result GetMembersResult
		callMethod("execute", params, &result)
		params["offset"] = strconv.Itoa(result.Response.Offset)
		sum += len(result.Response.Items)
		// <- может быть тут в chan записывать
		if result.Response.Done{
			break
		}
		counter++
		if counter == handler.RequestsPerSecond{
			time.Sleep(time.Second)
			counter = 0
		}
	}
	fmt.Println("Finished!", time.Now().Format("15:04:05"))
	fmt.Println("Всего получено:", sum)
}

// Делает вызов API-метода и преобразует полученный json в структуру parseTo
func callMethod(method string, params map[string]string, parseTo interface{}){
	url := addParameters(baseUrl + method, params)
	req, err := http.Get(url)
	if err != nil{
		log.Fatal(err)
	}
	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil{
		log.Fatal(err)
	}

	jsonErr := json.Unmarshal(body, &parseTo)
	if jsonErr != nil{
		fmt.Println(string(body))
		log.Fatal(jsonErr)
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