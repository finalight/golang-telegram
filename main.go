package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func setWebhook(url string, botURL string) {
	response, err := http.Get(url + "setWebhook?url=" + botURL)

	if err != nil {
		fmt.Println("webhook error", err)
	}

	fmt.Println("webhook response", response)
}

func sendMessage(baseURL string, chatID string, message string) {
	jsonString := fmt.Sprintf(`{"chat_id" : 184935795, "text": "%s"}`, message)
	jsonStr := []byte(jsonString)

	req, err := http.NewRequest("POST", baseURL+"sendMessage", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error", err)
	}

	defer resp.Body.Close()
	// _, _ := ioutil.ReadAll(resp.Body)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	botURL := os.Getenv("BOT_URL")
	telegramKey := os.Getenv("TELEGRAM_BOT_KEY")
	baseURL := "https://api.telegram.org/bot" + telegramKey + "/"

	setWebhook(baseURL, botURL)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/", func(c *gin.Context) {
		fmt.Println(c.Request.Body)
		c.JSON(200, gin.H{})
	})

	r.POST("/", func(c *gin.Context) {
		data, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			fmt.Println(err)
		}

		chatID, _, _, err := jsonparser.Get(data, "message", "from", "id")
		message, _, _, err := jsonparser.Get(data, "message", "text")

		sendMessage(baseURL, string(chatID), string(message))

		c.JSON(200, gin.H{})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
