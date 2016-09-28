package util

import (
	crand "crypto/rand"
	"encoding/binary"
	"log"
	"math/rand"

	"github.com/docker/orca/enzi/api/forms"
	"github.com/stretchr/testify/assert"
)

var Prng *rand.Rand

func init() {
	var seed int64
	if err := binary.Read(crand.Reader, binary.LittleEndian, &seed); err != nil {
		log.Fatalf("unable to get prng seed from /dev/urandom: %s", err)
	}

	Prng = rand.New(rand.NewSource(seed))
}

func GenerateRandomUsername(length uint) string {
	const chars = "abcdefghijklmnopqrstuvwxyz1234567890"
	runes := []rune(chars)

	username := make([]rune, length)
	for i := range username {
		username[i] = runes[Prng.Intn(len(runes))]
	}
	return string(username)
}

func GenerateRandomPassword(length uint) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890[];,./<>?:\"{}-_=+|\\`~*&^%$#@!()한조솨黄胜麟åéîøüπæ≤≥¥ƒç®œ…≈∫∑√Ω˙†˜ß"
	runes := []rune(chars)

	password := make([]rune, length)
	for i := range password {
		password[i] = runes[Prng.Intn(len(runes))]
	}
	return string(password)
}

type User struct {
	Name     string
	Password string
	IsAdmin  bool
}

// func AddFormUsers(newUsers []User, sqlUsers []*adminserver.FormUser) []*adminserver.FormUser {
// 	userIndex := make(map[string]int, len(sqlUsers))
// 	for i, sqlUser := range sqlUsers {
// 		userIndex[sqlUser.Username] = i
// 	}

// 	for _, user := range newUsers {
// 		i, ok := userIndex[user.Name]
// 		if !ok {
// 			i = len(sqlUsers)
// 			userIndex[user.Name] = i
// 			sqlUsers = append(sqlUsers, &adminserver.FormUser{
// 				Username: user.Name,
// 				IsNew:    true,
// 			})
// 		}

// 		// If the values are different then teams have changed.
// 		sqlUsers[i].TeamsChanged = sqlUsers[i].IsAdmin != *user.Account.IsAdmin

// 		sqlUsers[i].Password = user.Password
// 		sqlUsers[i].IsAdmin = *user.Account.IsAdmin
// 	}
// 	return sqlUsers
// }

func (u *Util) CreateActivateRandomUser() (User, func()) {
	retUser := User{GenerateRandomUsername(10), GenerateRandomPassword(10), false}
	retFunc := u.CreateUserWithChecks(retUser.Name, retUser.Password)
	if u.IsSuiteRunningInManagedMode() {
		u.API.EnziSession().UpdateAccount(retUser.Name, forms.UpdateAccount{IsActive: &[]bool{true}[0]})
	}
	return retUser, retFunc
}

func (u *Util) CreateUserWithChecks(username, password string) func() {
	retFunc := func() {}
	u.API.Login(u.Config.AdminUsername, u.Config.AdminPassword)
	if acc, err := u.CreateUser(username, password); err != nil {
		u.T().Fatalf("Failed to create user %s: %s", username, err)
	} else {
		retFunc = func() {
			u.API.Login(u.Config.AdminUsername, u.Config.AdminPassword)
			err := u.API.EnziSession().DeleteAccount(username)
			assert.Nil(u.T(), err, "%s", err)
		}
		assert.NotNil(u.T(), acc)
		assert.Equal(u.T(), username, acc.Name)
		assert.NotEmpty(u.T(), acc.ID)
		assert.False(u.T(), acc.IsOrg)
		if u.IsSuiteRunningInManagedMode() {
			assert.False(u.T(), *acc.IsActive)
		} else if u.IsSuiteRunningInLDAPMode() {
			assert.True(u.T(), *acc.IsActive)
		}
	}
	return retFunc
}
