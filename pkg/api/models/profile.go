package models

type Profile struct {
	Id           string        `json:"id"`
	Token        string        `json:"token"`
	Email        string        `json:"email"`
	MobileNumber string        `json:"mobile_number"`
	AreaCode     string        `json:"area_code"`
	Status       int           `json:"status"`
	IdProfile    interface{}   `json:"id_profile"`
	IdProfiles   []interface{} `json:"id_profiles"`
	Chat         struct {
		IsMuted   int `json:"is_muted"`
		MuteUntil int `json:"mute_until"`
	} `json:"chat"`
	HasLegacyChatHistory bool `json:"has_legacy_chat_history"`
	NeedBirthday         bool `json:"need_birthday"`
}

type ProfileResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		BizCode int     `json:"biz_code"`
		BizMsg  string  `json:"biz_msg"`
		BizData Profile `json:"biz_data"`
	} `json:"data"`
}

type ThinkingQuota struct {
	Quota int `json:"quota"`
	Used  int `json:"used"`
}

type QuotaResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		BizCode int    `json:"biz_code"`
		BizMsg  string `json:"biz_msg"`
		BizData struct {
			Thinking ThinkingQuota `json:"thinking"`
		} `json:"biz_data"`
	} `json:"data"`
}
