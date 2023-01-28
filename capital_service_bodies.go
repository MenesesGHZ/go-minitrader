package forexbot

type NewSessionBody struct {
	Identifier        string `json:"identifier"`
	Password          string `json:"password"`
	EncryptedPassword bool   `json:"encryptedPassword"`
}

type CreateWorkingOrderBody struct {
	Epic      string    `json:"epic"`
	Direction Signal    `json:"direction"`
	Type      OrderType `json:"type"`
	Size      float32   `json:"size"`
	Level     float32   `json:"level"`
}
