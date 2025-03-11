package whatsapp

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
)

func NewTemplateMessage(name string, languageCode string) *TemplateMessage {
	return &TemplateMessage{
		Name:     name,
		Language: TemplateMessageLanguage{languageCode},
	}
}

func (t *TemplateMessage) AddComponent(component TemplateMessageComponent) {
	t.Components = append(t.Components, component)
}

func (service *Service) SendTemplateMessage(fromPhoneNumberId string, to string, templateMessage *TemplateMessage) *errortools.Error {
	message := Message{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               to,
		Type:             "template",
		Template:         templateMessage,
	}

	var response SendMessageResponse

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPost,
		Url:           service.url(fmt.Sprintf("%s/messages", fromPhoneNumberId)),
		BodyModel:     message,
		ResponseModel: &response,
	}
	_, _, e := service.httpRequest(&requestConfig)

	return e
}

type Message struct {
	MessagingProduct string           `json:"messaging_product"`
	RecipientType    string           `json:"recipient_type"`
	To               string           `json:"to"`
	Type             string           `json:"type"`
	Template         *TemplateMessage `json:"template,omitempty"`
}

type TemplateMessage struct {
	Name       string                     `json:"name"`
	Language   TemplateMessageLanguage    `json:"language"`
	Components []TemplateMessageComponent `json:"components"`
}

type TemplateMessageLanguage struct {
	Code string `json:"code"`
}

type TemplateMessageComponent struct {
	Type       string                              `json:"type"`
	Parameters []TemplateMessageComponentParameter `json:"parameters"`
	SubType    string                              `json:"sub_type,omitempty"`
	Index      string                              `json:"index,omitempty"`
}

type TemplateMessageComponentParameter struct {
	Type  string `json:"type"`
	Image *struct {
		Link string `json:"link"`
	} `json:"image,omitempty"`
	Text     string `json:"text,omitempty"`
	Currency *struct {
		FallbackValue string      `json:"fallback_value"`
		Code          string      `json:"code"`
		Amount1000    interface{} `json:"amount_1000"`
	} `json:"currency,omitempty"`
	DateTime *struct {
		FallbackValue string `json:"fallback_value"`
	} `json:"date_time,omitempty"`
	Payload *string `json:"payload,omitempty"`
}

type SendMessageResponse struct {
	MessagingProduct string `json:"messaging_product"`
	Contacts         []struct {
		Input string `json:"input"`
		WaId  string `json:"wa_id"`
	} `json:"contacts"`
	Messages []struct {
		Id            string `json:"id"`
		MessageStatus string `json:"message_status"`
	} `json:"messages"`
}
