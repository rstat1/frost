package data

//User ...
type User struct {
	Id, PassHash string
	Username     string `storm:"id"`
	Group        string
}

//AuthRequest ...
type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//UserDetails ...
type UserDetails struct {
	Username    string           `json:"username"`
	Password    string           `json:"password"`
	Permissions []ServiceAuth `json:"access"`
}

//PermissionList ...
type PermissionList struct {
	Permissions []ServiceAuth `json:"permissions"`
}

//ServiceAuth ...
type ServiceAuth struct {
	Service     string            `json:"service"`
	Permissions []PermissionValue `json:"permissions"`
}

//PermissionValue ...
type PermissionValue struct {
	Name  string `json:"name"`
	Value bool   `json:"value"`
}

//Config ...
type Config struct {
	ID                 string `storm:"id"`
	UpdateAuthKey      string
	ActivationKey      string
	ActivationRequired bool `json:"ActivationRequired"`
}

//ServiceDetails ...
type ServiceDetails struct {
	AppName          string `storm:"id" json:"name" graph:"name"`
	BinName          string `json:"filename" graph:"serviceFilename"`
	APIName          string `json:"api_prefix" graph:"apiName"`
	ServiceID        string `storm:"unique"`
	ServiceKey       string
	RedirectURL      string
	ServiceAddress   string `json:"address" graph:"serviceAddress"`
	IsManagedService bool   `json:"managed" graph:"isManaged"`
	ServiceNameURLToUI bool `json:"serviceNameToUI"`
}

//ServiceAccess ...
type ServiceAccess struct {
	Service    string `json:"service"`
	Permission []struct {
		Name  string `json:"name"`
		Value bool   `json:"value"`
	}
}

//BingDailyImage ...
type BingDailyImage struct {
	Images []struct {
		Startdate     string        `json:"startdate"`
		Fullstartdate string        `json:"fullstartdate"`
		Enddate       string        `json:"enddate"`
		URL           string        `json:"url"`
		Urlbase       string        `json:"urlbase"`
		Copyright     string        `json:"copyright"`
		Copyrightlink string        `json:"copyrightlink"`
		Quiz          string        `json:"quiz"`
		Wp            bool          `json:"wp"`
		Hsh           string        `json:"hsh"`
		Drk           int           `json:"drk"`
		Top           int           `json:"top"`
		Bot           int           `json:"bot"`
		Hs            []interface{} `json:"hs"`
	} `json:"images"`
	Tooltips struct {
		Loading  string `json:"loading"`
		Previous string `json:"previous"`
		Next     string `json:"next"`
		Walle    string `json:"walle"`
		Walls    string `json:"walls"`
	} `json:"tooltips"`
}
