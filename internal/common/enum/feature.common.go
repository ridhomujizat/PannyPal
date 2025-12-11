package enum

type FeatureType string

const (
	FeatureTypeAIcashflow FeatureType = "AI_CASHFLOW"
)

type WebhookIncomingTag string

const (
	TagKeuangan WebhookIncomingTag = "#keuangan"
)

func TagForFeature(tag WebhookIncomingTag) FeatureType {
	switch tag {
	case TagKeuangan:
		return FeatureTypeAIcashflow
	default:
		return ""
	}
}
