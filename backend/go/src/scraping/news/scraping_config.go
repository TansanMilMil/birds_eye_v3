package news

import "github.com/birdseyeapi/birdseyeapi_v2/go/src/env"

func init() {
	if limit := env.GetEnvInt("SCRAPING_ARTICLES", 0); limit > 0 {
		MaxArticles = limit
		HatenaMaxArticles = limit
		CloudWatchMaxArticles = limit
		ZDNetMaxArticles = limit
	}
}
