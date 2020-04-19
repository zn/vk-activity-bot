package main

import (
	db "github.com/zn/vk-activity-bot/db-worker"
	vk "github.com/zn/vk-activity-bot/vk-api-handler"
)

const GROUP_ID string = "group id here"
const API_VERSION string = "5.103"
const API_TOKEN string = "api token here"

func main(){
	vkHandler := vk.VkApiHandler{
		GroupId: GROUP_ID,
		ApiToken: API_TOKEN,
		Version:  API_VERSION,
	}
	result := vkHandler.GetMembersList()
	db.UpdateSubscribers(result.Response.Items, GROUP_ID)
}
