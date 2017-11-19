package main

import (
	"log"

	"net/http"

	"fmt"
	"io"

	"os"

	"os/exec"

	"math/rand"

	"time"

	"gopkg.in/telegram-bot-api.v4"
)

const (
	token          = "477704380:AAFgpiVD2Zs0AkZv-VQ590K_dKCpjzKwYEE"
	welcomeMessage = "Welcome to Primitive Pic bot!\nThe original algorithm is provided by https://github.com/fogleman/primitive\n" +
		"I'm only writing this bot to " +
		"bring his art to everybody via a simple telegram bot.\n" +
		"Send me a photo and you'll see the magic!"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {

		switch {

		case update.Message == nil:
			continue

		case update.Message.Photo != nil:
			go handlePhoto(bot, update)

		case update.Message.Text == "/start":
			go sendMessage(welcomeMessage, bot, update)

		case update.Message.Photo == nil, update.Message.Text != "":
			go sendMessage("Send me a photo and you'll see the magic!", bot, update)
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
	}
}

func sendMessage(msgString string, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgString)

	bot.Send(msg)
}

func handlePhoto(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	picURL, err := bot.GetFileDirectURL((*update.Message.Photo)[1].FileID)
	if err != nil {
		log.Panic(err)
	}

	picName := fmt.Sprint(update.Message.From.UserName, update.Message.Chat.ID, update.Message.MessageID)
	downloadPhotoFromURL(picURL, picName)

	//TODO: say I'm only using a raspberry to take donations
	sendMessage("Modifying the photo.. It will take a while", bot, update)
	picModifiedName := modifyPhoto(picName)

	c := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, picModifiedName)
	bot.Send(c)

	println("Success")

}

func modifyPhoto(picName string) string {
	picModifiedName := fmt.Sprint(picName, "_mod.png")

	mode := getRandomMode()
	num := getRandomNum()

	cmd := exec.Command("./primitive", "-i", picName, "-o", picModifiedName, "-n", num, "-m", mode, "-r", "256")

	cmd.Run()

	cmd = exec.Command("rm", picName)
	go cmd.Run()

	return picModifiedName
}

func getRandomNum() string {
	rand.Seed(int64(time.Now().Second()))
	num := rand.Intn(400-170) + 170
	println("Num: ", num)
	return fmt.Sprintf("%d", num)
}

func getRandomMode() string {
	rand.Seed(int64(time.Now().Second()))
	mode := rand.Intn(9)
	println("Mode: ", mode)
	return fmt.Sprintf("%d", mode)
}

func downloadPhotoFromURL(url string, picName string) {
	response, e := http.Get(url)
	if e != nil {
		log.Panic(e)
	}

	defer response.Body.Close()

	file, err := os.Create(picName)
	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
}
