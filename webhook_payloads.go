package destinyhome

type webhookRequest struct {
	Device struct {
		Capabilities []string `json:"capabilities"`
	} `json:"device"`
	Handler struct {
		Name string `json:"name"`
	} `json:"handler"`
	Home struct {
		Params struct{} `json:"params"`
	} `json:"home"`
	Intent struct {
		Name   string           `json:"name"`
		Params map[string]param `json:"params"`
		Query  string           `json:"query"`
	} `json:"intent"`
	Scene struct {
		Name              string   `json:"name"`
		SlotFillingStatus string   `json:"slotFillingStatus"`
		Slots             struct{} `json:"slots"`
	} `json:"scene"`
	Session struct {
		ID            string        `json:"id"`
		Params        struct{}      `json:"params"`
		TypeOverrides []interface{} `json:"typeOverrides"`
	} `json:"session"`
	User struct {
		Locale string `json:"locale"`
		Params struct {
			VerificationStatus string `json:"verificationStatus"`
		} `json:"params"`
	} `json:"user"`
}

type webhookResponse struct {
	Prompt struct {
		FirstSimple struct {
			Speech string `json:"speech"`
			Text   string `json:"text"`
		} `json:"firstSimple"`
		Override bool `json:"override"`
	} `json:"prompt"`
	Scene struct {
		Name string `json:"name"`
		Next struct {
			Name string `json:"name"`
		} `json:"next"`
		Slots struct{} `json:"slots"`
	} `json:"scene"`
	Session struct {
		ID     string   `json:"id"`
		Params struct{} `json:"params"`
	} `json:"session"`
}

type param struct {
	Original string `json:"original"`
	Resolved string `json:"resolved"`
}
