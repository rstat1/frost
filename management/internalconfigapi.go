package management

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"git.m/svcman/common"
	"git.m/svcman/crypto"
	"git.m/svcman/data"
	"github.com/husobee/vestigo"
)

//InternalConfigAPI ...
type InternalConfigAPI struct {
	router *vestigo.Router
	data   *data.DataStore
	vault  *crypto.VaultClient
}

//NewInternalConfigAPI ...
func NewInternalConfigAPI(ds *data.DataStore, vc *crypto.VaultClient) *InternalConfigAPI {
	return &InternalConfigAPI{
		vault:  vc,
		router: vestigo.NewRouter(),
	}
}

//InitListener ...
func (icapi *InternalConfigAPI) InitListener() {
	go func() {
		icapi.setAPIRoutes()
		if err := http.ListenAndServe("localhost:5000", icapi.router); err != nil {
			common.LogError("", err)
		}
	}()
}

func (icapi *InternalConfigAPI) setAPIRoutes() {
	icapi.router.Handle("/api/icapi/config/set", common.RequestWrapper(common.Nothing, "POST", icapi.setConfigValue))
	icapi.router.Handle("/api/icapi/config/get/:key", common.RequestWrapper(common.Nothing, "GET", icapi.getConfigValue))
}

func (icapi *InternalConfigAPI) getConfigValue(resp http.ResponseWriter, r *http.Request) {
	var value string

	var valueToGet string = vestigo.Param(r, "key")

	switch valueToGet {
	case "dbAddr":
		value = common.CurrentConfig.DBAddr
		break
	case "dbName":
		value = common.CurrentConfig.DBName
		break
	case "vaultAddr":
		value = common.CurrentConfig.VaultAddr
		break
	case "dbCreds":
		if user, pw, e := icapi.vault.GetDBCredentials(); e == nil {
			creds := map[string]string{
				"username": user,
				"password": pw,
			}
			msgBytes, _ := json.Marshal(creds)
			value = string(msgBytes)
		} else {
			common.WriteFailureResponse(e, resp, "getConfigValue", 400)
			return
		}
		break
	default:
		common.WriteFailureResponse(errors.New("invalid config key specified"), resp, "setConfigValue", 400)
	}

	common.WriteAPIResponseStruct(resp, common.CreateAPIResponse(value, nil, 400))

}

func (icapi *InternalConfigAPI) setConfigValue(resp http.ResponseWriter, r *http.Request) {
	var changeMade bool
	var change data.ConfigChangeRequest
	if body, err := ioutil.ReadAll(r.Body); err == nil {
		json.Unmarshal(body, &change)
		switch change.Key {
		case "dbAddr":
			common.CurrentConfig.DBAddr = change.Value
			changeMade = true
			break
		case "dbName":
			common.CurrentConfig.DBName = change.Value
			changeMade = true
			break
		case "vaultToken":
			common.CurrentConfig.VaultToken = change.Value
			icapi.vault.SetAccessToken()
			changeMade = true
			break
		case "vaultAddr":
			common.CurrentConfig.VaultAddr = change.Value
			changeMade = true
			break
		case "vaultARID":
			common.CurrentConfig.VaultRoleID = change.Value
			break
		default:
			common.WriteFailureResponse(errors.New("invalid config key specified"), resp, "setConfigValue", 400)
			changeMade = false
		}
		if changeMade {
			config, _ := json.Marshal(common.CurrentConfig)
			common.LogError("", ioutil.WriteFile("config.json", config, 0600))
		}
		common.WriteAPIResponseStruct(resp, common.CreateAPIResponse("success", nil, 400))
	} else {
		common.WriteFailureResponse(err, resp, "setConfigValue", 500)
	}
}
