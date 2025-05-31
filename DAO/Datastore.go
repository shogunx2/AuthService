package dao

//"errors"
//"fmt"

type AuthRecord struct {
	UserId      string
	Password    string
	ApiKey      string
	ApiKeyValid bool
}

type AuthDatastore interface {
	Init()
	Insert(authRecord *AuthRecord) (*AuthRecord, error)
	Get(authRecord *AuthRecord) (*AuthRecord, error)
	Remove(authRecord *AuthRecord) (*AuthRecord, error)
	Update(authRecord *AuthRecord) (*AuthRecord, error)
}
