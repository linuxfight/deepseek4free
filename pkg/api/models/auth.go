package models

type AuthResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		BizCode int    `json:"biz_code"`
		BizMsg  string `json:"biz_msg"`
		BizData struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
			User struct {
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
				} `json:"Completion"`
				HasLegacyChatHistory bool `json:"has_legacy_chat_history"`
				NeedBirthday         bool `json:"need_birthday"`
			} `json:"user"`
		} `json:"biz_data"`
	} `json:"data"`
}
