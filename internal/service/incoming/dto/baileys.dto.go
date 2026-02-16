package dto

// ===============================
// Simplified Incoming Message
// ===============================

// SimplifiedIncomingMessage represents the clean webhook payload format
type SimplifiedIncomingMessage struct {
	// Session & Basic Info
	SessionID string `json:"sessionId"`
	MessageID string `json:"messageId"`
	Timestamp int64  `json:"timestamp"`

	// Sender Info
	From     string  `json:"from"`     // Clean phone number (no @s.whatsapp.net)
	FromName *string `json:"fromName"` // Sender's WhatsApp display name

	// Message Type
	MessageType string `json:"messageType"` // text, image, video, audio, document, location, contact, sticker, unknown

	// Chat Context
	IsGroup     bool    `json:"isGroup"`
	ChatID      string  `json:"chatId"`      // Full JID for reference
	Participant *string `json:"participant"` // Sender's number in group (for groups only)

	// Content (varies by message type)
	Content MessageContent `json:"content"`

	// Reply Context (optional)
	QuotedMessage *QuotedMessageInfo `json:"quotedMessage,omitempty"`
}

// ===============================
// Message Content
// ===============================

// MessageContent represents the message content (varies by type)
type MessageContent struct {
	Type string `json:"type"` // text, image, video, audio, document, location, contact, sticker, unknown

	// Text content
	Text *string `json:"text,omitempty"`

	// Media content (image, video, audio, document, sticker)
	Caption              *string               `json:"caption,omitempty"`
	Mimetype             *string               `json:"mimetype,omitempty"`
	FileName             *string               `json:"fileName,omitempty"`
	FileSize             *int64                `json:"fileSize,omitempty"`
	Duration             *int                  `json:"duration,omitempty"`
	DownloadInstructions *DownloadInstructions `json:"downloadInstructions,omitempty"`

	// Location content
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
	Name      *string  `json:"name,omitempty"`
	Address   *string  `json:"address,omitempty"`

	// Contact content
	DisplayName *string   `json:"displayName,omitempty"`
	Vcard       *string   `json:"vcard,omitempty"`
	Phones      *[]string `json:"phones,omitempty"`

	// Unknown content
	Description *string `json:"description,omitempty"`
}

// ===============================
// Download Instructions
// ===============================

type DownloadInstructions struct {
	Method   string                 `json:"method"`
	Endpoint string                 `json:"endpoint"`
	Body     map[string]interface{} `json:"body"`
}

// ===============================
// Quoted Message Info
// ===============================

type QuotedMessageInfo struct {
	MessageID string  `json:"messageId"`
	Text      *string `json:"text,omitempty"`
}

// ===============================
// Helper Methods
// ===============================

// GetText extracts text from the message content
func (m *SimplifiedIncomingMessage) GetText() string {
	if m.Content.Text != nil {
		return *m.Content.Text
	}
	if m.Content.Caption != nil {
		return *m.Content.Caption
	}
	return ""
}
