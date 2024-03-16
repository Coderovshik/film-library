package user

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/Coderovshik/film-library/internal/config"
	"github.com/Coderovshik/film-library/internal/util"
	gomock "go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

type userMatcher struct {
	user *User
}

func UserMatcher(u *User) gomock.Matcher {
	return &userMatcher{u}
}

func (um *userMatcher) Matches(x any) bool {
	other := x.(*User)
	err := bcrypt.CompareHashAndPassword([]byte(other.Passhash), []byte(um.user.Passhash))
	return err == nil && other.Username == um.user.Username && other.IsAdmin == um.user.IsAdmin
}

func (um *userMatcher) String() string {
	return fmt.Sprintf("%+v", um.user)
}

func TestService_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := NewMockUserRepository(ctrl)

	cfg := &config.Config{
		SigningKey: "key",
	}
	s := NewService(m, cfg)

	// valid request
	uIn := &User{
		Username: "user",
		Passhash: "user",
		IsAdmin:  false,
	}
	uOut := &User{
		ID:       1,
		Username: "user",
	}
	m.EXPECT().
		CreateUser(gomock.Any(), UserMatcher(uIn)).
		Return(uOut, nil).Times(1)

	req := &CreateUserRequest{
		Username: "user",
		Password: "user",
	}
	expRes := &CreateUserResponse{
		ID:       int(uOut.ID),
		Username: uOut.Username,
	}

	res, err := s.CreateUser(context.TODO(), req)
	if err != nil {
		t.Fatalf("No error expected, got %s", err.Error())
	}
	if res.ID != int(expRes.ID) || res.Username != expRes.Username {
		t.Errorf("Expected %+v, got %+v", expRes, res)
	}

	// user already exists
	uIn = &User{
		Username: "user",
		Passhash: "user",
		IsAdmin:  false,
	}
	m.EXPECT().
		CreateUser(gomock.Any(), UserMatcher(uIn)).
		Return(nil, ErrUserExist).Times(1)

	req = &CreateUserRequest{
		Username: "user",
		Password: "user",
	}

	_, err = s.CreateUser(context.TODO(), req)
	if !errors.Is(err, ErrUserExist) {
		t.Fatalf("Expected %s, got %s", ErrUserExist.Error(), err.Error())
	}

	// invalid request
	req = &CreateUserRequest{
		Username: "",
		Password: "",
	}

	_, err = s.CreateUser(context.TODO(), req)
	vErr := &util.ValidationError{}
	if !errors.As(err, &vErr) {
		t.Fatalf("Expected %s, got %s", vErr.Error(), err.Error())
	}
}

func TestService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := NewMockUserRepository(ctrl)

	cfg := &config.Config{
		SigningKey: "key",
	}
	s := NewService(m, cfg)

	// valid request
	in := "user"
	out := &User{
		ID:       1,
		Username: "user",
		Passhash: "$2y$10$zoMU2kA9pV4doHHwwPzTIOg746kIGKJc0CFHO9.ES1NzvBiBR0MLO",
		IsAdmin:  false,
	}
	m.EXPECT().
		GetUserByUsername(gomock.Any(), gomock.Eq(in)).
		Return(out, nil).Times(1)

	req := &LoginRequest{
		Username: "user",
		Password: "user",
	}
	expRes := util.NewUserClaims(int(out.ID), out.Username, out.IsAdmin)

	res, err := s.Login(context.TODO(), req)
	if err != nil {
		t.Fatalf("No error expected, got %s", err.Error())
	}

	uc, err := util.ParseUserClaims(res.AccessToken, cfg.SigningKey)
	if err != nil {
		t.Fatalf("No error expected, got %s", err.Error())
	}

	if expRes.ID != uc.ID || expRes.Username != uc.Username || expRes.IsAdmin != uc.IsAdmin {
		t.Errorf("Expected %+v, got %+v", expRes, uc)
	}

	// user does not exist
	in = "user69"
	m.EXPECT().
		GetUserByUsername(gomock.Any(), gomock.Eq(in)).
		Return(nil, ErrUserNotExist).Times(1)

	req = &LoginRequest{
		Username: "user69",
		Password: "6969",
	}

	_, err = s.Login(context.TODO(), req)
	if !errors.Is(err, ErrUserNotExist) {
		t.Fatalf("Expected %s, got %s", ErrUserExist.Error(), err.Error())
	}

	// incorrect password
	in = "user"
	out = &User{
		ID:       1,
		Username: "user",
		Passhash: "$2y$10$zoMU2kA9pV4doHHwwPzTIOg746kIGKJc0CFHO9.ES1NzvBiBR0MLO",
		IsAdmin:  false,
	}
	m.EXPECT().
		GetUserByUsername(gomock.Any(), gomock.Eq(in)).
		Return(out, nil).Times(1)

	req = &LoginRequest{
		Username: "user",
		Password: "resu",
	}

	_, err = s.Login(context.TODO(), req)
	if !errors.Is(err, ErrPasswordIncorrect) {
		t.Fatalf("Expected %s, got %s", ErrUserExist.Error(), err.Error())
	}

	// invalid request
	req = &LoginRequest{
		Username: "",
		Password: "",
	}

	_, err = s.Login(context.TODO(), req)
	vErr := &util.ValidationError{}
	if !errors.As(err, &vErr) {
		t.Fatalf("Expected %s, got %s", vErr.Error(), err.Error())
	}
}
