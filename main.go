package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	

)


func main() {

	err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file:", err)
    }

	InitDB()
    MigrateDB()

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
        log.Fatal("Telegram bot token is not set in .env file")
    }

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        dbHost, dbPort, dbUser, dbPassword, dbName)

	log.Println("Connecting to database:", psqlInfo)

	u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				userID := update.Message.From.ID
				nickname := update.Message.From.UserName
				userName := update.Message.From.FirstName + " " + update.Message.From.LastName
				
				user, isVerified := findUser(userID)
				
				if user != nil {
					if isVerified {
						keyboard := getMainMenuKeyboard()
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Привет, %s! \n Выбери пункт меню:", userName))
						msg.ReplyMarkup = keyboard
						bot.Send(msg)
					} else {
						bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Привет, %s! \n Ваш аккаунт не верифицирован. \n Нажмите /start для повторной проверки.", userName)))
					}
				} else {
					addUser(userID, nickname, userName)
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Привет, %s! Добавил вас в базу, ожидайте верификации от администратора!", userName)))
				}
			}
		}
	}
}

func addUser(userID int64, nickname, userName string) {
	user := User{
		ID: userID,
		Nickname: nickname,
		UserName: userName,
	}

	result := db.Create(&user)
	if result.Error != nil {
		log.Println("Error adding user:", result.Error)
		return
	}

	log.Printf("User %d added successfully!", userID)
}

func findUser(userID int64) (*User, bool) {
	var user User
	result := db.First(&user, "id = ?", userID)
	if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            return nil, false
        }
        log.Println("Error finding user:", result.Error)
        return nil, false
    }
	return &user, user.Verify
}

func getMainMenuKeyboard() tgbotapi.InlineKeyboardMarkup {
    keyboard := tgbotapi.NewInlineKeyboardMarkup(
        tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("Поиск книги", "search"),
            tgbotapi.NewInlineKeyboardButtonData("Просмотр библиотеки", "library"),
        ),
        tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("Добавить книгу", "add_book"),
        ),
    )

    return keyboard
}