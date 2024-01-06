package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
	zLog "github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Chat anything with ChatGPT
// Command: chat
func ExecuteChat(discord *discordgo.Session, message *discordgo.MessageCreate) {
	chatGptReq, err := prepareChatGptRequest(message) // Prepare Request Body for ChatGPT
	if err != nil {
		log.Fatal(err)
		zLog.Error().Msg(err.Error())
		discord.ChannelMessageSend(message.ChannelID, "Encountered Error!")
		return
	}

	chatGptRes, err := invokeAskChatGpt(chatGptReq) // Calling ChatGPT API
	if err != nil {
		log.Fatal(err)
		zLog.Error().Msg(err.Error())
		discord.ChannelMessageSend(message.ChannelID, "Encountered Error!")
		return
	}

	// Check whether ChatGPT return error
	if chatGptRes["error"] != nil {
		errMsg := chatGptRes["error"].(map[string]interface{})["message"]
		discord.ChannelMessageSend(message.ChannelID, errMsg.(string))
		return
	}

	// Extract the content from the JSON response
	content := chatGptRes["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
	discord.ChannelMessageSend(message.ChannelID, content)
}

// Prepare Request Body for ChatGPT
func prepareChatGptRequest(message *discordgo.MessageCreate) ([]byte, error) {
	actualMessage := strings.Split(message.Content, " ")[1]
	var chatGptReq map[string]interface{} = map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"message": []interface{}{
			map[string]interface{}{
				"role":    "system",
				"content": actualMessage,
			},
		}, "max_tokens": 50,
	}

	data, err := json.Marshal(chatGptReq)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// POST Request to communicate with ChatGPT
func invokeAskChatGpt(chatGptReq []byte) (map[string]interface{}, error) {
	req, err := http.NewRequest(http.MethodPost, viper.GetString("openai.chatgpt.endpoint"), bytes.NewBuffer(chatGptReq))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+viper.GetString("openai.chatgpt.api-key"))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	/******************************************** ChatGPT Sample Response ********************************************/
	// {
	// 	"choices": [
	// 		{
	// 			"finish_reason": "stop",
	// 			"index": 0,
	// 			"message": {
	// 				"content": "The 2020 World Series was played in Texas at Globe Life Field in Arlington.",
	// 				"role": "assistant"
	// 			},
	// 			"logprobs": null
	// 		}
	// 	],
	// 	"created": 1677664795,
	// 	"id": "chatcmpl-7QyqpwdfhqwajicIEznoc6Q47XAyW",
	// 	"model": "gpt-3.5-turbo-0613",
	// 	"object": "chat.completion",
	// 	"usage": {
	// 		"completion_tokens": 17,
	// 		"prompt_tokens": 57,
	// 		"total_tokens": 74
	// 	}
	// }
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
