package keyboards

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/reply"

	"github.com/insigmo/asisa_clinic_finder/internal/db"
	"github.com/insigmo/asisa_clinic_finder/internal/local_models"
)

const (
	ActionFindClinic     = "find_clinic"
	ActionChangeCity     = "change_city"
	ActionChangeLanguage = "change_language"

	RuCommandFindClinic     = "Найти поликлинику"
	RuCommandChangeCity     = "Поменять город"
	RuCommandChangeLanguage = "Поменять язык"

	EsCommandFindClinic     = "Buscar clínica"
	EsCommandChangeCity     = "Cambiar ciudad"
	EsCommandChangeLanguage = "Cambiar idioma"

	EnCommandFindClinic     = "Find clinic"
	EnCommandChangeCity     = "Change city"
	EnCommandChangeLanguage = "Change language"
)

type MenuTexts struct {
	Prompt                 string
	FindClinicButton       string
	ChangeCityButton       string
	ChangeLanguageButton   string
	ChooseLanguagePrompt   string
	LanguageSavedMessage   string
	UnknownMenuAction      string
	UnknownLanguageMessage string
	FindClinicHelpMessage  string
}

func menuTexts(languageCode string) MenuTexts {
	switch normalizeLanguage(languageCode) {
	case "ru":
		return MenuTexts{
			Prompt:                 "Выберите действие:",
			FindClinicButton:       "Найти поликлинику",
			ChangeCityButton:       "Поменять город",
			ChangeLanguageButton:   "Поменять язык",
			ChooseLanguagePrompt:   "Выберите язык:",
			LanguageSavedMessage:   "Язык сохранён.",
			UnknownMenuAction:      "Неизвестное действие меню.",
			UnknownLanguageMessage: "Неизвестный язык.",
			FindClinicHelpMessage:  "Используй команду: /find_clinic <медицинское направление>",
		}
	case "es":
		return MenuTexts{
			Prompt:                 "Elige una acción:",
			FindClinicButton:       "Buscar clínica",
			ChangeCityButton:       "Cambiar ciudad",
			ChangeLanguageButton:   "Cambiar idioma",
			ChooseLanguagePrompt:   "Elige un idioma:",
			LanguageSavedMessage:   "Idioma guardado.",
			UnknownMenuAction:      "Acción de menú desconocida.",
			UnknownLanguageMessage: "Idioma desconocido.",
			FindClinicHelpMessage:  "Usa el comando: /find_clinic <especialidad médica>",
		}
	default:
		return MenuTexts{
			Prompt:                 "Choose an action:",
			FindClinicButton:       "Find clinic",
			ChangeCityButton:       "Change city",
			ChangeLanguageButton:   "Change language",
			ChooseLanguagePrompt:   "Choose a language:",
			LanguageSavedMessage:   "Language saved.",
			UnknownMenuAction:      "Unknown menu action.",
			UnknownLanguageMessage: "Unknown language.",
			FindClinicHelpMessage:  "Use command: /find_clinic <medical direction>",
		}
	}
}

func normalizeLanguage(languageCode string) string {
	code := strings.TrimSpace(strings.ToLower(languageCode))
	switch {
	case strings.HasPrefix(code, "ru"):
		return "ru"
	case strings.HasPrefix(code, "es"):
		return "es"
	default:
		return "en"
	}
}

func resolveUserLanguage(ctx context.Context, update *models.Update) string {
	if update != nil && update.Message != nil && update.Message.From != nil {
		if languageCode := normalizeLanguage(update.Message.From.LanguageCode); languageCode != "" {
			return languageCode
		}
	}

	dbManager, ok := ctx.Value(local_models.DBManagerKey).(*db.Manager)
	if !ok || update == nil || update.Message == nil {
		return "en"
	}

	user, err := dbManager.GetUser(ctx, update.Message.Chat.ID)
	if err != nil || user == nil {
		return "en"
	}

	return normalizeLanguage(user.LanguageCode)
}

func BuildMainMenu(tgBot *bot.Bot, languageCode string) *reply.ReplyKeyboard {
	texts := menuTexts(languageCode)

	return reply.New(
		reply.WithPrefix("main_menu"),
		reply.IsSelective(),
		reply.ResizableKeyboard(),
	).
		Button(texts.FindClinicButton, tgBot, bot.MatchTypeExact, onMainMenuSelect).
		Row().
		Button(texts.ChangeCityButton, tgBot, bot.MatchTypeExact, onMainMenuSelect).
		Button(texts.ChangeLanguageButton, tgBot, bot.MatchTypeExact, onMainMenuSelect)
}

func BuildLanguageMenu(tgBot *bot.Bot) *reply.ReplyKeyboard {
	return reply.New(
		reply.WithPrefix("language_menu"),
		reply.IsSelective(),
		reply.ResizableKeyboard(),
	).
		Button("Русский", tgBot, bot.MatchTypeExact, onLanguageSelect).
		Button("Español", tgBot, bot.MatchTypeExact, onLanguageSelect).
		Button("English", tgBot, bot.MatchTypeExact, onLanguageSelect)
}

func SendMainMenu(ctx context.Context, tgBot *bot.Bot, update *models.Update) error {
	languageCode := resolveUserLanguage(ctx, update)
	texts := menuTexts(languageCode)

	_, err := tgBot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        texts.Prompt,
		ReplyMarkup: BuildMainMenu(tgBot, languageCode),
	})
	if err != nil {
		return fmt.Errorf("send main menu: %w", err)
	}

	return nil
}

func SendLanguageMenu(ctx context.Context, tgBot *bot.Bot, update *models.Update) error {
	languageCode := resolveUserLanguage(ctx, update)
	texts := menuTexts(languageCode)

	_, err := tgBot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        texts.ChooseLanguagePrompt,
		ReplyMarkup: BuildLanguageMenu(tgBot),
	})
	if err != nil {
		return fmt.Errorf("send language menu: %w", err)
	}

	return nil
}

func ResolveMainMenuAction(languageCode, text string) (string, bool) {
	texts := menuTexts(languageCode)
	value := strings.TrimSpace(text)

	switch value {
	case texts.FindClinicButton:
		return ActionFindClinic, true
	case texts.ChangeCityButton:
		return ActionChangeCity, true
	case texts.ChangeLanguageButton:
		return ActionChangeLanguage, true
	default:
		return "", false
	}
}

func ResolveLanguage(text string) (string, bool) {
	switch strings.TrimSpace(text) {
	case "Русский":
		return "ru", true
	case "Español":
		return "es", true
	case "English":
		return "en", true
	default:
		return "", false
	}
}

func FindClinicHelpMessage(languageCode string) string {
	return menuTexts(languageCode).FindClinicHelpMessage
}

func HandleLanguageSelection(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
	onLanguageSelect(ctx, tgBot, update)
}

func onMainMenuSelect(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
	languageCode := resolveUserLanguage(ctx, update)
	action, ok := ResolveMainMenuAction(languageCode, update.Message.Text)
	if !ok {
		texts := menuTexts(languageCode)
		_, _ = tgBot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   texts.UnknownMenuAction,
		})
		return
	}

	var response string
	switch action {
	case ActionFindClinic:
		response = menuTexts(languageCode).FindClinicHelpMessage
	case ActionChangeCity:
		response = menuTexts(languageCode).ChangeCityButton
	case ActionChangeLanguage:
		response = menuTexts(languageCode).ChangeLanguageButton
	default:
		response = menuTexts(languageCode).UnknownMenuAction
	}

	_, _ = tgBot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   response,
	})
}

func onLanguageSelect(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
	params := local_models.NewBaseParams(ctx, tgBot, update)

	dbManager, ok := ctx.Value(local_models.DBManagerKey).(*db.Manager)
	if !ok {
		params.Log.Error("dbManager is not set to context")
		return
	}

	newLanguage, ok := ResolveLanguage(update.Message.Text)
	if !ok {
		texts := menuTexts(resolveUserLanguage(ctx, update))
		_, _ = tgBot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   texts.UnknownLanguageMessage,
		})
		return
	}

	userInfo := update.Message.From
	user, err := dbManager.GetUser(ctx, update.Message.Chat.ID)
	if err != nil || user == nil {
		user = &db.User{
			ID:       userInfo.ID,
			Username: userInfo.Username,
			Name:     userInfo.FirstName,
			Lastname: userInfo.LastName,
			IsBot:    userInfo.IsBot,
		}
	}

	user.LanguageCode = newLanguage

	if err = dbManager.InsertOrUpdateUser(ctx, user); err != nil {
		params.Log.Error(err.Error())
		return
	}

	texts := menuTexts(newLanguage)

	_, _ = tgBot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        texts.LanguageSavedMessage,
		ReplyMarkup: BuildMainMenu(tgBot, newLanguage),
	})
}
