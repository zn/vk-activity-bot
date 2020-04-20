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
		ApiToken: API_TOKEN,
		Version:  API_VERSION,
	}
	vkHandler.GetMembersList()
	//db.UpdateSubscribers(result.Response.Items, GROUP_ID)

}
