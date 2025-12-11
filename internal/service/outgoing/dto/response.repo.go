package dto

type ResponseOutgoingwaha struct {
	Data            Data          `json:"_data"`
	ID              ID            `json:"id"`
	ACK             int64         `json:"ack"`
	HasMedia        bool          `json:"hasMedia"`
	Body            string        `json:"body"`
	Type            string        `json:"type"`
	Timestamp       int64         `json:"timestamp"`
	From            string        `json:"from"`
	To              string        `json:"to"`
	DeviceType      string        `json:"deviceType"`
	IsForwarded     bool          `json:"isForwarded"`
	ForwardingScore int64         `json:"forwardingScore"`
	IsStatus        bool          `json:"isStatus"`
	IsStarred       bool          `json:"isStarred"`
	FromMe          bool          `json:"fromMe"`
	HasQuotedMsg    bool          `json:"hasQuotedMsg"`
	HasReaction     bool          `json:"hasReaction"`
	VCards          []interface{} `json:"vCards"`
	MentionedIDS    []interface{} `json:"mentionedIds"`
	GroupMentions   []interface{} `json:"groupMentions"`
	IsGIF           bool          `json:"isGif"`
	Links           []interface{} `json:"links"`
}

type Data struct {
	ID                                    ID            `json:"id"`
	Viewed                                bool          `json:"viewed"`
	Body                                  string        `json:"body"`
	Type                                  string        `json:"type"`
	T                                     int64         `json:"t"`
	From                                  From          `json:"from"`
	To                                    From          `json:"to"`
	ACK                                   int64         `json:"ack"`
	IsNewMsg                              bool          `json:"isNewMsg"`
	Star                                  bool          `json:"star"`
	KicNotified                           bool          `json:"kicNotified"`
	IsFromTemplate                        bool          `json:"isFromTemplate"`
	IsAdsMedia                            bool          `json:"isAdsMedia"`
	PollInvalidated                       bool          `json:"pollInvalidated"`
	IsSentCagPollCreation                 bool          `json:"isSentCagPollCreation"`
	LatestEditMsgKey                      interface{}   `json:"latestEditMsgKey"`
	LatestEditSenderTimestampMS           interface{}   `json:"latestEditSenderTimestampMs"`
	QuotedMsg                             QuotedMsg     `json:"quotedMsg"`
	QuotedStanzaID                        string        `json:"quotedStanzaID"`
	QuotedRemoteJid                       From          `json:"quotedRemoteJid"`
	QuotedParticipant                     From          `json:"quotedParticipant"`
	MentionedJidList                      []interface{} `json:"mentionedJidList"`
	GroupMentions                         []interface{} `json:"groupMentions"`
	IsEventCanceled                       bool          `json:"isEventCanceled"`
	EventInvalidated                      bool          `json:"eventInvalidated"`
	IsVcardOverMmsDocument                bool          `json:"isVcardOverMmsDocument"`
	IsForwarded                           bool          `json:"isForwarded"`
	IsQuestion                            bool          `json:"isQuestion"`
	QuestionReplyQuotedMessage            interface{}   `json:"questionReplyQuotedMessage"`
	QuestionResponsesCount                int64         `json:"questionResponsesCount"`
	ReadQuestionResponsesCount            int64         `json:"readQuestionResponsesCount"`
	ForwardsCount                         int64         `json:"forwardsCount"`
	HasReaction                           bool          `json:"hasReaction"`
	DisappearingModeInitiator             string        `json:"disappearingModeInitiator"`
	DisappearingModeTrigger               string        `json:"disappearingModeTrigger"`
	DisappearingModeInitiatedByMe         bool          `json:"disappearingModeInitiatedByMe"`
	ProductHeaderImageRejected            bool          `json:"productHeaderImageRejected"`
	LastPlaybackProgress                  int64         `json:"lastPlaybackProgress"`
	IsDynamicReplyButtonsMsg              bool          `json:"isDynamicReplyButtonsMsg"`
	IsCarouselCard                        bool          `json:"isCarouselCard"`
	ParentMsgID                           interface{}   `json:"parentMsgId"`
	CallSilenceReason                     interface{}   `json:"callSilenceReason"`
	IsVideoCall                           bool          `json:"isVideoCall"`
	CallDuration                          interface{}   `json:"callDuration"`
	CallCreator                           interface{}   `json:"callCreator"`
	CallParticipants                      interface{}   `json:"callParticipants"`
	IsCallLink                            interface{}   `json:"isCallLink"`
	CallLinkToken                         interface{}   `json:"callLinkToken"`
	IsMdHistoryMsg                        bool          `json:"isMdHistoryMsg"`
	StickerSentTs                         int64         `json:"stickerSentTs"`
	IsAvatar                              bool          `json:"isAvatar"`
	LastUpdateFromServerTs                int64         `json:"lastUpdateFromServerTs"`
	InvokedBotWid                         interface{}   `json:"invokedBotWid"`
	BizBotType                            interface{}   `json:"bizBotType"`
	BotResponseTargetID                   interface{}   `json:"botResponseTargetId"`
	BotPluginType                         interface{}   `json:"botPluginType"`
	BotPluginReferenceIndex               interface{}   `json:"botPluginReferenceIndex"`
	BotPluginSearchProvider               interface{}   `json:"botPluginSearchProvider"`
	BotPluginSearchURL                    interface{}   `json:"botPluginSearchUrl"`
	BotPluginSearchQuery                  interface{}   `json:"botPluginSearchQuery"`
	BotPluginMaybeParent                  bool          `json:"botPluginMaybeParent"`
	BotReelPluginThumbnailCDNURL          interface{}   `json:"botReelPluginThumbnailCdnUrl"`
	BotMessageDisclaimerText              interface{}   `json:"botMessageDisclaimerText"`
	BotMsgBodyType                        interface{}   `json:"botMsgBodyType"`
	RequiresDirectConnection              bool          `json:"requiresDirectConnection"`
	BizContentPlaceholderType             interface{}   `json:"bizContentPlaceholderType"`
	HostedBizEncStateMismatch             bool          `json:"hostedBizEncStateMismatch"`
	SenderOrRecipientAccountTypeHosted    bool          `json:"senderOrRecipientAccountTypeHosted"`
	PlaceholderCreatedWhenAccountIsHosted bool          `json:"placeholderCreatedWhenAccountIsHosted"`
	GroupHistoryBundleMessageKey          interface{}   `json:"groupHistoryBundleMessageKey"`
	GroupHistoryBundleMetadata            interface{}   `json:"groupHistoryBundleMetadata"`
	NonJidMentions                        interface{}   `json:"nonJidMentions"`
	Links                                 []interface{} `json:"links"`
}

type From struct {
	Server     string `json:"server"`
	User       string `json:"user"`
	Serialized string `json:"_serialized"`
}

type ID struct {
	FromMe      bool   `json:"fromMe"`
	Remote      string `json:"remote"`
	ID          string `json:"id"`
	Participant From   `json:"participant"`
	Serialized  string `json:"_serialized"`
}

type QuotedMsg struct {
	Viewed                                bool               `json:"viewed"`
	Body                                  string             `json:"body"`
	Type                                  string             `json:"type"`
	ClientReceivedTsMillis                int64              `json:"clientReceivedTsMillis"`
	KicNotified                           bool               `json:"kicNotified"`
	IsFromTemplate                        bool               `json:"isFromTemplate"`
	Thumbnail                             string             `json:"thumbnail"`
	FaviconMMSMetadata                    interface{}        `json:"faviconMMSMetadata"`
	IsAdsMedia                            bool               `json:"isAdsMedia"`
	PollInvalidated                       bool               `json:"pollInvalidated"`
	IsSentCagPollCreation                 bool               `json:"isSentCagPollCreation"`
	MentionedJidList                      []interface{}      `json:"mentionedJidList"`
	GroupMentions                         []interface{}      `json:"groupMentions"`
	IsEventCanceled                       bool               `json:"isEventCanceled"`
	EventInvalidated                      bool               `json:"eventInvalidated"`
	IsVcardOverMmsDocument                bool               `json:"isVcardOverMmsDocument"`
	QuestionReplyQuotedMessage            interface{}        `json:"questionReplyQuotedMessage"`
	QuestionResponsesCount                int64              `json:"questionResponsesCount"`
	ReadQuestionResponsesCount            int64              `json:"readQuestionResponsesCount"`
	ForwardsCount                         int64              `json:"forwardsCount"`
	HasReaction                           bool               `json:"hasReaction"`
	ViewMode                              string             `json:"viewMode"`
	MessageSecret                         map[string]int64   `json:"messageSecret"`
	ProductHeaderImageRejected            bool               `json:"productHeaderImageRejected"`
	LastPlaybackProgress                  int64              `json:"lastPlaybackProgress"`
	IsDynamicReplyButtonsMsg              bool               `json:"isDynamicReplyButtonsMsg"`
	IsCarouselCard                        bool               `json:"isCarouselCard"`
	ParentMsgID                           interface{}        `json:"parentMsgId"`
	CallSilenceReason                     interface{}        `json:"callSilenceReason"`
	IsVideoCall                           bool               `json:"isVideoCall"`
	CallDuration                          interface{}        `json:"callDuration"`
	CallCreator                           interface{}        `json:"callCreator"`
	CallParticipants                      interface{}        `json:"callParticipants"`
	IsCallLink                            interface{}        `json:"isCallLink"`
	CallLinkToken                         interface{}        `json:"callLinkToken"`
	IsMdHistoryMsg                        bool               `json:"isMdHistoryMsg"`
	StickerSentTs                         int64              `json:"stickerSentTs"`
	IsAvatar                              bool               `json:"isAvatar"`
	LastUpdateFromServerTs                int64              `json:"lastUpdateFromServerTs"`
	InvokedBotWid                         interface{}        `json:"invokedBotWid"`
	BizBotType                            interface{}        `json:"bizBotType"`
	BotResponseTargetID                   interface{}        `json:"botResponseTargetId"`
	BotPluginType                         interface{}        `json:"botPluginType"`
	BotPluginReferenceIndex               interface{}        `json:"botPluginReferenceIndex"`
	BotPluginSearchProvider               interface{}        `json:"botPluginSearchProvider"`
	BotPluginSearchURL                    interface{}        `json:"botPluginSearchUrl"`
	BotPluginSearchQuery                  interface{}        `json:"botPluginSearchQuery"`
	BotPluginMaybeParent                  bool               `json:"botPluginMaybeParent"`
	BotReelPluginThumbnailCDNURL          interface{}        `json:"botReelPluginThumbnailCdnUrl"`
	BotMessageDisclaimerText              interface{}        `json:"botMessageDisclaimerText"`
	BotMsgBodyType                        interface{}        `json:"botMsgBodyType"`
	ReportingTokenInfo                    ReportingTokenInfo `json:"reportingTokenInfo"`
	RequiresDirectConnection              bool               `json:"requiresDirectConnection"`
	BizContentPlaceholderType             interface{}        `json:"bizContentPlaceholderType"`
	HostedBizEncStateMismatch             bool               `json:"hostedBizEncStateMismatch"`
	SenderOrRecipientAccountTypeHosted    bool               `json:"senderOrRecipientAccountTypeHosted"`
	PlaceholderCreatedWhenAccountIsHosted bool               `json:"placeholderCreatedWhenAccountIsHosted"`
	GroupHistoryBundleMessageKey          interface{}        `json:"groupHistoryBundleMessageKey"`
	GroupHistoryBundleMetadata            interface{}        `json:"groupHistoryBundleMetadata"`
}

type ReportingTokenInfo struct {
	ReportingToken map[string]int64 `json:"reportingToken"`
	Version        int64            `json:"version"`
	ReportingTag   map[string]int64 `json:"reportingTag"`
}
