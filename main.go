package main

import (
	"errors"
	"fmt"
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
	Init(ad AuthDatastore)
	Add(authReq *AuthRequest) (*AuthRequest, error)
	Remove(authReq *AuthRequest) (*AuthRequest, error)
	// UpdatePassword
	Authenticate(authReq *AuthRequest) (bool, error)
}

func mapAuthRequestToAuthRecord(authReq *AuthRequest) *AuthRecord {
	return &AuthRecord{
		ApiKeyValid: authReq.ApiKeyValid,
		ApiKey:      authReq.ApiKey,
		UserId:      authReq.UserId,
		Password:    authReq.Password,
	}
}

func mapAuthRecordToAuthResponse(authRecord *AuthRecord) *AuthResponse {
	return &AuthResponse{
		ApiKeyValid: authRecord.ApiKeyValid,
		ApiKey:      authRecord.ApiKey,
		UserId:      authRecord.UserId,
		Password:    authRecord.Password,
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
	if aRec.Password == authReq.Password{
		fmt.Println("Existing Authenticate (true) err: ", err)
	}else{
		fmt.Println("Existing Authenticate (false) err: ", err)
		return false, err
	}
	return true, nil
}

/*
 * Map containing Auth Record
 * TBD: implement update Password
 * TBD: reissue API Key
*/
type AuthRecord struct {
	ApiKeyValid bool
	ApiKey      string
	UserId      string
	Password    string
}

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
	fmt.Println("ApiKeyValid: ", ar.ApiKeyValid)
	fmt.Println("ApiKey: ", ar.ApiKey)
	fmt.Println("UserId: ", ar.UserId)
	fmt.Println("Password: ", ar.Password)
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

func (amd *AuthMapDatastore) getRecord(key string) *AuthRecord {
	authRecord, ok := amd.authMap[key]
	if !ok {
		return nil
	}
	return &authRecord
}

func (amd *AuthMapDatastore) setRecord(key string, ar *AuthRecord) {
	fmt.Println("Entering setRecord")
	fmt.Println("ar: ", ar)
	amd.authMap[key] = *ar
	fmt.Println("Exiting setRecord")
}

func (amd *AuthMapDatastore) Insert(authRecord *AuthRecord) (*AuthRecord, error) {
	fmt.Println("Entered Insert")
	var ar *AuthRecord
	var key string
	if authRecord.ApiKeyValid {
		key = authRecord.ApiKey
	} else {
		key = authRecord.UserId
	}
	if len(key) == 0 {
		return nil, errors.New("invalid key")
	}
	// check if the authRecord exists
	fmt.Println("key: ", key)
	ar = amd.getRecord(key)
	// if exists return error
	// return the authRecord
	if ar != nil {
		return nil, errors.New("authRecord already exists")
	}
	fmt.Println("set record")
	// else insert the authRecord into the map
	amd.setRecord(key, authRecord)
	amd.DumpDB()
	fmt.Println("Exiting Insert")
	return ar, nil
}

func (amd *AuthMapDatastore) Get(authRecordIn *AuthRecord) (*AuthRecord, error) {
	fmt.Println("Entered Get")
	// check if the authRecord exists
	var ar *AuthRecord
	var key string
	if authRecordIn.ApiKeyValid {
		key = authRecordIn.ApiKey
	} else {
		key = authRecordIn.UserId
	}
	fmt.Println("key:",key)
	if len(key) == 0 {
		return nil, errors.New("invalid key")
	}
	// check if the authRecord exists
	ar = amd.getRecord(key)
	// no record found, return error
	if ar == nil {
		amd.DumpDB()
		return nil, errors.New("AuthRecord does not exist")
	}
	// record exists return the authRecord
	amd.DumpDB()
	fmt.Println("Exiting Get")
	return ar, nil
}

func (amd *AuthMapDatastore) Remove(authRecord *AuthRecord) (*AuthRecord, error) {
	fmt.Println("Entered Remove")
	var ar *AuthRecord
	var key string

	if authRecord.ApiKeyValid {
		key = authRecord.ApiKey
	} else {
		key = authRecord.UserId
	}
	if len(key) == 0 {
		return nil, errors.New("invalid key")
	}
	// check if the authRecord exists
	ar = amd.getRecord(key)
	// no record found, return error
	if ar == nil {
		amd.DumpDB()
		return nil, errors.New("AuthRecord does not exist")
	}
	// if exists remove the authRecord
	delete(amd.authMap, key)
	amd.DumpDB()
	fmt.Println("Exiting Remove")
	return ar, nil
}

func (amd *AuthMapDatastore) Update(authRecord *AuthRecord) (*AuthRecord, error) {
	fmt.Println("Entered Update")
	var ar *AuthRecord
	var key string

	if authRecord.ApiKeyValid {
		key = authRecord.ApiKey
	} else {
		key = authRecord.UserId
	}
	if len(key) == 0 {
		return nil, errors.New("invalid key")
	}
	// check if the authRecord exists
	ar = amd.getRecord(key)
	// no record found, return error
	if ar == nil {
		amd.DumpDB()
		return nil, errors.New("AuthRecord does not exist")
	}
	amd.authMap[key] = *authRecord
	amd.DumpDB()
	fmt.Println("Exiting Update")
	return authRecord, nil
}

/*
 * Pending activitie:
 * 1. Map key should be User ID
 * And each record can contain User's API Key and User's Password
 *
 */
// func main() {
// 	fmt.Println("Entered main")

// 	var first int
// 	var ApiKey string
// 	var user string
// 	var Password string

// 	// init Authenticaion Map Datastore
// 	amd := AuthMapDatastore{}
// 	amd.Init()

// 	as := AuthService{}
// 	as.Init(&amd)

// 	for {
// 		fmt.Println("Enter 1 to add an api-key")
// 		fmt.Println("Enter 2 to add a user id and Password")
// 		fmt.Println("Enter 3 to authenticate with an api-key")
// 		fmt.Println("Enter 4 to authenticate with a user id and Password")
// 		fmt.Scanln(&first)
// 		fmt.Println("first: ", first)
// 		switch first {
// 		case 1:
// 			// TBD: correct the msg and say API Key
// 			fmt.Println("Add an auth-key")
// 			fmt.Scanln(&ApiKey)
// 			ar := AuthRequest{ApiKey: ApiKey, ApiKeyValid: true}
// 			arOut, err := as.Add(&ar)
// 			if err != nil {
// 				fmt.Println("Add error: ", err)
// 			} else {
// 				fmt.Println("Add result: ", arOut)
// 			}
// 		case 2:
// 			fmt.Println("Enter user Id")
// 			fmt.Scanln(&user)
// 			fmt.Println("Enter Password")
// 			fmt.Scanln(&Password)
// 			ar := AuthRequest{UserId: user, ApiKeyValid: false, Password: Password}
// 			fmt.Println("Invoking Add with user id and Password")
// 			arOut, err := as.Add(&ar)
// 			if err != nil {
// 				fmt.Println("Add error: ", err)
// 			} else {
// 				fmt.Println("Add result: ", arOut)
// 			}
// 		case 3:
// 			fmt.Println("Enter auth-key to authenticate")
// 			fmt.Scanln(&ApiKey)
// 			ar := AuthRequest{ApiKey: ApiKey, ApiKeyValid: true}
// 			authResult, err := as.Authenticate(&ar)
// 			if err != nil {
// 				fmt.Println("auth error: ", err)
// 			} else {
// 				fmt.Println("auth result: ", authResult)
// 			}
// 		case 4:
// 			fmt.Println("Case 4")
// 			fmt.Println("Enter user Id")
// 			fmt.Scanln(&user)
// 			fmt.Println("Enter Password")
// 			fmt.Scanln(&Password)
// 			ar := AuthRequest{UserId: user, ApiKeyValid: false, Password: Password}
// 			authResult, err := as.Authenticate(&ar)
// 			if err != nil {
// 				fmt.Println("auth error: ", err)
// 			} else {
// 				fmt.Println("auth result: ", authResult)
// 			}
// 		default:
// 		}
// 	}
// }
