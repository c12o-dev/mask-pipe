package patterns

import (
	"regexp"
	"strings"
)

var stripeKey = &Pattern{
	ID:         "stripe_key",
	Name:       "Stripe API Key",
	Regex:      regexp.MustCompile(`\b[sp]k_(?:live|test)_[A-Za-z0-9]{24,}`),
	CaptureIdx: 0,
	Examples:   buildStripeExamples(),
	NonExamples: []string{
		"sk_live_short",
		"rk_live_" + strings.Repeat("A", 25),
		"sk_prod_" + strings.Repeat("A", 25),
		"SK_LIVE_" + strings.Repeat("A", 25),
		"disk_test_1a2b3c4d5e6f7890abcdef123456",
		"sk_test_" + strings.Repeat("A", 23),
	},
	Source: "https://docs.stripe.com/keys",
}

func buildStripeExamples() []string {
	body := strings.Repeat("A", 25)
	return []string{
		"sk_live_" + body,
		"sk_test_" + body,
		"pk_live_" + body,
		"pk_test_" + body,
		"sk_live_" + strings.Repeat("Ab0C", 7)[:25],
	}
}
