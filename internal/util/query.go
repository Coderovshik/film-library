package util

import (
	"fmt"
	"strings"
)

type QueryableObject struct {
	keys   []string
	values []any
	count  int
}

func NewQueryableObject() *QueryableObject {
	return &QueryableObject{
		keys:   make([]string, 0),
		values: make([]any, 0),
	}
}

func (qo *QueryableObject) IsEmpty() bool {
	return len(qo.keys) == 0
}

func (qo *QueryableObject) Add(key string, val any) {
	qo.keys = append(qo.keys, key)
	qo.values = append(qo.values, val)
}

func (qo *QueryableObject) Args(start int) string {
	args := make([]string, len(qo.keys))
	for i, v := range qo.keys {
		args[i] = fmt.Sprintf("%s = $%d", v, i+start)
	}

	return strings.Join(args, ", ")
}

func (qo *QueryableObject) Values() []any {
	return qo.values
}

func (qo *QueryableObject) Len() int {
	return len(qo.keys)
}
