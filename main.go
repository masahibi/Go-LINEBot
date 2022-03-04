package main

import (
	"fmt"
	"github.com/line/line-bot-sdk-go/linebot"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {

	ChanellSecret := "667d346382f992671b4da40684f971bf"
	ChanellToken := "9Cdz5cvK4etUt5abkUeJ8OKiR0lOr3Ys/eYJzZp8bZL3srIDMkVhwe5GqAlYMBKkU41cSAMeNvhKnEOc711qvnoTpYRye4kk0asipZvwzrgoDTuT8LWIRZFEhtaUmJN85K+UbBinsI9VaaAxeAL99gdB04t89/1O/w1cDnyilFU="
	bot, err := linebot.New(ChanellSecret, ChanellToken) // LINEBotのインスタンスを生成
	if err != nil {                                      // エラーが発生した場合
		log.Fatal(err) // エラー内容を出力して終了
	}

	// Setup HTTP Server for receiving requests from LINE platform
	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		events, err := bot.ParseRequest(req) // ParseRequestでリクエストをパースして、eventsに格納
		if err != nil {                      // エラーが発生した場合
			if err == linebot.ErrInvalidSignature { // シグネチャが誤っている場合
				w.WriteHeader(400) // 400 Bad Requestを返す
			} else { // それ以外のエラーの場合
				w.WriteHeader(500) // 500 Internal Server Errorを返す
			}
			return // ParseRequestのエラーを返して終了
		}
		taskBook := NewTaskBook("taskbook.txt") // TaskBookのインスタンスを生成
		view := "タスク"
		add := "追加"
		addLimit := 3 // 追加、内容、日付
		del := "完了"
		delLimit := 2 // 完了、内容

		for _, event := range events { // ParseRequestでパースしたイベントをループ
			if event.Type == linebot.EventTypeMessage { // TypeがMessageの場合
				switch message := event.Message.(type) { // Messageの型を判定
				case *linebot.TextMessage: // Messageがテキストの場合
					testMessage := strings.Split(message.Text, ",") // ',' 区切りで分割してスライスにする
					if testMessage[0] == view {                     // 表示の場合
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(showItems(taskBook.tasks))).Do(); err != nil { // ReplyMessageで返信
							log.Print(err) // エラー内容を出力
						}
					} else if len(testMessage) == addLimit && testMessage[0] == add { // 追加の場合
						taskBook.AddTask(inputTask(testMessage[1], testMessage[2]))                                          // AddTaskでタスクを追加
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("追加しました！")).Do(); err != nil { // ReplyMessageで返信
							log.Print(err) // エラー内容を出力
						}
					} else if len(testMessage) == delLimit && testMessage[0] == del { // 完了の場合
						taskBook.DelTask(testMessage[1])                                                                                // DelTaskでタスクを削除
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("削除しました！\nお疲れさまでした！")).Do(); err != nil { // ReplyMessageで返信
							log.Print(err) // エラー内容を出力
						}
					} else {                                                                                                    // それ以外の場合
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil { // ReplyMessageで返信
							log.Print(err) // エラー内容を出力
						}
					}

					//if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil { // ReplyMessageで返信
					//	log.Print(err) // エラー内容を出力
					//}
				case *linebot.StickerMessage: // Messageがスタンプの場合
					replyMessage := fmt.Sprintf( // テキストを作成
						"sticker id is %s, stickerResourceType is %s", message.StickerID, message.StickerResourceType)      // スタンプIDとスタンプリソースタイプを出力
					replyMessage = fmt.Sprint("It's a nice sticker !!")                                                     // テキストを作成
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil { // ReplyMessageで返信
						log.Print(err) // エラー内容を出力
					}
				}
			}
		}
	})
	// This is just sample code.
	// For actual use, you must support HTTPS by using `ListenAndServeTLS`, a reverse proxy or something else.
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil { // ListenAndServeでHTTPサーバを起動
		log.Fatal(err) // ListenAndServeのエラーを出力して終了
	}
}

// Taskを入力し返す
func inputTask(category string, date string) *Task {
	var task Task

	task.Category = category
	task.Date = date

	return &task
}

// Taskの一覧を出力する
func showItems(items []*Task) string {
	//fmt.Println("===========")
	// itemsの要素を1つずつ取り出してitemに入れて繰り返す
	var text string
	var date time.Time
	for i, task := range items {
		date, _ = time.Parse("08/28", task.Date)
		text += fmt.Sprintf("[%d] %s : %d月%d日\n", i+1, task.Category, int(date.Month()), date.Day())
	}
	//fmt.Println("===========")
	return text
}
