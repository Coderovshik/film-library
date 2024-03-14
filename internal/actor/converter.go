package actor

import (
	"time"

	"github.com/Coderovshik/film-library/internal/util"
)

func ToQueryableObject(a *Actor) *util.QueryableObject {
	qo := util.NewQueryableObject()

	if len(a.Name) != 0 {
		qo.Add("actor_name", a.Name)
	}

	if len(a.Sex) != 0 {
		qo.Add("sex", a.Sex)
	}

	if !a.Birthday.IsZero() {
		qo.Add("birthday", a.Birthday)
	}

	return qo
}

func ToActorResponse(a *Actor) *ActorResponse {
	return &ActorResponse{
		ID: int(a.ID),
		Info: ActorInfo{
			Name:     a.Name,
			Sex:      a.Sex,
			Birthday: a.Birthday.Format("2006-01-02"),
		},
		Films: a.Films,
	}
}

func ToActor(ai *ActorInfo) *Actor {
	birthday, _ := time.Parse("2006-01-02", ai.Birthday)

	return &Actor{
		Name:     ai.Name,
		Sex:      ai.Name,
		Birthday: birthday,
	}
}
