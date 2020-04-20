package main

import (
	vk "github.com/zn/vk-activity-bot/vk-api-handler"
)

const GROUP_ID string = "englishyo"
const API_VERSION string = "5.103"
const API_TOKEN string = ""

func main(){
	vkHandler := vk.VkApiHandler{
		GroupId: GROUP_ID,
		Version:  API_VERSION,
		ApiToken: API_TOKEN,
		RequestsPerSecond: 3,
	}
	vkHandler.GetMembersList()
	//db.UpdateSubscribers(result.Response.Items, GROUP_ID)

}
