package twitter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

type Handler struct {
	*http.Client
	PingChannel      chan MonitorPing
	MonitorKillMap   map[int64]context.CancelFunc
	CurrentMonitored []int64

	// Secrets
	CSRFToken   string
	BearerToken string
	AuthToken   string
}

type MonitorFilter struct {
	PositiveKeywords []string
	NegativeKeywords []string
	TwitterId        int64
	LatestTweetTS    int64
}

type MonitorPing struct {
	Handle     string
	Title 	   string
	Message    string
	Image      string
	URL        string
	ParsedData []string // URLSs, public keys, etc!
}

var (
	URLRegex 	   = regexp.MustCompile(`https?:\/\/(?:www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b(?:[-a-zA-Z0-9()@:%_\+.~#?&\/=]*)`)
	BTCPubkeyRegex = regexp.MustCompile(``)
	ETHPubkeyRegex = regexp.MustCompile(`^0x[a-fA-F0-9]{40}$`)
	SOLPubkeyRegex = regexp.MustCompile(`[1-9A-HJ-NP-Za-km-z]{32,44}`)
)

func New(pingChan chan MonitorPing) *Handler {
	return &Handler{
		Client:           &http.Client{},
		PingChannel:      pingChan,
		MonitorKillMap:   make(map[int64]context.CancelFunc),
		CurrentMonitored: []int64{},
		CSRFToken:        os.Getenv("CSRF_TOKEN"),
		BearerToken:      os.Getenv("BEARER_TOKEN"),
		AuthToken:        os.Getenv("AUTH_TOKEN"),
	}
}

func (m *Handler) FetchTwitterID(twitterHandle string) int64 {

	req, err := http.NewRequest("GET", "https://twitter.com/i/api/graphql/ptQPCD7NrFS_TW71Lq07nw/UserByScreenName?variables=%7B%22screen_name%22%3A%22"+twitterHandle+"%22%2C%22withSafetyModeUserFields%22%3Atrue%2C%22withSuperFollowsUserFields%22%3Atrue%7D&features=%7B%22responsive_web_twitter_blue_verified_badge_is_enabled%22%3Atrue%2C%22verified_phone_label_enabled%22%3Afalse%2C%22responsive_web_graphql_timeline_navigation_enabled%22%3Atrue%7D", nil)

	if err != nil {
		return 0
	}

	req.Header = http.Header{
		"sec-ch-ua":                 {`"Chromium";v="112", "Google Chrome";v="112", "Not:A-Brand";v="99"'`},
		"x-twitter-client-language": {"en"},
		"x-csrf-token":              {m.CSRFToken},
		"sec-ch-ua-mobile":          {"en"},
		"authorization":             {"Bearer " + m.BearerToken},
		"user-agent":                {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36"},
		"content-type":              {"application/json"},
		"cookie":                    {fmt.Sprintf("auth_token=%s; ct0=%s;", m.AuthToken, m.CSRFToken)},
		"referer":                   {"https://twitter.com/home"},
		"x-twitter-auth-type":       {"OAuth2Session"},
		"x-twitter-active-user":     {"yes"},
		"sec-ch-ua-platform":        {`"macOS"`},
	}

	res, err := m.Client.Do(req)

	if err != nil {
		return 0
	}

	if res.StatusCode != 200 {
		return 0
	}

	bBytes, _ := io.ReadAll(res.Body)
	defer res.Body.Close()

	var twitterInfo TwitterDataResponse
	if err := json.Unmarshal(bBytes, &twitterInfo); err != nil {
		return 0
	}

	restId, _ := strconv.ParseInt(twitterInfo.Data.User.Result.RestId, 10, 64)
	return restId
}

func ParseExtras(text string) []string {
	return []string{}
}

func MonitorTweets(m *Handler, filter MonitorFilter) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("https://api.twitter.com/1.1/statuses/user_timeline.json?count=5&include_rts=0&user_id=%d&tweet_mode=extended", filter.TwitterId), nil)
	req.Header = http.Header{
		"sec-ch-ua":                 {`"Chromium";v="112", "Google Chrome";v="112", "Not:A-Brand";v="99"'`},
		"x-twitter-client-language": {"en"},
		"x-csrf-token":              {m.CSRFToken},
		"sec-ch-ua-mobile":          {"en"},
		"authorization":             {"Bearer " + m.BearerToken},
		"user-agent":                {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36"},
		"content-type":              {"application/json"},
		"cookie":                    {fmt.Sprintf("auth_token=%s; ct0=%s;", m.AuthToken, m.CSRFToken)},
		"referer":                   {"https://twitter.com/home"},
		"x-twitter-auth-type":       {"OAuth2Session"},
		"x-twitter-active-user":     {"yes"},
		"sec-ch-ua-platform":        {`"macOS"`},
	}

	res, err := m.Client.Do(req)

	if err != nil {
		return
	}

	if res.StatusCode != 200 {
		return
	}

	bBytes, _ := io.ReadAll(res.Body)
	defer res.Body.Close()


	var tweets FetchTweetsResponse
	if err := json.Unmarshal(bBytes, &tweets); err != nil {
		return
	}

	for indx, tweet := range tweets {
		parsedTime, _ := time.Parse(time.RubyDate, tweet.CreatedAt)
		unixTweetTime := parsedTime.Unix()
		if unixTweetTime >= filter.LatestTweetTS {
			parsedExtras := ParseExtras(tweet.FullText)
			m.PingChannel <- MonitorPing{
				Handle:  	tweet.User.ScreenName,
				Message: 	tweet.FullText,
				Image:   	tweet.User.ProfileImageURL,
				URL: 		fmt.Sprintf("https://twitter.com/%s/status/%s", tweet.User.ScreenName, tweet.IDStr),
				ParsedData: parsedExtras,
			}
		}
		if indx == 0 {
			filter.LatestTweetTS = unixTweetTime
		}

	}
}
