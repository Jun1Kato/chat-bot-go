package main

import (
	"net/http"
	"os"

	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
	"math/rand"
	"time"
)

const (
	defaultMessage = "ちょっと聞こえないです"
)

func main() {
	port := os.Args[1]
	channelSecret := os.Args[2]
	accessToken := os.Args[3]

	if port == "" || channelSecret == "" || accessToken == "" {
		panic("$PORT and $CHANNEL_SECRET and $ACCESS_TOKEN must be set")
	}
	fmt.Println("chat-bot-server start")

	router := gin.New()
	router.Use(gin.Logger())

	// webhook
	router.POST("/webhook", func(c *gin.Context) {
		client := &http.Client{Timeout: time.Duration(15 * time.Second)}
		bot, err := linebot.New(channelSecret, accessToken, linebot.WithHTTPClient(client))
		if err != nil {
			fmt.Println(err)
			return
		}
		received, err := bot.ParseRequest(c.Request)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, event := range received {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					source := event.Source
					if source.Type == linebot.EventSourceTypeUser {
						repMessage := defaultMessage
						if resMessage := getResMessage(message.Text); resMessage != "" {
							repMessage = resMessage
						}
						postMessage := linebot.NewTextMessage(repMessage)
						if _, err = bot.ReplyMessage(event.ReplyToken, postMessage).Do(); err != nil {
							fmt.Print(err)
						}
					}
				default:
					fmt.Printf("This Message is not linebot.TextMessage. event.MessageType : %v \n", event.Message)
				}
			} else {
				fmt.Printf("event.Type is not linebot.EventTypeMessage. event.Type : %v \n", event.Type)
			}
		}
	})

	// server setting
	s := &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}

// getResMessage : メッセージ処理
func getResMessage(reqMessage string) (message string) {
	fmt.Println("getResMessage start")
	defer fmt.Println("getResMessage end")
	resMessages := [3]string{"わかるわかる", "それで？それで？", "からの〜？"}

	rand.Seed(time.Now().UnixNano())
	if rand.Intn(2) == 0 {
		if math := rand.Intn(4); math != 3 {
			message = resMessages[math]
		} else {
			message = reqMessage + "じゃねーよw"
		}
	}
	return
}
