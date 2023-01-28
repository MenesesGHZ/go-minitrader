package forexbot

type NewSessionBody struct {
	Identifier        string `json:"identifier"`
	Password          string `json:"password"`
	EncryptedPassword bool   `json:"encryptedPassword"`
}

type CreatePositionBody struct {
	Epic           string    `json:"epic"`
	Direction      Direction `json:"direction"`
	Size           string    `json:"size"`
	GuaranteedStop bool      `json:"guaranteedStop"`
	StopLevel      float64   `json:"stopLevel"`
	ProfitLevel    float64   `json:"profitLevel"`
}
