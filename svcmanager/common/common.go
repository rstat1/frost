package common

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"os"

	"github.com/evalphobia/logrus_sentry"
	"github.com/getsentry/raven-go"
	"github.com/kavu/go_reuseport"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"

	"crypto/rand"
	"net/http"
	"time"
)

type APIResponse struct {
	Status         string `json:"status"`
	Response       string `json:"response"`
	HttpStatusCode int    `json:"-"`
}

var (
	SentryClient *raven.Client
	Logger       *logrus.Logger
	httpServer   = &http.Server{
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
	}
)

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
func WritePlainStringResponse(writer http.ResponseWriter, resp string, failCode int) {
	writeCommonHeaders(writer)
	writer.WriteHeader(failCode)
	apiResp, _ := json.Marshal(resp)
	writer.Write([]byte(apiResp))
}
func WriteAPIResponseStruct(writer http.ResponseWriter, resp APIResponse) {
	writeCommonHeaders(writer)
	writer.WriteHeader(resp.HttpStatusCode)
	apiResp, _ := json.Marshal(resp)
	writer.Write([]byte(apiResp))
}
func ValidateRequest(validator func(*http.Request) APIResponse, handler func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(raven.RecoveryHandler(func(writer http.ResponseWriter, request *http.Request) {
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
func RequestWrapper(validator func(*http.Request) APIResponse, validMethod string, handler func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(raven.RecoveryHandler(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != validMethod {
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
func writeCommonHeaders(writer http.ResponseWriter) {
	writer.Header().Add("Content-Type", "application/json")
	//writer.Header().Add("Access-Control-Allow-Origin", "http://192.168.1.12:4200")
	//writer.Header().Add("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
}
func Nothing(r *http.Request) APIResponse {
	return CreateAPIResponse("success", nil, 200)
}
func LogDebug(functionName, extra string, details interface{}) {
	Logger.WithFields(logrus.Fields{
		"func":  functionName,
		"extra": extra,
	}).Debugln(details)
}
func LogWarn(functionName, extra string, details interface{}) {
	Logger.WithFields(logrus.Fields{
		"func":  functionName,
		"extra": extra,
	}).Warnln(details)
}
func LogError(functionName, extra string, details interface{}) {
	Logger.WithFields(logrus.Fields{
		"func":  functionName,
		"extra": extra,
	}).Errorln(details)
}
func SetupHTTPSListener(handler http.Handler, port int) {
	m := autocert.Manager{
		Prompt:      autocert.AcceptTOS,
		Cache:       autocert.DirCache("certs"),
		RenewBefore: 5 * time.Hour,
		HostPolicy: autocert.HostWhitelist(
			"localhost",
			"gemini-svc.m.rdro.us",
			"gemini.rdro.us",
		),
		Email: "rstat1@gmail.com",
	}
	httpServer.TLSConfig = &tls.Config{
		GetCertificate: m.GetCertificate,
	}

	listener, err := reuseport.Listen("tcp", ":443")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	httpServer.Handler = handler
	httpServer.Addr = ":443"
	err = httpServer.ServeTLS(listener, "", "")

	if err != nil {
		Logger.WithField("func", "main").Errorln(err)
	}
}
func InitLogrus() {
	Logger = logrus.New()
	Logger.Out = os.Stdout
	Logger.SetLevel(logrus.DebugLevel)

	client, err := raven.New("http://57ad78f1ed984ff2bcb5a6d40760431d:33389d1773c049c1abc69dd246b6ae2a@sentry.m/8")

	if err != nil {
		Logger.Fatal(err)
	} else {
		SentryClient = client
	}

	hook, err := logrus_sentry.NewWithClientSentryHook(client, []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
	})

	if err == nil {
		Logger.Hooks.Add(hook)
	}
}
func CommonProcessInit() {
	InitLogrus()
	if os.Getenv("PWD") == "" {
		Logger.Warnln("pwd not set")
		os.Chdir("/gemini")
	}
}
func CreateFailureResponse(err error, functionName string, status int) APIResponse {
	Logger.WithField("func", functionName).Errorln(err)
	return CreateAPIResponse("failed", err, status)
}
func CreateFailureResponseWithFields(err error, status int, fields logrus.Fields) APIResponse {
	Logger.WithFields(fields).Errorln(err)
	return CreateAPIResponse("failed", err, status)
}
func WriteFailureResponse(err error, resp http.ResponseWriter, functionName string, status int) {
	Logger.WithField("func", functionName).Errorln(err)
	WriteAPIResponseStruct(resp, CreateAPIResponse("failed", err, status))
}

//https://stackoverflow.com/questions/12771930

func RandomID(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}
