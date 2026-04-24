package keyboards

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/reply"

	"github.com/insigmo/asisa_clinic_finder/internal/db"
	"github.com/insigmo/asisa_clinic_finder/internal/helpers"
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

func menuTexts(languageCode string) (*MenuTexts, error) {
	if !slices.Contains([]string{"ru", "es", "en"}, languageCode) {
		return &MenuTexts{}, fmt.Errorf("unknown language: %s", languageCode)
	}
	switch normalizeLanguage(languageCode) {
	case "ru":
		return &MenuTexts{
			Prompt:                 "Выберите действие:",
			FindClinicButton:       "Найти поликлинику",
			ChangeCityButton:       "Поменять город",
			ChangeLanguageButton:   "Поменять язык",
			ChooseLanguagePrompt:   "Выберите язык:",
			LanguageSavedMessage:   "Язык сохранён.",
			UnknownMenuAction:      "Неизвестное действие меню.",
			UnknownLanguageMessage: "Неизвестный язык.",
			FindClinicHelpMessage:  "Используй команду: /find_clinic <медицинское направление>",
		}, nil
	case "es":
		return &MenuTexts{
			Prompt:                 "Elige una acción:",
			FindClinicButton:       "Buscar clínica",
			ChangeCityButton:       "Cambiar ciudad",
			ChangeLanguageButton:   "Cambiar idioma",
			ChooseLanguagePrompt:   "Elige un idioma:",
			LanguageSavedMessage:   "Idioma guardado.",
			UnknownMenuAction:      "Acción de menú desconocida.",
			UnknownLanguageMessage: "Idioma desconocido.",
			FindClinicHelpMessage:  "Usa el comando: /find_clinic <especialidad médica>",
		}, nil
	case "en":
		return &MenuTexts{
			Prompt:                 "Choose an action:",
			FindClinicButton:       "Find clinic",
			ChangeCityButton:       "Change city",
			ChangeLanguageButton:   "Change language",
			ChooseLanguagePrompt:   "Choose a language:",
			LanguageSavedMessage:   "Language saved.",
			UnknownMenuAction:      "Unknown menu action.",
			UnknownLanguageMessage: "Unknown language.",
			FindClinicHelpMessage:  "Use command: /find_clinic <medical direction>",
		}, nil
	}
	return &MenuTexts{}, nil
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

func resolveUserLanguage(ctx context.Context, tgBot *bot.Bot, update *models.Update) string {
	params := local_models.NewBaseParams(ctx, tgBot, update)
	dbManager := helpers.GetDbManager(params)
	user := helpers.FetchUser(params, dbManager)
	if user == nil {
		params.Log.Error("User not found")
		return ""
	}
	if user.LanguageCode == "" {
		if languageCode := normalizeLanguage(update.Message.From.LanguageCode); languageCode != "" {
			return languageCode
		}
	}

	return normalizeLanguage(user.LanguageCode)
}

func BuildMainMenu(tgBot *bot.Bot, languageCode string) *reply.ReplyKeyboard {
	texts, err := menuTexts(languageCode)
	if err != nil || texts.Prompt == "" {
		return nil
	}

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
	params := local_models.NewBaseParams(ctx, tgBot, update)
	dbManager := helpers.GetDbManager(params)
	user := helpers.FetchUser(params, dbManager)
	if user == nil {
		params.Log.Error("User not found")
		return errors.New("user not found")
	}
	texts, err := menuTexts(user.LanguageCode)
	if err != nil || texts.Prompt == "" {
		return nil
	}

	_, err = tgBot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        texts.Prompt,
		ReplyMarkup: BuildMainMenu(tgBot, user.LanguageCode),
	})
	if err != nil {
		return fmt.Errorf("send main menu: %w", err)
	}

	return nil
}

func SendLanguageMenu(ctx context.Context, tgBot *bot.Bot, update *models.Update) error {
	params := local_models.NewBaseParams(ctx, tgBot, update)
	dbManager := helpers.GetDbManager(params)
	user := helpers.FetchUser(params, dbManager)
	if user == nil {
		params.Log.Error("User not found")
		return errors.New("user not found")
	}

	texts, err := menuTexts(user.LanguageCode)
	if err != nil || texts.Prompt == "" {
		return nil
	}
	_, err = tgBot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        texts.ChooseLanguagePrompt,
		ReplyMarkup: BuildLanguageMenu(tgBot),
	})
	if err != nil {
		return fmt.Errorf("send language menu: %w", err)
	}

	return nil
}

func ResolveMainMenuAction(params *local_models.BaseParams) (string, bool) {
	dbManager := helpers.GetDbManager(params)
	user := helpers.FetchUser(params, dbManager)
	if user == nil {
		params.Log.Error("User not found")
		return "", false
	}
	texts, err := menuTexts(user.LanguageCode)
	if err != nil || texts.Prompt == "" {
		return "", false
	}
	value := strings.TrimSpace(params.Update.Message.Text)

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

func onMainMenuSelect(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
	params := local_models.NewBaseParams(ctx, tgBot, update)
	dbManager := helpers.GetDbManager(params)
	user := helpers.FetchUser(params, dbManager)
	if user == nil {
		params.Log.Error("User not found")
		return
	}

	languageCode := user.LanguageCode
	texts, err := menuTexts(languageCode)
	if err != nil || texts.Prompt == "" {
		params.Log.Error("Menu text Select Error: ")
		return
	}

	action, ok := ResolveMainMenuAction(params)
	if !ok {
		texts, err = menuTexts(languageCode)
		if err != nil || texts.Prompt == "" {
			params.Log.Error("Menu text Select Error: ")
			return
		}
		_, _ = tgBot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   texts.UnknownMenuAction,
		})
		return
	}

	var response string
	menuAllText, err := menuTexts(languageCode)
	if err != nil || menuAllText == nil {
		params.Log.Error(fmt.Sprintf("languageCode not found %s", languageCode))
		return
	}
	switch action {
	case ActionFindClinic:
		response = menuAllText.FindClinicHelpMessage
	case ActionChangeCity:
		response = menuAllText.ChangeCityButton
	case ActionChangeLanguage:
		response = menuAllText.ChangeLanguageButton
	default:
		response = menuAllText.UnknownMenuAction
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
		texts, err := menuTexts(resolveUserLanguage(ctx, tgBot, update))
		if err != nil || texts.Prompt == "" {
			params.Log.Error(fmt.Sprintf("languageCode not found %s", update.Message.Text))
			return
		}
		_, _ = tgBot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   texts.UnknownLanguageMessage,
		})
		return
	}

	user, err := dbManager.GetUser(ctx, update.Message.Chat.ID)
	if err != nil || user == nil {
		params.Log.Error("user not found")
		return
	}

	user.LanguageCode = newLanguage

	if err = dbManager.InsertOrUpdateUser(ctx, user); err != nil {
		params.Log.Error(err.Error())
		return
	}

	texts, err := menuTexts(newLanguage)
	if err != nil || texts.Prompt == "" {
		params.Log.Error(fmt.Sprintf("languageCode not found %s", newLanguage))
		return
	}

	_, _ = tgBot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        texts.LanguageSavedMessage,
		ReplyMarkup: BuildMainMenu(tgBot, newLanguage),
	})
}
