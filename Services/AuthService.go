package services

import (
	"fmt"

	DAO "github.com/shogunx2/AuthService/dao"
)

type AuthRequest struct {
	ApiKeyValid bool   `json:"api_key_valid"`
	ApiKey      string `json:"api_key"`
	UserId      string `json:"user_id"`
	Password    string `json:"Password"`
}

type AuthResponse struct {
	ApiKeyValid bool
	ApiKey      string
	UserId      string
	Password    string
}

type AuthServiceIf interface {
	Init(ad DAO.AuthDatastore)
	Add(authReq *AuthRequest) (*AuthRequest, error)
	Remove(authReq *AuthRequest) (*AuthRequest, error)
	// UpdatePassword
	Authenticate(authReq *AuthRequest) (bool, error)
}

func mapAuthRequestToAuthRecord(authReq *AuthRequest) *DAO.AuthRecord {
	return &DAO.AuthRecord{
		ApiKeyValid: authReq.ApiKeyValid,
		ApiKey:      authReq.ApiKey,
		UserId:      authReq.UserId,
		Password:    authReq.Password,
	}
}

func mapAuthRecordToAuthResponse(authRecord *DAO.AuthRecord) *AuthResponse {
	return &AuthResponse{
		ApiKeyValid: authRecord.ApiKeyValid,
		ApiKey:      authRecord.ApiKey,
		UserId:      authRecord.UserId,
		Password:    authRecord.Password,
	}
}

type AuthService struct {
	authDatastore DAO.AuthDatastore
}

func (as *AuthService) Init(ad DAO.AuthDatastore) {
	as.authDatastore = ad
}

func (as *AuthService) Add(authReq *AuthRequest) (*AuthResponse, error) {
	aRec := mapAuthRequestToAuthRecord(authReq)
	_, err := as.authDatastore.Insert(aRec)
	fmt.Println("Add: ", aRec)
	fmt.Println("Existing Add err: ", err)
	authRsp := mapAuthRecordToAuthResponse(aRec)
	return authRsp, err
}

func (as *AuthService) Remove(authReq *AuthRequest) (*AuthResponse, error) {
	aRec := mapAuthRequestToAuthRecord(authReq)
	_, err := as.authDatastore.Remove(aRec)
	fmt.Println("Existing Remove err: ", err)
	authRsp := mapAuthRecordToAuthResponse(aRec)
	return authRsp, err
}

func (as *AuthService) Authenticate(authReq *AuthRequest) (bool, error) {
	aRec := mapAuthRequestToAuthRecord(authReq)
	aRec, err := as.authDatastore.Get(aRec)
	if err != nil {
		fmt.Println("Existing Authenticate (false) err: ", err)
		return false, err
	}
	if aRec.Password == authReq.Password {
		fmt.Println("Existing Authenticate (true) err: ", err)
	} else {
		fmt.Println("Existing Authenticate (false) err: ", err)
		return false, err
	}
	return true, nil
}
