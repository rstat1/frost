package data

type User struct {
	Id, PassHash string
	Username     string `storm:"id"`
	Group        string
}
type ClientDetails struct {
	ID           string `storm:"id"`
	ClientID     string
	ClientSecret string
}
type Config struct {
	ID                 string `storm:"id"`
	UpdateAuthKey      string
	ActivationKey      string
	ActivationRequired bool `json:"ActivationRequired"`
}
type KnownRoute struct {
	AppName          string `storm:"id" json:"name"`
	BinName          string `json:"filename"`
	APIName          string `json:"api_prefix"`
	ServiceAddress   string `json:"address"`
	IsManagedService bool   `json:"managed"`
}
