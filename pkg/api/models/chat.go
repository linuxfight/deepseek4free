package models

type ChatHistoryResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		BizCode int         `json:"biz_code"`
		BizMsg  string      `json:"biz_msg"`
		BizData ChatHistory `json:"biz_data"`
	} `json:"data"`
}

type ChatHistory struct {
	ChatSession  ChatSession   `json:"chat_session"`
	ChatMessages []ChatMessage `json:"chat_messages"`
	CacheValid   bool          `json:"cache_valid"`
	RouteId      interface{}   `json:"route_id"`
}

type ChatCreateResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		BizCode int         `json:"biz_code"`
		BizMsg  string      `json:"biz_msg"`
		BizData ChatSession `json:"biz_data"`
	} `json:"data"`
}

type ChatEditResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		BizCode int    `json:"biz_code"`
		BizMsg  string `json:"biz_msg"`
		BizData struct {
			ChatSessionUpdatedAt float64 `json:"chat_session_updated_at"`
		} `json:"biz_data"`
	} `json:"data"`
}

type ChatSession struct {
	Id               string  `json:"id"`
	SeqId            int     `json:"seq_id"`
	Title            *string `json:"title"`
	TitleType        *string `json:"title_type"`
	UpdatedAt        float64 `json:"updated_at"`
	Agent            string  `json:"agent"`
	Version          int32   `json:"version"`
	CurrentMessageId *int32  `json:"current_message_id"`
	InsertedAt       float64 `json:"inserted_at"`
	Character        *string `json:"character"`
}

type ChatMessage struct {
	MessageId             int              `json:"message_id"`
	ParentId              *int             `json:"parent_id"`
	Model                 string           `json:"model"`
	Role                  string           `json:"role"`
	Content               string           `json:"content"`
	ThinkingEnabled       bool             `json:"thinking_enabled"`
	ThinkingContent       *string          `json:"thinking_content"`
	ThinkingElapsedSecs   *int             `json:"thinking_elapsed_secs"`
	BanEdit               bool             `json:"ban_edit"`
	BanRegenerate         bool             `json:"ban_regenerate"`
	Status                string           `json:"status"`
	AccumulatedTokenUsage int              `json:"accumulated_token_usage"`
	Files                 []interface{}    `json:"files"`
	Tips                  []interface{}    `json:"tips"`
	InsertedAt            float64          `json:"inserted_at"`
	SearchEnabled         bool             `json:"search_enabled"`
	SearchStatus          *string          `json:"search_status"`
	SearchResults         []SearchResponse `json:"search_results"`
}

type SearchResponse struct {
	Url         string      `json:"url"`
	Title       string      `json:"title"`
	Snippet     string      `json:"snippet"`
	CiteIndex   *int        `json:"cite_index"`
	PublishedAt interface{} `json:"published_at"`
	SiteName    interface{} `json:"site_name"`
	SiteIcon    string      `json:"site_icon"`
}

type ChatListResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		BizCode int    `json:"biz_code"`
		BizMsg  string `json:"biz_msg"`
		BizData struct {
			ChatSessions []ChatSession `json:"chat_sessions"`
			HasMore      bool          `json:"has_more"`
		} `json:"biz_data"`
	} `json:"data"`
}
