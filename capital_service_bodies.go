package forexbot

type NewSessionBody struct {
	Identifier        string `json:"identifier"`
	Password          string `json:"password"`
	EncryptedPassword bool   `json:"encryptedPassword"`
}
