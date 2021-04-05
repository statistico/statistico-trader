package bootstrap

import (
	"github.com/statistico/statistico-strategy/internal/trader/classify"
)

func (c Container) FilterMatcher() classify.FilterMatcher {
	return classify.NewFilterMatcher(c.DataServiceFixtureClient(), c.ResultClassifier(), c.StatClassifier())
}

func (c Container) ResultClassifier() classify.ResultFilterClassifier {
	return classify.NewResultFilterClassifier(c.DataServiceResultClient())
}

func (c Container) StatClassifier() classify.StatFilterClassifier {
	return classify.NewStatFilterClassifier(c.DataServiceResultClient())
}
