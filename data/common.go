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
	AppName          string `storm:"id" json:"name" graph:"name"`
	BinName          string `json:"filename" graph:"serviceFilename"`
	APIName          string `json:"api_prefix" graph:"apiName"`
	ServiceAddress   string `json:"address" graph:"serviceAddress"`
	IsManagedService bool   `json:"managed" graph:"isManaged"`
}
