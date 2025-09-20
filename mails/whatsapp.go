package mails

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pocketbase/pocketbase/core"
)

// WhatsAppConfig holds the configuration for WhatsApp Business API
type WhatsAppConfig struct {
	AccessToken string `json:"accessToken"`
	PhoneNumberID string `json:"phoneNumberId"`
	APIURL      string `json:"apiUrl"`
}

// WhatsAppMessage represents a WhatsApp message structure
type WhatsAppMessage struct {
	MessagingProduct string `json:"messaging_product"`
	To               string `json:"to"`
	Type             string `json:"type"`
	Template         *WhatsAppTemplate `json:"template,omitempty"`
	Text             *WhatsAppText     `json:"text,omitempty"`
}

// WhatsAppTemplate represents a WhatsApp template message
type WhatsAppTemplate struct {
	Name       string                 `json:"name"`
	Language   WhatsAppLanguage       `json:"language"`
	Components []WhatsAppComponent    `json:"components,omitempty"`
}

// WhatsAppText represents a simple text message
type WhatsAppText struct {
	Body string `json:"body"`
}

// WhatsAppLanguage represents the language for template messages
type WhatsAppLanguage struct {
	Code string `json:"code"`
}

// WhatsAppComponent represents a component in a template message
type WhatsAppComponent struct {
	Type       string                 `json:"type"`
	Parameters []WhatsAppParameter    `json:"parameters,omitempty"`
}

// WhatsAppParameter represents a parameter in a template component
type WhatsAppParameter struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// WhatsAppResponse represents the response from WhatsApp API
type WhatsAppResponse struct {
	MessagingProduct string `json:"messaging_product"`
	Contacts         []struct {
		Input string `json:"input"`
		WaID  string `json:"wa_id"`
	} `json:"contacts"`
	Messages []struct {
		ID string `json:"id"`
	} `json:"messages"`
}

// SendWhatsAppMessage sends a message via WhatsApp Business API
func SendWhatsAppMessage(config WhatsAppConfig, to, message string) error {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	whatsappMsg := WhatsAppMessage{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "text",
		Text: &WhatsAppText{
			Body: message,
		},
	}

	jsonData, err := json.Marshal(whatsappMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal WhatsApp message: %w", err)
	}

	url := fmt.Sprintf("%s/%s/messages", config.APIURL, config.PhoneNumberID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create WhatsApp request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+config.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send WhatsApp message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("WhatsApp API returned status %d", resp.StatusCode)
	}

	return nil
}

// SendRecordOTPWhatsApp sends OTP via WhatsApp to the specified auth record
func SendRecordOTPWhatsApp(app core.App, authRecord *core.Record, otpId string, pass string) error {
	// Get WhatsApp configuration from app settings
	// This would need to be added to the settings model
	config := WhatsAppConfig{
		AccessToken:   app.Settings().Meta.WhatsAppAccessToken, // This field needs to be added
		PhoneNumberID: app.Settings().Meta.WhatsAppPhoneNumberID, // This field needs to be added
		APIURL:        "https://graph.facebook.com/v18.0",
	}

	// Get phone number from record (assuming there's a phone field)
	phoneNumber := authRecord.GetString("phone")
	if phoneNumber == "" {
		return fmt.Errorf("no phone number found for record %s", authRecord.Id)
	}

	// Resolve WhatsApp template
	message, err := resolveWhatsAppTemplate(app, authRecord, authRecord.Collection().OTP.WhatsAppTemplate, map[string]any{
		core.EmailPlaceholderOTPId: otpId,
		core.EmailPlaceholderOTP:   pass,
	})
	if err != nil {
		return err
	}

	// Send WhatsApp message
	err = SendWhatsAppMessage(config, phoneNumber, message)
	if err != nil {
		return fmt.Errorf("failed to send WhatsApp OTP: %w", err)
	}

	// Update OTP sentTo field
	otp, err := app.FindOTPById(otpId)
	if err != nil {
		app.Logger().Warn(
			"Unable to find OTP to update its sentTo field",
			"error", err,
			"otpId", otpId,
		)
		return nil
	}

	if otp.SentTo() != "" {
		return nil // was already sent to another target
	}

	otp.SetSentTo(phoneNumber)
	if err = app.Save(otp); err != nil {
		app.Logger().Error(
			"Failed to update OTP sentTo field",
			"error", err,
			"otpId", otpId,
			"to", phoneNumber,
		)
	}

	return nil
}

// resolveWhatsAppTemplate resolves the WhatsApp message template with placeholders
func resolveWhatsAppTemplate(app core.App, authRecord *core.Record, template core.MessageTemplate, placeholders map[string]any) (string, error) {
	// Add common placeholders
	placeholders[core.EmailPlaceholderAppName] = app.Settings().Meta.AppName
	placeholders[core.EmailPlaceholderAppURL] = app.Settings().Meta.AppURL
	placeholders["{RECORD_ID}"] = authRecord.Id
	placeholders["{RECORD_EMAIL}"] = authRecord.Email()

	// Resolve template
	message := template.Resolve(placeholders)

	return message, nil
}
