package film

import (
	"fmt"
	"strings"
	"time"

	"github.com/Coderovshik/film-library/internal/util"
)

func ToQueryableLists(fa *FilmActors, format string) (string, []any) {
	n := len(fa.ActorIDs)

	values := make([]any, 0, n+1)
	values = append(values, fa.ID)

	qs := make([]string, 0, n)

	for i := 0; i < n; i++ {
		qs = append(qs, fmt.Sprintf(format, i+2))
		values = append(values, fa.ActorIDs[i])
	}

	return strings.Join(qs, ", "), values
}

func ToQueryableObject(a *Film) *util.QueryableObject {
	qo := util.NewQueryableObject()

	if len(a.Name) != 0 {
		qo.Add("movie_name", a.Name)
	}

	if len(a.Description) != 0 {
		qo.Add("movie_description", a.Description)
	}

	if !a.ReleaseDate.IsZero() {
		qo.Add("releasedate", a.ReleaseDate)
	}

	if a.Rating >= 0 {
		qo.Add("rating", a.Rating)
	}

	return qo
}

var sortMap = map[string]string{
	"name":        "movie_name",
	"rating":      "rating",
	"releasedate": "releasedate",
}

func ToQueryConditions(q *Query) [3]string {
	var conditions [3]string

	var filmCon string
	if len(q.Film) != 0 {
		pattern := "'%" + q.Film + "%'"
		filmCon = "WHERE movie_name LIKE " + pattern
	}
	conditions[0] = filmCon

	var actorCon string
	if len(q.Actor) != 0 {
		pattern := "'%" + q.Actor + "%'"
		actorCon = "HAVING STRING_AGG (a.actor_name, ';') LIKE " + pattern
	}
	conditions[1] = actorCon

	var sortCon string
	if q.Sort == nil {
		sortCon = "ORDER BY rating DESC"
	} else {
		sortCon = fmt.Sprintf("ORDER BY %s %s",
			sortMap[q.Sort[0]], strings.ToUpper(q.Sort[1]))
	}
	conditions[2] = sortCon

	return conditions
}

func ToQuery(req *GetFilmsRequest) *Query {
	var sort []string
	if len(req.SortQuery) != 0 {
		sort = strings.Split(req.SortQuery, ",")
	}

	return &Query{
		Sort:  sort,
		Film:  req.FilmQuery,
		Actor: req.ActorQuery,
	}
}

func ToFilmResponse(f *Film) *FilmResponse {
	return &FilmResponse{
		ID: int(f.ID),
		Info: FilmInfo{
			Name:        f.Name,
			Description: f.Description,
			ReleaseDate: f.ReleaseDate.Format("2006-01-02"),
			Rating:      int(f.Rating),
		},
		Actors: f.Actors,
	}
}

func ToFilm(fi *FilmInfo) *Film {
	releaseDate, _ := time.Parse("2006-01-02", fi.ReleaseDate)

	return &Film{
		Name:        fi.Name,
		Description: fi.Description,
		ReleaseDate: releaseDate,
		Rating:      int32(fi.Rating),
	}
}

func ToActorIDs32(ids []int) []int32 {
	a := make([]int32, 0, len(ids))
	for _, v := range ids {
		a = append(a, int32(v))
	}

	return a
}

func ToActorsShortRespose(a []*ActorShort) []*ActorShortResponse {
	res := make([]*ActorShortResponse, 0, len(a))
	for _, v := range a {
		res = append(res, &ActorShortResponse{
			ID:   int(v.ID),
			Name: v.Name,
		})
	}

	return res
}
