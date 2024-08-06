package main

import (
	"errors"
	"fmt"
)

type AuthRequest struct {
	apiKeyValid bool
	apiKey      string
	userId      string
	password    string
}

type AuthResponse struct {
	apiKeyValid bool
	apiKey      string
	userId      string
	password    string
}

type AuthServiceIf interface {
	Init(ad AuthDatastore)
	Add(authReq *AuthRequest) (*AuthRequest, error)
	Remove(authReq *AuthRequest) (*AuthRequest, error)
	// UpdatePassword
	Authenticate(authReq *AuthRequest) (bool, error)
}

func mapAuthRequestToAuthRecord(authReq *AuthRequest) *AuthRecord {
	return &AuthRecord{
		apiKeyValid: authReq.apiKeyValid,
		apiKey:      authReq.apiKey,
		userId:      authReq.userId,
		password:    authReq.password,
	}
}

func mapAuthRecordToAuthResponse(authRecord *AuthRecord) *AuthResponse {
	return &AuthResponse{
		apiKeyValid: authRecord.apiKeyValid,
		apiKey:      authRecord.apiKey,
		userId:      authRecord.userId,
		password:    authRecord.password,
	}
}

type AuthService struct {
	authDatastore AuthDatastore
}

func (as *AuthService) Init(ad AuthDatastore) {
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
	fmt.Println("Existing Authenticate (true) err: ", err)
	return true, nil
}

/*
 * Map containing Auth Record
 * Key : apiKey or userId
 * Value: AuthRecord  (also has apiKey or User ID which duplication of key)
 * TBD: Prepare a new record which is having apiKeyValue, apiKey and password only
 * and the key to the database should be user id
 * TBD: implement update password
 * TBD: reissue API Key
 * Divide the project into 3 files: main.go, auth.go, auth_datastore.go
 */
type AuthRecord struct {
	apiKeyValid bool
	apiKey      string
	userId      string
	password    string
}

// TBD: Read about Go Interfaces
type AuthDatastore interface {
	Init()
	Insert(authRecord *AuthRecord) (*AuthRecord, error)
	Get(authRecord *AuthRecord) (*AuthRecord, error)
	Remove(authRecord *AuthRecord) (*AuthRecord, error)
	Update(authRecord *AuthRecord) (*AuthRecord, error)
}

type AuthMapDatastore struct {
	authMap map[string]AuthRecord
}

func printAuthRecord(key string, ar *AuthRecord) {
	fmt.Println("key: ", key)
	fmt.Println("apiKeyValid: ", ar.apiKeyValid)
	fmt.Println("apiKey: ", ar.apiKey)
	fmt.Println("userId: ", ar.userId)
	fmt.Println("password: ", ar.password)
}

func (amd *AuthMapDatastore) DumpDB() {
	fmt.Println("--- Entered DumpDB ---")
	for k, v := range amd.authMap {
		printAuthRecord(k, &v)
	}
	fmt.Println("--- Exiting DumpDB ---")
}

func (amd *AuthMapDatastore) Init() {
	// make a map of string to AuthRecord
	amd.authMap = make(map[string]AuthRecord)
}

func (amd *AuthMapDatastore) Insert(authRecord *AuthRecord) (*AuthRecord, error) {
	fmt.Println("Entered Insert")
	// check if the authRecord exists
	var ok bool
	if authRecord.apiKeyValid {
		_, ok = amd.authMap[authRecord.apiKey]
	} else {
		_, ok = amd.authMap[authRecord.userId]
	}
	// if exists return error
	// else insert the authRecord into the map
	// return the authRecord
	if ok {
		return nil, errors.New("AuthRecord already exists")
	}
	if authRecord.apiKeyValid {
		amd.authMap[authRecord.apiKey] = *authRecord
	} else {
		amd.authMap[authRecord.userId] = *authRecord
	}
	amd.DumpDB()
	fmt.Println("Exiting Insert")
	return authRecord, nil
}

func (amd *AuthMapDatastore) Get(authRecordIn *AuthRecord) (*AuthRecord, error) {
	fmt.Println("Entered Get")
	// check if the authRecord exists
	var ok bool
	var authRecord AuthRecord
	if authRecordIn.apiKeyValid {
		authRecord, ok = amd.authMap[authRecordIn.apiKey]
	} else {
		authRecord, ok = amd.authMap[authRecordIn.userId]
	}
	// no record found, return error
	if !ok {
		amd.DumpDB()
		return nil, errors.New("AuthRecord does not exist")
	}
	// record exists return the authRecord
	amd.DumpDB()
	fmt.Println("Exiting Get")
	return &authRecord, nil
}

func (amd *AuthMapDatastore) Remove(authRecord *AuthRecord) (*AuthRecord, error) {
	fmt.Println("Entered Remove")
	// check if the authRecord exists
	// if exists remove the authRecord
	// else return error
	if _, ok := amd.authMap[authRecord.apiKey]; ok {
		delete(amd.authMap, authRecord.apiKey)
		amd.DumpDB()
		return authRecord, nil
	}
	amd.DumpDB()
	fmt.Println("Exiting Remove")
	return nil, errors.New("AuthRecord does not exist")
}

func (amd *AuthMapDatastore) Update(authRecord *AuthRecord) (*AuthRecord, error) {
	fmt.Println("Entered Update")
	// check if the authRecord exists
	// if exists update the authRecord
	// else return error
	if _, ok := amd.authMap[authRecord.apiKey]; ok {
		amd.authMap[authRecord.apiKey] = *authRecord
		amd.DumpDB()
		return authRecord, nil
	}
	amd.DumpDB()
	fmt.Println("Exiting Update")
	return nil, errors.New("AuthRecord does not exist")
}

/*
 * Pending activitie:
 * 1. Map key should be User ID
 * And each record can contain User's API Key and User's Password
 *
 */
func main() {
	fmt.Println("Entered main")

	var first int
	var apiKey string
	var user string
	var password string

	// init Authenticaion Map Datastore
	amd := AuthMapDatastore{}
	amd.Init()

	as := AuthService{}
	as.Init(&amd)

	for {
		fmt.Println("Enter 1 to add an api-key")
		fmt.Println("Enter 2 to add a user id and password")
		fmt.Println("Enter 3 to authenticate with an api-key")
		fmt.Println("Enter 4 to authenticate with a user id and password")
		fmt.Scanln(&first)
		fmt.Println("first: ", first)
		switch first {
		case 1:
			// TBD: correct the msg and say API Key
			fmt.Println("Add an auth-key")
			fmt.Scanln(&apiKey)
			ar := AuthRequest{apiKey: apiKey, apiKeyValid: true}
			arOut, err := as.Add(&ar)
			if err != nil {
				fmt.Println("Add error: ", err)
			} else {
				fmt.Println("Add result: ", arOut)
			}
		case 2:
			fmt.Println("Enter user Id")
			fmt.Scanln(&user)
			fmt.Println("Enter password")
			fmt.Scanln(&password)
			ar := AuthRequest{userId: user, apiKeyValid: false, password: password}
			fmt.Println("Invoking Add with user id and password")
			arOut, err := as.Add(&ar)
			if err != nil {
				fmt.Println("Add error: ", err)
			} else {
				fmt.Println("Add result: ", arOut)
			}
		case 3:
			fmt.Println("Enter auth-key to authenticate")
			fmt.Scanln(&apiKey)
			ar := AuthRequest{apiKey: apiKey, apiKeyValid: true}
			authResult, err := as.Authenticate(&ar)
			if err != nil {
				fmt.Println("auth error: ", err)
			} else {
				fmt.Println("auth result: ", authResult)
			}
		case 4:
			fmt.Println("Case 4")
			fmt.Println("Enter user Id")
			fmt.Scanln(&user)
			fmt.Println("Enter password")
			fmt.Scanln(&password)
			ar := AuthRequest{userId: user, apiKeyValid: false, password: password}
			authResult, err := as.Authenticate(&ar)
			if err != nil {
				fmt.Println("auth error: ", err)
			} else {
				fmt.Println("auth result: ", authResult)
			}
		default:
		}
	}
}
