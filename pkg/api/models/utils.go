package models

type NullResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		BizCode int         `json:"biz_code"`
		BizMsg  string      `json:"biz_msg"`
		BizData interface{} `json:"biz_data"`
	} `json:"data"`
}

type PowChallenge struct {
	Algorithm   string `json:"algorithm"`
	Challenge   string `json:"challenge"`
	Salt        string `json:"salt"`
	Signature   string `json:"signature"`
	Difficulty  int    `json:"difficulty"`
	ExpireAt    int64  `json:"expire_at"`
	ExpireAfter int    `json:"expire_after"`
	TargetPath  string `json:"target_path"`
}

type PowResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		BizCode int    `json:"biz_code"`
		BizMsg  string `json:"biz_msg"`
		BizData struct {
			Challenge PowChallenge `json:"challenge"`
		} `json:"biz_data"`
	} `json:"data"`
}
