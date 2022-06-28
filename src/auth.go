package core

//import (
//	"errors"
//	"net/http"
//	"strings"
//	"time"
//
//	mantisEncoding "github.com/sphireinc/mantis/encoding"
//	mantisEncryption "github.com/sphireinc/mantis/encryption"
//)
//
//type AccessToken struct {
//	Scope       string
//	AccessToken string
//	TokenType   string
//	AppId       string
//	ExpiresIn   time.Duration
//	Nonce       string
//}
//
//type Auth struct {
//	Username    string
//	Password    string
//	BearerToken string
//}
//
//// GetAccessToken returns an accessToken response
//func GetAccessToken(w http.ResponseWriter, req *http.Request) {
//	auth, err := getAuthFromHeader(req)
//	if err != nil {
//		App.Log.HandleError("GetAccessToken error", err)
//	}
//	nonce := mantisEncryption.CreateRandomString(16)
//	hash := mantisEncryption.Hash{
//		Input:     mantisEncoding.Base64EncodeUrl(auth.String() + nonce),
//		Algorithm: mantisEncryption.Sha512,
//	}
//	hash.Hash()
//
//	accessToken := AccessToken{
//		Scope:       "*",
//		AccessToken: hash.Output,
//		TokenType:   "Bearer",
//		AppId:       "",
//		ExpiresIn:   30 * time.Minute,
//		Nonce:       nonce,
//	}
//	HandleResponse(w, req, accessToken, http.StatusOK)
//}
//
//// getAuthFromHeader returns an Auth struct filled with the parameters from the Auth header
//func getAuthFromHeader(req *http.Request) (Auth, error) {
//	authorization := req.Header.Get("Authorization")
//	isBasic := strings.Index(authorization, "Basic") >= 0
//	isBearer := strings.Index(authorization, "Bearer") >= 0
//
//	// check for either Basic authorization or Bearer authorization
//	if !isBasic && !isBearer {
//		return Auth{}, errors.New("neither basic nor bearer designation")
//	}
//
//	authParams := strings.Split(authorization, " ")
//
//	auth := Auth{}
//	if isBasic {
//		authParams = strings.Split(authParams[1], ":")
//		auth.Username = authParams[0]
//		auth.Password = authParams[1]
//	}
//
//	if isBearer {
//		auth.BearerToken = authParams[1]
//	}
//
//	return auth, nil
//}
