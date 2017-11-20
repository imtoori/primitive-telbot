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
	"bytes"
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
	picModifiedName, modeString, num := modifyPhoto(picName)

	c := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, picModifiedName)
	bot.Send(c)
	sendMessage(fmt.Sprint(num, " ", modeString), bot, update)

	println("Success")

}

func modifyPhoto(picName string) (picModifiedName string, modeString string, num string) {
	picModifiedName := fmt.Sprint(picName, "_mod.png")

	mode := getRandomMode()
	num := getRandomNum()

	cmd := exec.Command("./primitive", "-i", picName, "-o", picModifiedName, "-n", num, "-m", mode, "-r", "256")

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Println(err)
		log.Println(stderr.String())
	}

	switch mode {
	case "0":
		modeString = "combo"
	case "1":
		modeString = "triangles"
	case "2":
		modeString = "rectangles"
	case "3":
		modeString = "ellipses"
	case "4":
		modeString = "circles"
	case "5":
		modeString = "rotated rectangles"
	case "6":
		modeString = "beziers"
	case "7":
		modeString = "rotated ellipses"
	case "8":
		modeString = "polygons"
	}

	cmd = exec.Command("rm", picName)
	go cmd.Run()

	return
}

func getRandomNum() string {
	rand.Seed(int64(time.Now().Nanosecond()))
	num := rand.Intn(1300-170) + 170
	println("Num: ", num)
	return fmt.Sprintf("%d", num)
}

func getRandomMode() string {
	rand.Seed(int64(time.Now().Nanosecond()))
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
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}
}
