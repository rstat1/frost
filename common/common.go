package common

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	// "github.com/evalphobia/logrus_sentry"

	"github.com/sirupsen/logrus"

	"crypto/rand"
	"net/http"
	"time"
)

//APIResponse ...
type APIResponse struct {
	Status         string `json:"status"`
	Response       string `json:"response"`
	HttpStatusCode int    `json:"-"`
}

//LocalAppConfig ...
type LocalAppConfig struct {
	BaseURL               string `json:"baseURL"`
	RedisServerAddr       string `json:"redisAddr"`
	VaultToken            string `json:"-"`
	VaultKeyID            string `json:"vaultKeyID"`
	VaultAddr             string `json:"vaultAddress"`
	VaultFrostAppRoleName string `json:"vaultAppRoleName"`
	VaultServiceAppRole   string `json:"vaultServicesAppRole"`
}

var (
	BaseURL       string
	CurrentConfig LocalAppConfig
	DevMode       bool
	Logger        *logrus.Logger
	httpServer    = &http.Server{
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
	}
)

//CreateAPIResponse ...
func CreateAPIResponse(response string, err error, failureCode int) APIResponse {
	if err == nil {
		return APIResponse{
			Status:         "success",
			Response:       response,
			HttpStatusCode: http.StatusOK,
		}
	} else {
		return APIResponse{
			Status:         "failed",
			Response:       err.Error(),
			HttpStatusCode: failureCode,
		}
	}
}

//WritePlainStringResponse ...
func WritePlainStringResponse(writer http.ResponseWriter, resp string, failCode int) {
	writeCommonHeaders(writer)
	writer.WriteHeader(failCode)
	writer.Write([]byte(resp))
}

//WriteAPIResponseStruct ...
func WriteAPIResponseStruct(writer http.ResponseWriter, resp APIResponse) {
	writeCommonHeaders(writer)
	writer.WriteHeader(resp.HttpStatusCode)
	apiResp, _ := json.Marshal(resp)
	writer.Write([]byte(apiResp))
}

//ValidatePOSTRequest ...
func ValidatePOSTRequest(validator func(*http.Request) APIResponse, handler func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(httpErrorHandler(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == "POST" && request.Header.Get("Content-Length") == "" {
			WriteAPIResponseStruct(writer, CreateAPIResponse("", errors.New("request body empty"), 400))
		} else {
			if resp := validator(request); resp.Status == "success" {
				handler(writer, request)
			} else {
				WriteAPIResponseStruct(writer, resp)
			}
		}
	}))
}

//RequestWrapper ...
func RequestWrapper(validator func(*http.Request) APIResponse, validMethod string, handler func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(httpErrorHandler(func(writer http.ResponseWriter, request *http.Request) {
		if validMethod != "" && request.Method != validMethod {
			WriteAPIResponseStruct(writer, APIResponse{
				Status:         "failed",
				Response:       "method not allowed",
				HttpStatusCode: http.StatusMethodNotAllowed,
			})
		} else {
			if resp := validator(request); resp.Status == "success" {
				handler(writer, request)
			} else {
				WriteAPIResponseStruct(writer, resp)
			}
		}
	}))
}
func httpErrorHandler(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
		// err, _ := SentryClient.CapturePanic(func() {
		// }, nil)
		// if err != nil {
		// 	WriteAPIResponseStruct(w, CreateAPIResponse("failed", errors.New("something serious happened."), 500))
		// }
	}
}

func writeCommonHeaders(writer http.ResponseWriter) {
	writer.Header().Add("Content-Type", "application/json")
	//writer.Header().Add("Access-Control-Allow-Origin", "http://192.168.1.12:4200")
	//writer.Header().Add("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
}

//Nothing ...
func Nothing(r *http.Request) APIResponse {
	return CreateAPIResponse("success", nil, 200)
}

// //SetupHTTPSListener ...
// func SetupHTTPSListener(handler http.Handler, port int) {
// 	m := autocert.Manager{
// 		Prompt:      autocert.AcceptTOS,
// 		Cache:       autocert.DirCache("certs"),
// 		RenewBefore: 5 * time.Hour,
// 		HostPolicy: autocert.HostWhitelist(
// 			"localhost",
// 			"gemini-svc.m.rdro.us",
// 			"gemini.rdro.us",
// 		),
// 		Email: "rstat1@gmail.com",
// 	}
// 	httpServer.TLSConfig = &tls.Config{
// 		GetCertificate: m.GetCertificate,
// 	}

// 	listener, err := reuseport.Listen("tcp", ":443")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer listener.Close()

// 	httpServer.Handler = handler
// 	httpServer.Addr = ":443"
// 	err = httpServer.ServeTLS(listener, "", "")

// 	if err != nil {
// 		Logger.WithField("func", "main").Errorln(err)
// 	}
// }

//InitLogrus ...
func InitLogrus() {
	Logger = logrus.New()
	Logger.Out = os.Stdout
	Logger.SetLevel(logrus.DebugLevel)
}

//CommonProcessInit ...
func CommonProcessInit(dev, loadConfig bool) {
	DevMode = dev
	InitLogrus()
	if os.Getenv("PWD") == "" {
		Logger.Warnln("pwd not set")
		os.Chdir("/webservices")
	}
	if loadConfig {
		if file, err := ioutil.ReadFile("config.json"); err == nil {
			err = json.Unmarshal([]byte(file), &CurrentConfig)
			if err != nil {
				panic(err)
			}
			BaseURL = "." + CurrentConfig.BaseURL
		} else {
			// common.Logger.Errorln()
			panic(errors.New("Please create a config.json file"))
		}
	}
}

//CreateFailureResponse ...
func CreateFailureResponse(err error, functionName string, status int) APIResponse {
	Logger.WithField("func", functionName).Errorln(err)
	return CreateAPIResponse("failed", err, status)
}

//CreateFailureResponseWithFields ...
func CreateFailureResponseWithFields(err error, status int, fields logrus.Fields) APIResponse {
	Logger.WithFields(fields).Errorln(err)
	return CreateAPIResponse("failed", err, status)
}

//WriteFailureResponse ..
func WriteFailureResponse(err error, resp http.ResponseWriter, functionName string, status int) {
	Logger.WithField("func", functionName).Errorln(err)
	WriteAPIResponseStruct(resp, CreateAPIResponse("failed", err, status))
}

//RandomID https://stackoverflow.com/questions/12771930
func RandomID(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

//Unzip https://golangcode.com/unzip-files-in-go/
func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			Logger.Errorln(err)
			return err
		}
		defer rc.Close()
		fpath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			var fdir string
			if lastIndex := strings.LastIndex(fpath, string(os.PathSeparator)); lastIndex > -1 {
				fdir = fpath[:lastIndex]
			}
			err = os.MkdirAll(fdir, os.ModePerm)
			if err != nil {
				Logger.Errorln(err)
				return err
			}
			f, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				Logger.Errorln(err)
				return err
			}
			defer f.Close()
			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

//LogError ...
func LogError(extra string, err error) error {
	if err != nil {
		pc, _, line, _ := runtime.Caller(1)
		funcObj := runtime.FuncForPC(pc)
		runtimeFunc := regexp.MustCompile(`^.*\.(.*)$`)
		name := runtimeFunc.ReplaceAllString(funcObj.Name(), "$1")

		if extra != "" {
			Logger.WithFields(logrus.Fields{"func": name, "line": line, "extra": extra}).Errorln(err)
		} else {
			Logger.WithFields(logrus.Fields{"func": name, "line": line}).Errorln(err)
		}
		return err
	}
	return nil
}

//LogDebug ...
func LogDebug(extraKey string, extraValue interface{}, entry interface{}) {
	pc, _, line, _ := runtime.Caller(1)
	funcObj := runtime.FuncForPC(pc)
	runtimeFunc := regexp.MustCompile(`^.*\.(.*)$`)
	name := runtimeFunc.ReplaceAllString(funcObj.Name(), "$1")

	if extraKey != "" {
		Logger.WithFields(logrus.Fields{"func": name, "line": line, extraKey: extraValue}).Debugln(entry)
	} else {
		Logger.WithFields(logrus.Fields{"func": name, "line": line}).Debugln(entry)
	}
}

//LogInfo ...
func LogInfo(extraKey string, extraValue interface{}, entry interface{}) {
	pc, _, line, _ := runtime.Caller(1)
	funcObj := runtime.FuncForPC(pc)
	runtimeFunc := regexp.MustCompile(`^.*\.(.*)$`)
	name := runtimeFunc.ReplaceAllString(funcObj.Name(), "$1")

	if extraKey != "" {
		Logger.WithFields(logrus.Fields{"func": name, "line": line, extraKey: extraValue}).Infoln(entry)
	} else {
		Logger.WithFields(logrus.Fields{"func": name, "line": line}).Infoln(entry)
	}
}

//LogWarn ...
func LogWarn(extraKey, extraValue string, entry interface{}) {
	pc, _, line, _ := runtime.Caller(1)
	funcObj := runtime.FuncForPC(pc)
	runtimeFunc := regexp.MustCompile(`^.*\.(.*)$`)
	name := runtimeFunc.ReplaceAllString(funcObj.Name(), "$1")

	if extraKey != "" {
		Logger.WithFields(logrus.Fields{"func": name, "line": line, extraKey: extraValue}).Warnln(entry)
	} else {
		Logger.WithFields(logrus.Fields{"func": name, "line": line}).Warnln(entry)
	}
}

//HasServiceCreds ...
func HasServiceCreds(r *http.Request) APIResponse {
	var resp APIResponse

	hasSid := hasRequiredParam("sid", r)
	hasSKey := hasRequiredParam("skey", r)
	if hasSid.Status == "success" && hasSKey.Status == "success" {
		resp = hasSid
	} else {
		resp = CreateFailureResponse(errors.New("missing required parameter"), "HasServiceCreds", 400)
	}
	return resp
}
func hasRequiredParam(param string, r *http.Request) APIResponse {
	value := r.Header.Get(param)
	if value == "" {
		return CreateFailureResponse(errors.New("missing required param"), "hasRequiredParam", 500)
	}
	return CreateAPIResponse(value, nil, 200)
}

//TimeTrack ...
//https://stackoverflow.com/questions/45766572
func TimeTrack(start time.Time) {
	elapsed := time.Since(start)
	pc, _, _, _ := runtime.Caller(1)
	funcObj := runtime.FuncForPC(pc)
	runtimeFunc := regexp.MustCompile(`^.*\.(.*)$`)
	name := runtimeFunc.ReplaceAllString(funcObj.Name(), "$1")
	Logger.WithFields(logrus.Fields{"elaspsed": elapsed, "func": name}).Debugln("done")
}

//NewFalse This only exists because the Vault client API is stupid.
func NewFalse() *bool {
	b := false
	return &b
}
