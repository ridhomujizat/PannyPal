package dto

// ===============================
// Root Incoming Payload
// ===============================

type BaileysIncomingMessage struct {
	SessionID   string             `json:"sessionId"`
	From        string             `json:"from"`
	MessageType string             `json:"messageType"`
	Message     BaileysMessageBody `json:"message"`
	Timestamp   int64              `json:"timestamp"`
	Key         BaileysMessageKey  `json:"key"`
}

// ===============================
// Message Key
// ===============================

type BaileysMessageKey struct {
	RemoteJid   string  `json:"remoteJid"`
	ID          string  `json:"id"`
	FromMe      bool    `json:"fromMe"`
	Participant *string `json:"participant,omitempty"` // group only
}

// ===============================
// Message Body
// ===============================

type BaileysMessageBody struct {
	Conversation        string                      `json:"conversation,omitempty"`
	ExtendedTextMessage *BaileysExtendedTextMessage `json:"extendedTextMessage,omitempty"`
	MessageContextInfo  *BaileysMessageContextInfo  `json:"messageContextInfo,omitempty"`
}

// ===============================
// Extended Text Message
// ===============================

type BaileysExtendedTextMessage struct {
	Text                  string              `json:"text"`
	ContextInfo           *BaileysContextInfo `json:"contextInfo,omitempty"`
	InviteLinkGroupTypeV2 string              `json:"inviteLinkGroupTypeV2,omitempty"`
}

// ===============================
// Context Info (Reply / Ephemeral)
// ===============================

type BaileysContextInfo struct {
	StanzaID                  string                   `json:"stanzaId,omitempty"`
	Participant               string                   `json:"participant,omitempty"`
	QuotedMessage             *BaileysQuotedMessage    `json:"quotedMessage,omitempty"`
	EphemeralSettingTimestamp string                   `json:"ephemeralSettingTimestamp,omitempty"`
	DisappearingMode          *BaileysDisappearingMode `json:"disappearingMode,omitempty"`
}

// ===============================
// Quoted / Reply Message
// ===============================

type BaileysQuotedMessage struct {
	Conversation       string                 `json:"conversation,omitempty"`
	MessageContextInfo map[string]interface{} `json:"messageContextInfo,omitempty"`
}

// ===============================
// Disappearing Mode
// ===============================

type BaileysDisappearingMode struct {
	Initiator     string `json:"initiator"`
	Trigger       string `json:"trigger"`
	InitiatedByMe bool   `json:"initiatedByMe"`
}

// ===============================
// Message Context Metadata
// ===============================

type BaileysMessageContextInfo struct {
	DeviceListMetadata        *BaileysDeviceListMetadata `json:"deviceListMetadata,omitempty"`
	DeviceListMetadataVersion int                        `json:"deviceListMetadataVersion,omitempty"`
	MessageSecret             string                     `json:"messageSecret,omitempty"`
	LimitSharingV2            *BaileysLimitSharingV2     `json:"limitSharingV2,omitempty"`
}

// ===============================
// Device Metadata
// ===============================

type BaileysDeviceListMetadata struct {
	SenderKeyHash       string `json:"senderKeyHash"`
	SenderTimestamp     string `json:"senderTimestamp"`
	SenderAccountType   string `json:"senderAccountType"`
	ReceiverAccountType string `json:"receiverAccountType"`
	RecipientKeyHash    string `json:"recipientKeyHash"`
	RecipientTimestamp  string `json:"recipientTimestamp"`
}

// ===============================
// Limit Sharing
// ===============================

type BaileysLimitSharingV2 struct {
	SharingLimited               bool   `json:"sharingLimited"`
	Trigger                      string `json:"trigger"`
	LimitSharingSettingTimestamp string `json:"limitSharingSettingTimestamp"`
	InitiatedByMe                bool   `json:"initiatedByMe"`
}

func (m *BaileysIncomingMessage) GetText() string {
	if m.Message.Conversation != "" {
		return m.Message.Conversation
	}

	if m.Message.ExtendedTextMessage != nil {
		return m.Message.ExtendedTextMessage.Text
	}

	return ""
}
