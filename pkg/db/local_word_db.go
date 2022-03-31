package db

import (
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
)

type LocalWordDB struct {
	words []string
	rand  *rand.Rand
}

type Params struct {
	Filename string
	Sep      string
}

func (l *LocalWordDB) GetRandomWord() string {
	return l.words[l.rand.Intn(len(l.words))]
}

func NewWordDBFromFile(params Params) (*LocalWordDB, error) {
	content, err := ioutil.ReadFile(params.Filename)
	if err != nil {
		return nil, err
	}
	return &LocalWordDB{
		words: strings.Split(string(content), params.Sep),
		rand:  rand.New(rand.NewSource(time.Now().Unix())),
	}, nil
}
