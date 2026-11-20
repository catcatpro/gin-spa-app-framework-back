package requests

type SystemSettingsRequest struct {
	Id       string `json:"id"`
	SetKey   string `json:"set_key"`
	SetValue string `json:"set_value"`
}
