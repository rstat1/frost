package management

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/husobee/vestigo"
	"github.com/minio/sio"
	"go.alargerobot.dev/frost/common"
	"go.alargerobot.dev/frost/crypto"
	"go.alargerobot.dev/frost/data"
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
		data:   ds,
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

//ValidateServicesCreds ...
func (icapi *InternalConfigAPI) ValidateServicesCreds(r *http.Request) common.APIResponse {
	serviceName := vestigo.Param(r, "service")
	sid := r.Header.Get("X-Frost-ServiceID")
	skey := r.Header.Get("X-Frost-ServiceKey")
	if sid != "" && skey != "" {
		if details, err := icapi.data.GetRoute(serviceName); err == nil {
			if details.ServiceID == sid && details.ServiceKey == skey {
				return common.CreateAPIResponse("success", nil, 200)
			} else {
				return common.CreateAPIResponse("failed", errors.New("wrong sid or skey"), 400)
			}
		} else {
			return common.CreateAPIResponse("failed", err, 400)
		}
	}
	return common.CreateAPIResponse("failed", errors.New("missing required header(s)"), 400)
}

func (icapi *InternalConfigAPI) setAPIRoutes() {
	icapi.router.Handle("/api/icapi/config/:service/set/:key", common.ValidatePOSTRequest(icapi.ValidateServicesCreds, icapi.setConfigValue))
	icapi.router.Handle("/api/icapi/config/:service/get/:key", common.RequestWrapper(icapi.ValidateServicesCreds, "GET", icapi.getConfigValue))
}

func (icapi *InternalConfigAPI) getConfigValue(resp http.ResponseWriter, r *http.Request) {
	var entryKey Key
	var entryCryptoKey ConfigEncryptionKey
	var key string = vestigo.Param(r, "key")
	var serviceName string = vestigo.Param(r, "service")

	if serviceName == "watchdog" {
		common.WriteAPIResponseStruct(resp, common.CreateAPIResponse("failed", errors.New("no"), 400))
		return
	}

	if v, e := icapi.data.GetServiceConfigValue(key, serviceName); e == nil {
		valueStr, _ := base64.StdEncoding.DecodeString(v)
		value := bytes.NewBuffer([]byte(valueStr))
		if encipheredValue, err := icapi.vault.ReadKeyFromKV("service-config/" + serviceName + "/" + key); err == nil {
			data, _ := base64.StdEncoding.DecodeString(string(encipheredValue))
			common.LogError("", json.Unmarshal(data, &entryCryptoKey))
			if masterKey, err := icapi.vault.UnsealKey(crypto.FrostKeyID, entryCryptoKey.SealedMasterKey, crypto.Context{"key": key}); err == nil {
				entryKey.Unseal(masterKey[:], entryCryptoKey.EntryKey)
				decipheredRead, err := sio.DecryptReader(value, sio.Config{Key: entryKey[:], MinVersion: sio.Version20})
				if err != nil {
					common.WriteFailureResponse(err, resp, "setConfigValue", 500)
					return
				}
				encipheredValue, err := ioutil.ReadAll(decipheredRead)
				common.WriteAPIResponseStruct(resp, common.CreateAPIResponse(string(encipheredValue), nil, 200))
			} else {
				common.WriteFailureResponse(err, resp, "setConfigValue", 500)
			}
		} else {
			common.WriteFailureResponse(err, resp, "setConfigValue", 500)
		}
	} else {
		common.WriteAPIResponseStruct(resp, common.CreateAPIResponse("failed", errors.New("invalid key specified"), 400))
	}

}

func (icapi *InternalConfigAPI) setConfigValue(resp http.ResponseWriter, r *http.Request) {
	var value string

	key := vestigo.Param(r, "key")

	serviceName := vestigo.Param(r, "service")

	if body, err := ioutil.ReadAll(r.Body); err == nil {
		value = string(body)
	} else {
		common.WriteFailureResponse(err, resp, "setConfigValue", 500)
		return
	}
	if serviceName == "watchdog" {
		common.WriteAPIResponseStruct(resp, icapi.setFrostConfigValue(key, value))
	} else {
		if vKey, vSealed, err := icapi.vault.GenerateKey(crypto.FrostKeyID, crypto.Context{"key": key}); err == nil {
			cryptoKey := GenerateKey(vKey[:], "service-config/"+serviceName+"/"+key)
			sealed, _ := cryptoKey.Seal(vKey[:], "service-config/"+serviceName+"/"+key)
			value := bytes.NewBuffer([]byte(value))
			if encipheredRead, err := sio.EncryptReader(value, sio.Config{Key: cryptoKey[:], MinVersion: sio.Version20}); err == nil {
				encipheredValue, err := ioutil.ReadAll(encipheredRead)
				if err != nil {
					common.WriteFailureResponse(err, resp, "setConfigValue", 500)
					return
				}
				entryCryptoKey := ConfigEncryptionKey{EntryKey: sealed, SealedMasterKey: vSealed}
				ecKey, _ := json.Marshal(entryCryptoKey)
				if e := icapi.vault.WriteKeyToKVStorage(ecKey, "service-config/"+serviceName+"/"+key); e != nil {
					common.WriteFailureResponse(e, resp, "setConfigValue", 500)
					return
				}
				common.WriteAPIResponseStruct(resp, common.CreateAPIResponse("success", icapi.data.SetConfigValue(key, serviceName, encipheredValue), 200))
			}
		} else {
			common.WriteFailureResponse(err, resp, "setConfigValue", 500)
		}
	}
}

func (icapi *InternalConfigAPI) setFrostConfigValue(key, value string) (resp common.APIResponse) {
	var changeMade bool

	switch key {
	case "dbAddr":
		common.CurrentConfig.DBAddr = value
		changeMade = true
		break
	case "dbName":
		common.CurrentConfig.DBName = value
		changeMade = true
		break
	case "vaultToken":
		common.CurrentConfig.VaultToken = value
		icapi.vault.SetAccessToken()
		changeMade = true
		break
	case "vaultAddr":
		common.CurrentConfig.VaultAddr = value
		changeMade = true
		break
	case "vaultARID":
		common.CurrentConfig.VaultRoleID = value
		break
	default:
		resp = common.CreateAPIResponse("success", errors.New("invalid config key specified"), 400)
		changeMade = false
	}
	if changeMade {
		config, _ := json.Marshal(common.CurrentConfig)
		common.LogError("", ioutil.WriteFile("config.json", config, 0600))
	}
	resp = common.CreateAPIResponse("success", nil, 400)
	return resp
}
