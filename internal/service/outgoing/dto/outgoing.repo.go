package dto

type Payloadwaha struct {
	ID          string      `json:"id"`
	Me          Me          `json:"me"`
	Event       string      `json:"event"`
	Engine      string      `json:"engine"`
	Payload     Payload     `json:"payload"`
	Session     string      `json:"session"`
	Metadata    Metadata    `json:"metadata"`
	Timestamp   int64       `json:"timestamp"`
	Environment Environment `json:"environment"`
}

type Environment struct {
	Tier    string `json:"tier"`
	Engine  string `json:"engine"`
	Browser string `json:"browser"`
	Version string `json:"version"`
}

type Me struct {
	ID       string `json:"id"`
	PushName string `json:"pushName"`
}

type Metadata struct {
}

type Payload struct {
	ID          string        `json:"id"`
	To          string        `json:"to"`
	ACK         int64         `json:"ack"`
	Body        string        `json:"body"`
	From        string        `json:"from"`
	Data        Data          `json:"_data"`
	Media       interface{}   `json:"media"`
	FromMe      bool          `json:"fromMe"`
	Source      string        `json:"source"`
	VCards      []interface{} `json:"vCards"`
	ACKName     string        `json:"ackName"`
	ReplyTo     *ReplyTo      `json:"replyTo"`
	HasMedia    bool          `json:"hasMedia"`
	Location    interface{}   `json:"location"`
	Timestamp   int64         `json:"timestamp"`
	Participant string        `json:"participant"`
}

type Data struct {
	T                                     int64              `json:"t"`
	ID                                    ID                 `json:"id"`
	To                                    string             `json:"to"`
	ACK                                   int64              `json:"ack"`
	Body                                  string             `json:"body"`
	From                                  string             `json:"from"`
	Star                                  bool               `json:"star"`
	Type                                  string             `json:"type"`
	Invis                                 bool               `json:"invis"`
	Links                                 []interface{}      `json:"links"`
	Author                                string             `json:"author"`
	Viewed                                bool               `json:"viewed"`
	IsAvatar                              bool               `json:"isAvatar"`
	IsNewMsg                              bool               `json:"isNewMsg"`
	ViewMode                              string             `json:"viewMode"`
	QuotedMsg                             QuotedMsg          `json:"quotedMsg"`
	RecvFresh                             bool               `json:"recvFresh"`
	Thumbnail                             string             `json:"thumbnail"`
	BizBotType                            interface{}        `json:"bizBotType"`
	IsAdsMedia                            bool               `json:"isAdsMedia"`
	IsCallLink                            interface{}        `json:"isCallLink"`
	NotifyName                            string             `json:"notifyName"`
	QuotedType                            int64              `json:"quotedType"`
	CallCreator                           interface{}        `json:"callCreator"`
	HasReaction                           bool               `json:"hasReaction"`
	IsVideoCall                           bool               `json:"isVideoCall"`
	KicNotified                           bool               `json:"kicNotified"`
	ParentMsgID                           interface{}        `json:"parentMsgId"`
	CallDuration                          interface{}        `json:"callDuration"`
	BotPluginType                         interface{}        `json:"botPluginType"`
	CallLinkToken                         interface{}        `json:"callLinkToken"`
	ForwardsCount                         int64              `json:"forwardsCount"`
	GroupMentions                         []interface{}      `json:"groupMentions"`
	InvokedBotWid                         interface{}        `json:"invokedBotWid"`
	MessageSecret                         map[string]int64   `json:"messageSecret"`
	StickerSentTs                         int64              `json:"stickerSentTs"`
	BotMsgBodyType                        interface{}        `json:"botMsgBodyType"`
	IsCarouselCard                        bool               `json:"isCarouselCard"`
	IsFromTemplate                        bool               `json:"isFromTemplate"`
	IsMdHistoryMsg                        bool               `json:"isMdHistoryMsg"`
	QuotedStanzaID                        string             `json:"quotedStanzaID"`
	IsEventCanceled                       bool               `json:"isEventCanceled"`
	PollInvalidated                       bool               `json:"pollInvalidated"`
	CallParticipants                      interface{}        `json:"callParticipants"`
	EventInvalidated                      bool               `json:"eventInvalidated"`
	LatestEditMsgKey                      interface{}        `json:"latestEditMsgKey"`
	MentionedJidList                      []interface{}      `json:"mentionedJidList"`
	CallSilenceReason                     interface{}        `json:"callSilenceReason"`
	QuotedParticipant                     string             `json:"quotedParticipant"`
	BotPluginSearchURL                    interface{}        `json:"botPluginSearchUrl"`
	FaviconMMSMetadata                    interface{}        `json:"faviconMMSMetadata"`
	ReportingTokenInfo                    ReportingTokenInfo `json:"reportingTokenInfo"`
	BotResponseTargetID                   interface{}        `json:"botResponseTargetId"`
	BotPluginMaybeParent                  bool               `json:"botPluginMaybeParent"`
	BotPluginSearchQuery                  interface{}        `json:"botPluginSearchQuery"`
	LastPlaybackProgress                  int64              `json:"lastPlaybackProgress"`
	IsSentCagPollCreation                 bool               `json:"isSentCagPollCreation"`
	ClientReceivedTsMillis                int64              `json:"clientReceivedTsMillis"`
	IsVcardOverMmsDocument                bool               `json:"isVcardOverMmsDocument"`
	LastUpdateFromServerTs                int64              `json:"lastUpdateFromServerTs"`
	QuestionResponsesCount                int64              `json:"questionResponsesCount"`
	BotPluginReferenceIndex               interface{}        `json:"botPluginReferenceIndex"`
	BotPluginSearchProvider               interface{}        `json:"botPluginSearchProvider"`
	BotMessageDisclaimerText              interface{}        `json:"botMessageDisclaimerText"`
	IsDynamicReplyButtonsMsg              bool               `json:"isDynamicReplyButtonsMsg"`
	RequiresDirectConnection              interface{}        `json:"requiresDirectConnection"`
	BizContentPlaceholderType             interface{}        `json:"bizContentPlaceholderType"`
	HostedBizEncStateMismatch             bool               `json:"hostedBizEncStateMismatch"`
	GroupHistoryBundleMetadata            interface{}        `json:"groupHistoryBundleMetadata"`
	ProductHeaderImageRejected            bool               `json:"productHeaderImageRejected"`
	QuestionReplyQuotedMessage            interface{}        `json:"questionReplyQuotedMessage"`
	ReadQuestionResponsesCount            int64              `json:"readQuestionResponsesCount"`
	LatestEditSenderTimestampMS           interface{}        `json:"latestEditSenderTimestampMs"`
	BotReelPluginThumbnailCDNURL          interface{}        `json:"botReelPluginThumbnailCdnUrl"`
	GroupHistoryBundleMessageKey          interface{}        `json:"groupHistoryBundleMessageKey"`
	SenderOrRecipientAccountTypeHosted    bool               `json:"senderOrRecipientAccountTypeHosted"`
	PlaceholderCreatedWhenAccountIsHosted bool               `json:"placeholderCreatedWhenAccountIsHosted"`
}

type ID struct {
	ID          string `json:"id"`
	FromMe      bool   `json:"fromMe"`
	Remote      string `json:"remote"`
	Serialized  string `json:"_serialized"`
	Participant string `json:"participant"`
}

type QuotedMsg struct {
	Body string `json:"body"`
	Kind string `json:"kind"`
	Type string `json:"type"`
}

type ReportingTokenInfo struct {
	Version        int64            `json:"version"`
	ReportingTag   map[string]int64 `json:"reportingTag"`
	ReportingToken map[string]int64 `json:"reportingToken"`
}

type ReplyTo struct {
	Body string    `json:"body"`
	Data QuotedMsg `json:"_data"`
}
