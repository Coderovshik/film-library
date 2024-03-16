package actor

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/Coderovshik/film-library/internal/util"
	gomock "go.uber.org/mock/gomock"
)

func TestService_GetActors(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := NewMockActorRepository(ctrl)

	s := NewService(m)

	// valid request
	bd1, _ := time.Parse(time.DateOnly, "1995-05-03")
	bd2, _ := time.Parse(time.DateOnly, "1996-06-04")
	bd3, _ := time.Parse(time.DateOnly, "1997-07-05")
	out := []*Actor{
		{
			ID:       1,
			Name:     "actor1",
			Sex:      "male",
			Birthday: bd1,
			Films:    []string{"film1", "film2"},
		},
		{
			ID:       2,
			Name:     "actor2",
			Sex:      "female",
			Birthday: bd2,
			Films:    []string{"film2", "film3", "film4"},
		},
		{
			ID:       3,
			Name:     "actor3",
			Sex:      "male",
			Birthday: bd3,
			Films:    []string{"film1"},
		},
	}
	m.EXPECT().
		GetActors(gomock.Any()).
		Return(out, nil).Times(1)

	res, err := s.GetActors(context.TODO())
	if err != nil {
		t.Fatalf("No error expected, got %s", err.Error())
	}
	for i := range out {
		ar := ToActorResponse(out[i])
		if !reflect.DeepEqual(ar, res[i]) {
			t.Errorf("Item %d: expected %+v, got %+v", i, ar, res[i])
		}
	}
}

func TestService_AddActor(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := NewMockActorRepository(ctrl)

	s := NewService(m)

	// valid request
	bd, _ := time.Parse(time.DateOnly, "1995-05-03")
	in := &Actor{
		Name:     "actor1",
		Sex:      "male",
		Birthday: bd,
	}
	out := &Actor{
		ID:       1,
		Name:     "actor1",
		Sex:      "male",
		Birthday: bd,
	}
	m.EXPECT().
		AddActor(gomock.Any(), gomock.Eq(in)).
		Return(out, nil).Times(1)

	req := &ActorInfo{
		Name:     "actor1",
		Sex:      "male",
		Birthday: "1995-05-03",
	}
	expRes := ToActorResponse(out)

	res, err := s.AddActor(context.TODO(), req)
	if err != nil {
		t.Fatalf("No error expected, got %s", err.Error())
	}

	if !reflect.DeepEqual(expRes, res) {
		t.Errorf("Expected %+v, got %+v", expRes, res)
	}

	// invalid (empty) request
	req = &ActorInfo{
		Name:     "",
		Sex:      "",
		Birthday: "",
	}

	_, err = s.AddActor(context.TODO(), req)
	vErr := &util.ValidationError{}
	if !errors.As(err, &vErr) {
		t.Fatalf("Expected %s, got %s", vErr.Error(), err.Error())
	}

	// invalid (format) request
	req = &ActorInfo{
		Name:     "actor1",
		Sex:      "apple",
		Birthday: "2021.25.32",
	}

	_, err = s.AddActor(context.TODO(), req)
	vErr = &util.ValidationError{}
	if !errors.As(err, &vErr) {
		t.Fatalf("Expected %s, got %s", vErr.Error(), err.Error())
	}
}

func TestService_GetActor(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := NewMockActorRepository(ctrl)

	s := NewService(m)

	// valid requset
	in := int32(1)
	bd, _ := time.Parse(time.DateOnly, "1995-05-03")
	out := &Actor{
		ID:       1,
		Name:     "actor1",
		Sex:      "male",
		Birthday: bd,
		Films:    []string{"film1", "film2"},
	}

	m.EXPECT().
		GetActor(gomock.Any(), gomock.Eq(in)).
		Return(out, nil).Times(1)

	req := &ActorIdRequest{
		ID: "1",
	}
	expRes := ToActorResponse(out)

	res, err := s.GetActor(context.TODO(), req)
	if err != nil {
		t.Fatalf("No error expected, got %s", err.Error())
	}

	if !reflect.DeepEqual(expRes, res) {
		t.Errorf("Expected %+v, got %+v", expRes, res)
	}

	// actor does not exist
	in = int32(1)

	m.EXPECT().
		GetActor(gomock.Any(), gomock.Eq(in)).
		Return(nil, ErrActorNotExist).Times(1)

	req = &ActorIdRequest{
		ID: "1",
	}

	_, err = s.GetActor(context.TODO(), req)
	if !errors.Is(err, ErrActorNotExist) {
		t.Fatalf("Expected %s, got %s", ErrActorNotExist.Error(), err.Error())
	}

	// invalid actor id
	req = &ActorIdRequest{
		ID: "apple",
	}

	_, err = s.GetActor(context.TODO(), req)
	if !errors.Is(err, ErrIdInvalid) {
		t.Fatalf("Expected %s, got %s", ErrIdInvalid.Error(), err.Error())
	}
}

func TestService_UpdateActor(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := NewMockActorRepository(ctrl)

	s := NewService(m)

	//valid request
	bd, _ := time.Parse(time.DateOnly, "1995-05-03")
	in := &Actor{
		ID:       1,
		Birthday: bd,
	}
	out := &Actor{
		ID:       1,
		Name:     "actor1",
		Sex:      "male",
		Birthday: bd,
		Films:    []string{"film1", "film2"},
	}

	gomock.InOrder(
		m.EXPECT().UpdateActor(gomock.Any(), gomock.Eq(in)).Return(nil).Times(1),
		m.EXPECT().GetActor(gomock.Any(), gomock.Eq(in.ID)).Return(out, nil).Times(1),
	)

	req := &ActorIdInfoRequest{
		ID: "1",
		Info: ActorInfo{
			Birthday: bd.Format(time.DateOnly),
		},
	}
	expRes := ToActorResponse(out)

	res, err := s.UpdateActor(context.TODO(), req)
	if err != nil {
		t.Fatalf("No error expected, got %s", err.Error())
	}

	if !reflect.DeepEqual(expRes, res) {
		t.Errorf("Expected %+v, got %+v", expRes, res)
	}

	// empty update
	in = &Actor{
		ID: 1,
	}

	m.EXPECT().
		UpdateActor(gomock.Any(), gomock.Eq(in)).
		Return(ErrEmptyUpdate).Times(1)

	req = &ActorIdInfoRequest{
		ID: "1",
	}

	_, err = s.UpdateActor(context.TODO(), req)
	if !errors.Is(err, ErrEmptyUpdate) {
		t.Fatalf("Expected %s, got %s", ErrEmptyUpdate.Error(), err.Error())
	}

	// actor does not exist
	in = &Actor{
		ID:   6969,
		Name: "new name",
	}

	m.EXPECT().
		UpdateActor(gomock.Any(), gomock.Eq(in)).
		Return(ErrActorNotExist).Times(1)

	req = &ActorIdInfoRequest{
		ID: "6969",
		Info: ActorInfo{
			Name: "new name",
		},
	}

	_, err = s.UpdateActor(context.TODO(), req)
	if !errors.Is(err, ErrActorNotExist) {
		t.Fatalf("Expected %s, got %s", ErrActorNotExist.Error(), err.Error())
	}

	// invalid actor id

	req = &ActorIdInfoRequest{
		ID: "apple",
		Info: ActorInfo{
			Name: "new name",
		},
	}

	_, err = s.UpdateActor(context.TODO(), req)
	if !errors.Is(err, ErrIdInvalid) {
		t.Fatalf("Expected %s, got %s", ErrIdInvalid.Error(), err.Error())
	}

	// invalid (format) request
	req = &ActorIdInfoRequest{
		ID: "1",
		Info: ActorInfo{
			Birthday: "apple",
		},
	}

	_, err = s.UpdateActor(context.TODO(), req)
	vErr := &util.ValidationError{}
	if !errors.As(err, &vErr) {
		t.Fatalf("Expected %s, got %s", vErr.Error(), err.Error())
	}
}

func TestService_DeleteActor(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := NewMockActorRepository(ctrl)

	s := NewService(m)

	//valid request
	in := int32(1)
	bd, _ := time.Parse(time.DateOnly, "1995-05-03")
	out := &Actor{
		ID:       1,
		Name:     "actor1",
		Sex:      "male",
		Birthday: bd,
		Films:    []string{"film1", "film2"},
	}

	gomock.InOrder(
		m.EXPECT().GetActor(gomock.Any(), gomock.Eq(in)).Return(out, nil).Times(1),
		m.EXPECT().DeleteActor(gomock.Any(), gomock.Eq(in)).Return(nil).Times(1),
	)

	req := &ActorIdRequest{
		ID: "1",
	}
	expRes := ToActorResponse(out)

	res, err := s.DeleteActor(context.TODO(), req)
	if err != nil {
		t.Fatalf("No error expected, got %s", err.Error())
	}

	if !reflect.DeepEqual(expRes, res) {
		t.Errorf("Expected %+v, got %+v", expRes, res)
	}

	// actor does not exist
	in = int32(1)

	m.EXPECT().
		GetActor(gomock.Any(), gomock.Eq(in)).
		Return(nil, ErrActorNotExist).Times(1)

	req = &ActorIdRequest{
		ID: "1",
	}

	_, err = s.DeleteActor(context.TODO(), req)
	if !errors.Is(err, ErrActorNotExist) {
		t.Fatalf("Expected %s, got %s", ErrActorNotExist.Error(), err.Error())
	}

	// invalid actor id
	req = &ActorIdRequest{
		ID: "-1",
	}

	_, err = s.DeleteActor(context.TODO(), req)
	if !errors.Is(err, ErrIdInvalid) {
		t.Fatalf("Expected %s, got %s", ErrIdInvalid.Error(), err.Error())
	}
}
