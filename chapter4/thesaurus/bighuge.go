package thesaurus

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type BigHuge struct {
	APIKEY string
}

type synonyms struct {
	Noun *words `json:"noun"`
	Verb *words `json:"verb"`
}

type words struct {
	Syn []string `json:"syn"`
}

func (b *BigHuge) Synonyms(term string) ([]string, error) {
	var syns []string
	resp, err := http.Get("http://words.bighugelabs.com/api/2/" + b.APIKEY + "/" + term + "/json")
	if err != nil {
		return syns, fmt.Errorf("bighuge: %qの類語検索に失敗しました: %v", term, err)
	}
	var data synonyms
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return syns, err
	}
	syns = append(syns, data.Noun.Syn...)
	syns = append(syns, data.Verb.Syn...)
	return syns, nil
}
