package controllers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/luckydevil2007/audionotes/usecases"
)

type TelegramBot struct {
	bot      *tgbotapi.BotAPI // Экземпляр бота API Telegram
	chatID   int64            // ID чата, куда будут отправляться сообщения
	vicinity int64
	//checkAuth *usecases.AuthUseCase
	note *usecases.NoteUseCase
	//repo      *repositories.Repository
	//producer  *producers.EventProducer
}

func NewTelegramBot(token string, note *usecases.NoteUseCase) (*TelegramBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	// Возвращаем инициализированный адаптер с ботом и ID чата
	return &TelegramBot{bot: bot, note: note}, nil
}

func (t *TelegramBot) Test(ctx context.Context) error {
	note, err := t.note.OpenNearest(ctx, 60.0, 30.0, 1000.0)
	if err != nil {
		return err
	}
	if len(note.Data) == 0 {
		return fmt.Errorf("data==0")
	}
	return nil
}

func (t *TelegramBot) Run(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	updates := t.bot.GetUpdatesChan(u)
	errChan := make(chan error)
	//go func() {
	for update := range updates {
		if update.Message == nil {
			continue
		}

		// Handle /start command
		if update.Message.IsCommand() && update.Message.Command() == "start" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome! Use /newtour to create a tour.")
			t.bot.Send(msg)
		}

		// Handle location sharing
		if update.Message.Location != nil {
			lat := float64(update.Message.Location.Latitude)
			lon := float64(update.Message.Location.Longitude)
			// Check if near a tour point (pseudo-c
			note, err := t.note.OpenNearest(ctx, lat, lon, 1.0)

			if err == nil {
				file := tgbotapi.FileBytes{
					Name:  note.Title,
					Bytes: note.Data,
				}

				audioConfig := tgbotapi.NewVoice(update.Message.Chat.ID, file)
				//audio := tgbotapi.NewAudioShare(update.Message.Chat.ID, audioURL)
				t.bot.Send(audioConfig)
			}
		}

		if update.Message.Voice != nil {
			fileID := update.Message.Voice.FileID
			fileConfig := tgbotapi.FileConfig{FileID: fileID}
			file, err := t.bot.GetFile(fileConfig)
			if err != nil {
				continue
			}
			nameTmp := t.bot.Token + "/" + file.FilePath
			url := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", t.bot.Token, file.FilePath)
			resp, err := http.Get(url)
			var data []byte
			data, err = io.ReadAll(resp.Body)
			t.note.Upload(ctx, strings.Replace(nameTmp, "/", "", -1), data, 1)
		}

		if update.Message.Audio != nil {
			fileID := update.Message.Audio.FileID
			fileConfig := tgbotapi.FileConfig{FileID: fileID}
			file, err := t.bot.GetFile(fileConfig)
			if err != nil {
				continue
			}
			nameTmp := t.bot.Token + "/" + file.FilePath
			url := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", t.bot.Token, file.FilePath)
			resp, err := http.Get(url)

			var data []byte
			data, err = io.ReadAll(resp.Body)

			t.note.Upload(ctx, strings.Replace(nameTmp, "/", "", -1), data, 1)
		}
	}
	//}()
	err := <-errChan
	return err
}

/*func (t *TelegramBot) DownloadFile(url string, filePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}*/

/*
func (t *TelegramBot) UploadNote(ctx tgbotapi.BotAPI.context.Context) error {
	voice := ctx.Message().Voice
	filePath := fmt.Sprintf("downloads/voice/%d_%s.ogg", c.Chat().ID, voice.FileID)

	if err := b.Download(&voice.File, filePath); err != nil {
		return c.Send("❌ Failed to download voice message")
	}

	return c.Send(fmt.Sprintf(
		"🎤 Voice message saved!\nDuration: %d sec", voice.Duration))
}

// SendMessage отправляет текстовое сообщение в Telegram-чат.
// Принимает контекст (ctx) и текст сообщения (message).
// Возвращает ошибку, если сообщение не удалось отправить.
func (t *TelegramAdapter) SendMessage(ctx context.Context, message string) error {
	// Создаём новое текстовое сообщение для отправки в указанный чат
	msg := tgbotapi.NewMessage(t.chatID, message)

	// Отправляем сообщение через API Telegram
	_, err := t.bot.Send(msg)

	// Возвращаем ошибку, если отправка не удалась
	return err
}
*/
