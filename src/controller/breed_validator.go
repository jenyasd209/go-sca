package controller

import (
	"encoding/json"
	"io"
	"net/http"
)

const catApi = "https://api.thecatapi.com/v1/breeds"

type catApiResp []struct {
	Name string `json:"name"`
}

type breedValidator struct {
	api string

	savedNames map[string]struct{}
}

func NewBreedValidator() BreedValidator {
	return &breedValidator{
		api:        catApi,
		savedNames: make(map[string]struct{}),
	}
}

func (bv *breedValidator) Init() error {
	resp, err := http.Get(bv.api)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	cats := catApiResp{}
	err = json.Unmarshal(body, &cats)
	if err != nil {
		return err
	}

	for _, cat := range cats {
		bv.savedNames[cat.Name] = struct{}{}
	}

	return nil
}

func (bv *breedValidator) Validate(name string) bool {
	_, ok := bv.savedNames[name]
	return ok
}
