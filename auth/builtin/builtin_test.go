package builtin

import (
	"fmt"

	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	libkvmock "github.com/docker/libkv/store/mock"
	"github.com/docker/orca/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

const (
	TEST_PASS = "orca"
)

func getTestStore() (store.Store, error) {
	client := "localhost:9999"
	libkv.AddStore("mock", libkvmock.New)
	return libkv.NewStore(
		"mock",
		[]string{client},
		&store.Config{
			ConnectionTimeout: 10 * time.Second,
		},
	)
}

func TestAuthenticate(t *testing.T) {
	kvStore, err := getTestStore()
	if err != nil {
		t.Fatal(err)
	}

	j := []byte(`{"first_name":"Test","last_name":"User","username":"testuser","password":"$2a$10$9S73Yfjr3Fyf7lbU6Mgcn.yT7MxQKGffaOwOzI27./WgQ4xCn01DO","roles":["admin"]}`)

	kvPair := &store.KVPair{
		Key:   "key",
		Value: j,
	}

	s := kvStore.(*libkvmock.Mock)
	s.On("Exists", datastoreAccounts+"/testuser").Return(true, nil)
	s.On("Get", datastoreAccounts+"/testuser").Return(kvPair, nil)
	s.On("List", datastoreAccounts).Return([]*store.KVPair{}, nil)
	s.On("List", datastoreAccountTeams+"/testuser").Return([]*store.KVPair{}, nil)
	s.On("Put", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("*store.WriteOptions")).Return(nil)

	a := NewAuthenticator(&kvStore, "orca", nil, nil)
	token, err := a.AuthenticateUsernamePassword("testuser", TEST_PASS)
	if err != nil {
		t.Error(err)
	}
	if token == nil {
		t.Error("expected password orca")
	}

	token, err = a.AuthenticateUsernamePassword("testuser", "BADpass")
	if err != auth.ErrInvalidPassword {
		t.Error(err)
	}
	if token != nil {
		t.Error("expected passwords to not match")
	}
}

func TestAuthenticateFailIfNoPass(t *testing.T) {
	kvStore, err := getTestStore()
	if err != nil {
		t.Fatal(err)
	}

	j := []byte(`{"first_name":"Test","last_name":"User","username":"testuser","password":"$2a$10$9S73Yfjr3Fyf7lbU6Mgcn.yT7MxQKGffaOwOzI27./WgQ4xCn01DO","roles":["admin"]}`)

	kvPair := &store.KVPair{
		Key:   "key",
		Value: j,
	}

	s := kvStore.(*libkvmock.Mock)
	s.On("Exists", datastoreAccounts+"/testuser").Return(true, nil)
	s.On("Get", datastoreAccounts+"/testuser").Return(kvPair, nil)
	s.On("List", datastoreAccounts).Return([]*store.KVPair{}, nil)
	s.On("List", datastoreAccountTeams+"/testuser").Return([]*store.KVPair{}, nil)
	s.On("Put", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("*store.WriteOptions")).Return(nil)

	a := NewAuthenticator(&kvStore, "orca", nil, nil)
	token, err := a.AuthenticateUsernamePassword("testuser", "")
	if err != auth.ErrInvalidPassword {
		t.Error(err)
	}
	if token != nil {
		t.Error("empty password should not match")
	}
}

func TestGenerateToken(t *testing.T) {
	kvStore, err := getTestStore()
	if err != nil {
		t.Fatal(err)
	}

	s := kvStore.(*libkvmock.Mock)
	s.On("Put", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("*store.WriteOptions")).Return(nil)
	s.On("List", datastoreAccounts).Return([]*store.KVPair{}, nil)

	a := NewAuthenticator(&kvStore, "orca", nil, nil)
	assert.NotNil(t, a)

	tokenStr, err := a.GenerateToken("testuser", nil)
	if err != nil {
		t.Errorf("Error generating token: %s", err)
	}
	assert.NotEqual(t, tokenStr, "", "Didn't get a token string")

	tok, err := a.parseToken(tokenStr)
	assert.Equal(t, tok.Claims["iss"], "orca", "Issuer was set incorrectly")
	assert.Equal(t, tok.Claims["sub"], "testuser", "Subject was set incorrectly")

	t.Logf("exp = %q", tok.Claims["exp"])
	tt := time.Now().Add(TokenExp).Unix()
	exp := tok.Claims["exp"]

	if float64(tt) < exp.(float64) {
		t.Errorf("Token time was set too short")
	}
}

func TestAuthTeam(t *testing.T) {
	kvStore, err := getTestStore()
	if err != nil {
		t.Fatal(err)
	}

	teamName := "testTeam"
	teamID := "12345"

	j := []byte(fmt.Sprintf(`{"id":"%s","name":"%s","description":"Test"}`, teamID, teamName))

	kvPair := &store.KVPair{
		Key:   "key",
		Value: j,
	}

	s := kvStore.(*libkvmock.Mock)
	s.On("Exists", datastoreTeams+"/"+teamID).Return(true, nil)
	s.On("Get", datastoreTeams+"/"+teamID).Return(kvPair, nil)
	s.On("Put", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("*store.WriteOptions")).Return(nil)
	s.On("List", datastoreAccounts).Return([]*store.KVPair{}, nil)
	s.On("List", datastoreTeams).Return([]*store.KVPair{kvPair}, nil)
	s.On("List", datastoreTeamMembers+"/"+teamID).Return([]*store.KVPair{}, nil)

	a := NewAuthenticator(&kvStore, "orca", nil, nil)
	g, err := a.GetTeam(nil, teamID)
	if err != nil {
		t.Fatal(err)
	}

	if g.Name != teamName {
		t.Fatalf("expected team name %s; received %s", teamName, g.Name)
	}

	if g.Id != teamID {
		t.Fatalf("expected team id %s; received %s", teamID, g.Id)
	}
}

func TestAuthTeams(t *testing.T) {
	kvStore, err := getTestStore()
	if err != nil {
		t.Fatal(err)
	}

	teamID := "12345"

	j := []byte(fmt.Sprintf(`{"id":"%s","name":"testTeam","description":"Test"}`, teamID))

	kvPair := &store.KVPair{
		Key:   "key",
		Value: j,
	}

	s := kvStore.(*libkvmock.Mock)
	s.On("Exists", datastoreTeams+"/"+teamID).Return(true, nil)
	s.On("Get", datastoreTeams+"/"+teamID).Return(kvPair, nil)
	s.On("Put", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("*store.WriteOptions")).Return(nil)
	s.On("List", datastoreAccounts).Return([]*store.KVPair{}, nil)
	s.On("List", datastoreTeams).Return([]*store.KVPair{kvPair}, nil)
	s.On("List", datastoreTeamMembers+"/"+teamID).Return([]*store.KVPair{}, nil)

	a := NewAuthenticator(&kvStore, "orca", nil, nil)
	teams, err := a.ListTeams(nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(teams) != 1 {
		t.Fatalf("expected 1 team; received %d", len(teams))
	}

	g := teams[0]
	if g.Id != teamID {
		t.Fatalf("expected team id %s; received %s", teamID, g.Id)
	}
}
