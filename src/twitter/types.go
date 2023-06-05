package twitter


type TwitterDataResponse struct {
	Data struct {
		User struct {
			Result struct {
				Typename                   string `json:"__typename"`
				Id                         string `json:"id"`
				RestId                     string `json:"rest_id"`
				AffiliatesHighlightedLabel struct {
				} `json:"affiliates_highlighted_label"`
				HasNftAvatar   bool `json:"has_nft_avatar"`
				IsBlueVerified bool `json:"is_blue_verified"`
				Legacy         struct {
					BlockedBy           bool   `json:"blocked_by"`
					Blocking            bool   `json:"blocking"`
					CanDm               bool   `json:"can_dm"`
					CanMediaTag         bool   `json:"can_media_tag"`
					CreatedAt           string `json:"created_at"`
					DefaultProfile      bool   `json:"default_profile"`
					DefaultProfileImage bool   `json:"default_profile_image"`
					Description         string `json:"description"`
					Entities            struct {
						Description struct {
							Urls []interface{} `json:"urls"`
						} `json:"description"`
					} `json:"entities"`
					FastFollowersCount      int           `json:"fast_followers_count"`
					FavouritesCount         int           `json:"favourites_count"`
					FollowRequestSent       bool          `json:"follow_request_sent"`
					FollowedBy              bool          `json:"followed_by"`
					FollowersCount          int           `json:"followers_count"`
					Following               bool          `json:"following"`
					FriendsCount            int           `json:"friends_count"`
					HasCustomTimelines      bool          `json:"has_custom_timelines"`
					IsTranslator            bool          `json:"is_translator"`
					ListedCount             int           `json:"listed_count"`
					Location                string        `json:"location"`
					MediaCount              int           `json:"media_count"`
					Muting                  bool          `json:"muting"`
					Name                    string        `json:"name"`
					NeedsPhoneVerification  bool          `json:"needs_phone_verification"`
					NormalFollowersCount    int           `json:"normal_followers_count"`
					Notifications           bool          `json:"notifications"`
					PinnedTweetIdsStr       []string      `json:"pinned_tweet_ids_str"`
					PossiblySensitive       bool          `json:"possibly_sensitive"`
					ProfileBannerUrl        string        `json:"profile_banner_url"`
					ProfileImageUrlHttps    string        `json:"profile_image_url_https"`
					ProfileInterstitialType string        `json:"profile_interstitial_type"`
					Protected               bool          `json:"protected"`
					ScreenName              string        `json:"screen_name"`
					StatusesCount           int           `json:"statuses_count"`
					TranslatorType          string        `json:"translator_type"`
					Verified                bool          `json:"verified"`
					WantRetweets            bool          `json:"want_retweets"`
					WithheldInCountries     []interface{} `json:"withheld_in_countries"`
				} `json:"legacy"`
				SmartBlockedBy        bool `json:"smart_blocked_by"`
				SmartBlocking         bool `json:"smart_blocking"`
				SuperFollowEligible   bool `json:"super_follow_eligible"`
				SuperFollowedBy       bool `json:"super_followed_by"`
				SuperFollowing        bool `json:"super_following"`
				LegacyExtendedProfile struct {
					Birthdate struct {
						Day            int    `json:"day"`
						Month          int    `json:"month"`
						Year           int    `json:"year"`
						Visibility     string `json:"visibility"`
						YearVisibility string `json:"year_visibility"`
					} `json:"birthdate"`
				} `json:"legacy_extended_profile"`
				IsProfileTranslatable bool `json:"is_profile_translatable"`
			} `json:"result"`
		} `json:"user"`
	} `json:"data"`
}
type FetchTweetsResponse []struct {
	CreatedAt        string `json:"created_at"`
	ID               int64  `json:"id"`
	IDStr            string `json:"id_str"`
	FullText         string `json:"full_text"`
	Truncated        bool   `json:"truncated"`
	DisplayTextRange []int  `json:"display_text_range"`
	Entities         struct {
		Hashtags     []any `json:"hashtags"`
		Symbols      []any `json:"symbols"`
		UserMentions []struct {
			ScreenName string `json:"screen_name"`
			Name       string `json:"name"`
			ID         int64  `json:"id"`
			IDStr      string `json:"id_str"`
			Indices    []int  `json:"indices"`
		} `json:"user_mentions"`
		Urls []any `json:"urls"`
		Media []struct {
			ID            int64  `json:"id,omitempty"`
			IDStr         string `json:"id_str,omitempty"`
			Indices       []int  `json:"indices,omitempty"`
			MediaURL      string `json:"media_url,omitempty"`
			MediaURLHTTPS string `json:"media_url_https,omitempty"`
			URL           string `json:"url,omitempty"`
			DisplayURL    string `json:"display_url,omitempty"`
			ExpandedURL   string `json:"expanded_url,omitempty"`
			Type          string `json:"type,omitempty"`
			OriginalInfo  any `json:"original_info,omitempty"`
			Sizes any `json:"sizes,omitempty"`
			Features any `json:"features,omitempty"`
		} `json:"media,omitempty"`
	} `json:"entities"`
	Source               string `json:"source"`
	InReplyToStatusID    int64  `json:"in_reply_to_status_id"`
	InReplyToStatusIDStr string `json:"in_reply_to_status_id_str"`
	InReplyToUserID      int64  `json:"in_reply_to_user_id"`
	InReplyToUserIDStr   string `json:"in_reply_to_user_id_str"`
	InReplyToScreenName  string `json:"in_reply_to_screen_name"`
	User                 struct {
		ID          int64  `json:"id"`
		IDStr       string `json:"id_str"`
		Name        string `json:"name"`
		ScreenName  string `json:"screen_name"`
		Location    string `json:"location"`
		Description string `json:"description"`
		URL         any    `json:"url"`
		Entities    struct {
			Description struct {
				Urls []any `json:"urls"`
			} `json:"description"`
		} `json:"entities"`
		Protected                      bool     `json:"protected"`
		FollowersCount                 int      `json:"followers_count"`
		FastFollowersCount             int      `json:"fast_followers_count"`
		NormalFollowersCount           int      `json:"normal_followers_count"`
		FriendsCount                   int      `json:"friends_count"`
		ListedCount                    int      `json:"listed_count"`
		CreatedAt                      string   `json:"created_at"`
		FavouritesCount                int      `json:"favourites_count"`
		UtcOffset                      any      `json:"utc_offset"`
		TimeZone                       any      `json:"time_zone"`
		GeoEnabled                     bool     `json:"geo_enabled"`
		Verified                       bool     `json:"verified"`
		StatusesCount                  int      `json:"statuses_count"`
		MediaCount                     int      `json:"media_count"`
		Lang                           any      `json:"lang"`
		ContributorsEnabled            bool     `json:"contributors_enabled"`
		IsTranslator                   bool     `json:"is_translator"`
		IsTranslationEnabled           bool     `json:"is_translation_enabled"`
		ProfileBackgroundColor         string   `json:"profile_background_color"`
		ProfileBackgroundImageURL      any      `json:"profile_background_image_url"`
		ProfileBackgroundImageURLHTTPS any      `json:"profile_background_image_url_https"`
		ProfileBackgroundTile          bool     `json:"profile_background_tile"`
		ProfileImageURL                string   `json:"profile_image_url"`
		ProfileImageURLHTTPS           string   `json:"profile_image_url_https"`
		ProfileBannerURL               string   `json:"profile_banner_url"`
		ProfileLinkColor               string   `json:"profile_link_color"`
		ProfileSidebarBorderColor      string   `json:"profile_sidebar_border_color"`
		ProfileSidebarFillColor        string   `json:"profile_sidebar_fill_color"`
		ProfileTextColor               string   `json:"profile_text_color"`
		ProfileUseBackgroundImage      bool     `json:"profile_use_background_image"`
		HasExtendedProfile             bool     `json:"has_extended_profile"`
		DefaultProfile                 bool     `json:"default_profile"`
		DefaultProfileImage            bool     `json:"default_profile_image"`
		PinnedTweetIds                 []int64  `json:"pinned_tweet_ids"`
		PinnedTweetIdsStr              []string `json:"pinned_tweet_ids_str"`
		HasCustomTimelines             bool     `json:"has_custom_timelines"`
		CanMediaTag                    bool     `json:"can_media_tag"`
		FollowedBy                     bool     `json:"followed_by"`
		Following                      bool     `json:"following"`
		FollowRequestSent              bool     `json:"follow_request_sent"`
		Notifications                  bool     `json:"notifications"`
		AdvertiserAccountType          string   `json:"advertiser_account_type"`
		AdvertiserAccountServiceLevels []any    `json:"advertiser_account_service_levels"`
		BusinessProfileState           string   `json:"business_profile_state"`
		TranslatorType                 string   `json:"translator_type"`
		WithheldInCountries            []any    `json:"withheld_in_countries"`
		RequireSomeConsent             bool     `json:"require_some_consent"`
	} `json:"user"`
	Geo                  any    `json:"geo"`
	Coordinates          any    `json:"coordinates"`
	Place                any    `json:"place"`
	Contributors         any    `json:"contributors"`
	IsQuoteStatus        bool   `json:"is_quote_status"`
	RetweetCount         int    `json:"retweet_count"`
	FavoriteCount        int    `json:"favorite_count"`
	ConversationID       int64  `json:"conversation_id"`
	ConversationIDStr    string `json:"conversation_id_str"`
	Favorited            bool   `json:"favorited"`
	Retweeted            bool   `json:"retweeted"`
	Lang                 string `json:"lang"`
	SupplementalLanguage any    `json:"supplemental_language"`
}