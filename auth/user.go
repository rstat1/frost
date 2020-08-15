package auth

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.alargerobot.dev/frost/common"
	"go.alargerobot.dev/frost/data"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

var (
	issuerName string = "https://trinity.alargerobot.dev"
)

//User ..
type User struct {
	datastore *data.DataStore
	hmacKey   []byte
}

//AuthToken ...
type AuthToken struct {
	Issuer      string `json:"iss"`
	Subject     string `json:"sub"`
	Expires     int64  `json:"exp"`
	ForApp      string `json:"app"`
	UserID      string `json:"uid"`
	AccessLevel string `json:"lvl"`
}

//NewUserService ...
func NewUserService(db *data.DataStore) *User {
	return &User{
		datastore: db,
		hmacKey:   generateSymKey(),
	}
}

//NewUser Creates a new user and returns an auth token.
func (u *User) NewUser(request data.AuthRequest, p []data.ServiceAuth) common.APIResponse {
	var apiResp common.APIResponse
	if _, err := u.datastore.NewUser(request, p); err == nil {
		apiResp = common.CreateAPIResponse("success", nil, 400)
	} else {
		apiResp = common.CreateFailureResponse(err, "NewUser", 400)
	}

	return apiResp
}

//GetUserFromToken ...
func (u *User) GetUserFromToken(r *http.Request) data.User {
	var user data.User
	success, response := u.GetUserHeader(r)
	if success {
		if token, err := jwt.ParseSigned(response); err == nil {
			var defaultClaims jwt.Claims
			username := struct {
				Name        string `json:"sub"`
				AccessLevel string `json:"lvl"`
			}{}
			token.Claims(u.hmacKey, &defaultClaims, &username)
			user = u.datastore.GetUser(username.Name)
			user.Group = username.AccessLevel
		} else {
			common.Logger.WithField("func", "GetUsernameFromToken").Errorln(errors.New("invalid token"))
			user = data.User{}
		}
	} else {
		common.Logger.WithField("func", "GetUsernameFromToken").Errorln(errors.New("token not provided"))
		user = data.User{}
	}
	user.PassHash = ""
	return user
}

//GetUserFromProvidedToken ...
func (u *User) GetUserFromProvidedToken(token, serviceName string) data.User {
	var user data.User
	if token, err := jwt.ParseSigned(token); err == nil {
		var defaultClaims jwt.Claims
		tokenClaims := struct {
			Name        string `json:"sub"`
			AccessLevel string `json:"lvl"`
			AppName     string `json:"app"`
		}{}
		token.Claims(u.hmacKey, &defaultClaims, &tokenClaims)
		user = u.datastore.GetUser(tokenClaims.Name)
		user.Group = tokenClaims.AccessLevel
		// if tokenClaims.AppName != serviceName {
		// 	common.LogError("", errors.New("invalid token"))
		// 	user = data.User{}
		// }
	} else {
		common.LogError("", errors.New("invalid token"))
		user = data.User{}
	}
	user.PassHash = ""
	return user
}

//GetUserByID ...
func (u *User) GetUserByID(id string) data.User {
	return u.datastore.GetUserByID(id)
}

//ValidateLoginRequest Checks that the provided username and password are known and returns a signed JWT. Also checks for site access
//permission.
func (u *User) ValidateLoginRequest(request data.AuthRequest) common.APIResponse {
	var userInfo data.User
	var response common.APIResponse

	passHasher := sha512.New512_256()
	hash := hex.EncodeToString(passHasher.Sum([]byte(request.Password)))
	userInfo = u.datastore.GetUser(request.Username)
	if userInfo.PassHash == string(hash) {
		response = common.CreateAPIResponse("success", nil, 500) //u.GenerateAuthToken(request, userInfo)
	} else {
		response = common.CreateFailureResponse(errors.New("incorrect credentials"), "ValidateLoginRequest", 401)
	}
	return response
}

//ValidateToken Checks that provided token was sent by this service, hasn't expired and was signed by the current instance. Also can check if a user is root or not.
func (u *User) ValidateToken(token string, sudo bool, requireUserToken bool) (bool, string) {
	var user data.User
	if token, err := jwt.ParseSigned(token); err == nil {
		var defaultClaims jwt.Claims
		customClaims := struct {
			ForApp string `json:"app"`
			UserID string `json:"uid"`
			Group  string `json:"lvl"`
		}{}
		token.Claims(u.hmacKey, &defaultClaims, &customClaims)

		if customClaims.ForApp == "" && customClaims.UserID == "" && customClaims.Group == "" && defaultClaims.Issuer == "" {
			return false, "token invalid"
		}

		user = u.datastore.GetUserByID(customClaims.UserID)
		if user.Id == "" {
			return false, "User ID doesn't belong to a valid user."
		}

		if defaultClaims.Issuer != issuerName {
			return false, "Invalid Issuer."
		}

		if defaultClaims.Expiry.Time().Unix() < time.Now().Unix() {
			return false, "Token expired."
		}

		if sudo {
			if customClaims.Group != "root" {
				return false, "user is not root"
			}
		}
		if requireUserToken {
			if customClaims.Group == "client" {
				return false, "not a human"
			}
		}
	} else {
		return false, "Invalid Token"
	}
	return true, ""
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
func (u *User) GenerateAuthToken(username, app string) common.APIResponse {
	var response common.APIResponse
	if userDetails := u.datastore.GetUser(username); userDetails.Username != "" {
		sig, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.HS256, Key: u.hmacKey}, (&jose.SignerOptions{}).WithType("JWT"))
		if err == nil {
			token := AuthToken{
				Issuer:      issuerName,
				Subject:     username,
				Expires:     time.Now().Add(744 * time.Hour).Unix(),
				ForApp:      app,
				UserID:      userDetails.Id,
				AccessLevel: u.getAppAccessForUser(username, app),
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
	} else {
		return common.CreateFailureResponse(errors.New("could find user with that name"), "GenerateAuthToken", 500)
	}
}

//ParsePermissionList ...
func (u *User) ParsePermissionList(p []data.ServiceAuth) {
	common.Logger.Debugln(p[0])
}
func (u *User) getAppAccessForUser(user, app string) string {
	if u.datastore.DoesUserHavePermission(user, app, "hasRoot") == true {
		return "root"
	}
	return "user"
}
func generateSymKey() []byte {
	k := make([]byte, 64)
	if _, e := rand.Read(k); e == nil {
		return k
	}
	return nil
}
