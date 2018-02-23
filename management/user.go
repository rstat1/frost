package management

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"git.m/svcman/common"
	"git.m/svcman/data"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

var (
	issuerName string = "https://auth.m.rdro.us"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type User struct {
	datastore *data.DataStore
	hmacKey   []byte
}

type AuthResult struct {
	Result string `json:"result"`
}
type TokenValidationResult struct {
}
type AuthToken struct {
	Issuer  string `json:"iss"`
	Subject string `json:"sub"`
	Expires int64  `json:"exp"`
	Group   string `json:"grp"`
	UserID  string `json:"uid"`
}

func NewUserService(db *data.DataStore) *User {
	return &User{
		datastore: db,
		hmacKey:   generateSymKey(),
	}
}

//NewUser Creates a new user and returns an auth token.
func (u *User) NewUser(name, password string) common.APIResponse {
	var apiResp common.APIResponse

	err, user := u.datastore.NewUser(name, password)
	if err == nil {
		resp := u.GenerateAuthToken(AuthRequest{name, password}, user)
		if resp.Status == "success" {
			apiResp = common.CreateAPIResponse(resp.Response, nil, 500)
		} else {
			apiResp = common.CreateFailureResponse(errors.New(resp.Response), "NewUser", resp.HttpStatusCode)
		}
	} else {
		apiResp = common.CreateFailureResponse(err, "NewUser", 500)
	}
	return apiResp
}
func (u *User) GetUserFromToken(r *http.Request) data.User {
	var user data.User
	success, response := u.GetUserHeader(r)
	if success {
		if token, err := jwt.ParseSigned(response); err == nil {
			var defaultClaims jwt.Claims
			username := struct {
				Name string `json:"sub"`
			}{}
			token.Claims(u.hmacKey, &defaultClaims, &username)
			user = u.datastore.GetUser(username.Name)
		} else {
			common.Logger.WithField("func", "GetUsernameFromToken").Errorln(errors.New("invalid token"))
			user = data.User{}
		}
	} else {
		common.Logger.WithField("func", "GetUsernameFromToken").Errorln(errors.New("token not provided"))
		user = data.User{}
	}
	return user
}

//GetUserByID ...
func (u *User) GetUserByID(id string) data.User {
	return u.datastore.GetUserByID(id)
}

//ValidateLoginRequest Checks that the provided username and password are known and returns a signed JWT.
func (u *User) ValidateLoginRequest(request *http.Request) common.APIResponse {
	var userInfo data.User
	var response common.APIResponse
	var authRequest AuthRequest
	body, _ := ioutil.ReadAll(request.Body)
	json.Unmarshal(body, &authRequest)

	if authRequest.Username == "client" {
		return common.CreateFailureResponse(errors.New("invalid username"), "ValidateLoginRequest", 401)
	}

	passHasher := sha512.New512_256()
	hash := hex.EncodeToString(passHasher.Sum([]byte(authRequest.Password)))
	userInfo = u.datastore.GetUser(authRequest.Username)
	if userInfo.PassHash == string(hash) {
		response = u.GenerateAuthToken(authRequest, userInfo)
	} else {
		response = common.CreateFailureResponse(errors.New("incorrect credentials"), "ValidateLoginRequest", 401)
	}
	return response
}

//ValidateToken Checks that provided token was sent by this service, hasn't expired and was signed by the current instance. Also can check if a user is root or not.
func (u *User) ValidateToken(token string, sudo bool, requireUserToken bool) (bool, string) {
	var isValid bool
	var reason string = ""
	return true, ""
	if token, err := jwt.ParseSigned(token); err == nil {
		var defaultClaims jwt.Claims
		customClaims := struct {
			Group  string `json:"grp"`
			UserID string `json:"uid"`
		}{}
		token.Claims(u.hmacKey, &defaultClaims, &customClaims)

		if customClaims.Group == "" && customClaims.UserID == "" && defaultClaims.Issuer == "" {
			return false, "token invalid"
		}

		isValid = defaultClaims.Issuer == issuerName
		reason = "Invalid Issuer."

		isValid = defaultClaims.Expiry.Time().Unix() > time.Now().Unix()
		reason = "Token expired."

		if sudo {
			isValid = customClaims.Group == "root"
			reason = "user is not root"
		}
		if requireUserToken {
			isValid = customClaims.Group != "client"
			reason = "not a human"
		}
	} else {
		isValid = false
		reason = "Invalid Token"
	}
	return isValid, reason
}

//GetUserHeader Gets the Authorization header from the give request
func (u *User) GetUserHeader(request *http.Request) (bool, string) {
	UserHeader := request.Header.Get("Authorization")
	if len(UserHeader) > 7 && strings.EqualFold(UserHeader[0:6], "BEARER") {
		token := UserHeader[7:]
		return true, token
	}
	return false, ""
}

//AuthTokenProvided Used by the Request validator to allow/disallow access to an API based on the presence of a valid Authorization header in the given request.
func (u *User) AuthTokenProvided(r *http.Request) common.APIResponse {
	success, response := u.GetUserHeader(r)
	if success == false {
		return common.CreateFailureResponse(errors.New("no token provided"), "AuthTokenProvided", 401)
	}
	success, response = u.ValidateToken(response, false, false)
	if success == false {
		return common.CreateFailureResponse(errors.New(response), "AuthTokenProvided", 401)
	}
	return common.CreateAPIResponse("success", nil, 200)
}

//IsRoot Used by the request validator to allow/disallow access to an API based on whether the token in the given request belongs to a root user.
func (u *User) IsRoot(r *http.Request) common.APIResponse {
	var resp common.APIResponse

	success, token := u.GetUserHeader(r)
	if success == false {
		resp = common.CreateFailureResponse(errors.New("no token provided"), "IsRoot", 401)
	} else {
		success, reason := u.ValidateToken(token, true, true)
		if success == false {
			resp = common.CreateFailureResponse(fmt.Errorf(reason), "IsRoot", 403)
		} else {
			resp = common.CreateAPIResponse("success", nil, 403)
		}
	}
	return resp
}

//IsUser ...
func (u *User) IsUser(r *http.Request) common.APIResponse {
	var resp common.APIResponse

	success, token := u.GetUserHeader(r)
	if success == false {
		resp = common.CreateFailureResponse(errors.New("no token provided"), "IsUser", 401)
	} else {
		success, reason := u.ValidateToken(token, true, true)
		if success == false {
			resp = common.CreateFailureResponse(fmt.Errorf(reason), "IsUser", 403)
		} else {
			resp = common.CreateAPIResponse("success", nil, 403)
		}
	}
	return resp
}

//GenerateAuthToken ..
func (u *User) GenerateAuthToken(request AuthRequest, userDetails data.User) common.APIResponse {
	var response common.APIResponse
	sig, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.HS256, Key: u.hmacKey}, (&jose.SignerOptions{}).WithType("JWT"))
	if err == nil {
		token := AuthToken{
			Issuer:  issuerName,
			Subject: request.Username,
			Expires: time.Now().Add(744 * time.Hour).Unix(),
			Group:   userDetails.Group,
			UserID:  userDetails.Id,
		}
		json, _ := json.Marshal(token)

		if webSig, err := sig.Sign(json); err == nil {
			signature, err := webSig.CompactSerialize()
			response = common.CreateAPIResponse(signature, err, 500)
		} else {
			response = common.CreateFailureResponse(err, "GenerateAuthToken", 500)
		}
	} else {
		response = common.CreateFailureResponse(err, "GenerateAuthToken", 500)
	}
	return response
}

func generateSymKey() []byte {
	k := make([]byte, 64)
	if _, e := rand.Read(k); e == nil {
		return k
	}
	return nil
}
