package server

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"testing"

	"go-sca/src/controller"
	"go-sca/src/database"
	"go-sca/src/model"
	"go-sca/src/repos"
	"go-sca/src/server/handlers"
	"gorm.io/gorm"
)

// TODO slightly test with bed tests :D

func TestName(t *testing.T) {
	cat := &model.Mission{
		Model:      gorm.Model{},
		ExecutorID: 0,
		Executor: &model.SpyCat{
			Model:      gorm.Model{},
			Name:       "MissionTestName",
			Breed:      "Abyssinian",
			Experience: 3,
			Salary:     2222,
		},
		Targets: []*model.Target{
			{
				Model:     gorm.Model{},
				Name:      "Target1",
				Country:   "Country1",
				Notes:     "some notes; and some more",
				Completed: false,
				MissionID: 0,
				Mission:   nil,
			},
		},
		Completed: false,
	}

	payload, err := json.Marshal(cat)
	if err != nil {
		t.Fatal(err)
	}
	println(string(payload))
}

func TestServer_SpyCat(t *testing.T) {
	tmpFile, err := ioutil.TempFile(".", "example-*.db")
	if err != nil {
		t.Log("Error generating random name:", err)
		t.FailNow()
	}
	defer os.Remove(tmpFile.Name())

	s, url := initDefaultServer(t, tmpFile.Name())
	defer s.Shutdown()

	client := &http.Client{}

	cat := &model.SpyCat{
		Name:       "TestName",
		Breed:      "Abyssinian",
		Experience: 5,
		Salary:     3000,
	}

	payload, err := json.Marshal(cat)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", url+handlers.CatsRoute, bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	resp, err = client.Get(url + handlers.CatsRoute)
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var cats []*model.SpyCat
	err = json.Unmarshal(body, &cats)
	if err != nil {
		t.Fatal(err)
	}
}

func initDefaultServer(t *testing.T, dbFile string) (*Server, string) {
	t.Helper()

	db, err := database.NewDatabase(dbFile, &gorm.Config{})
	if err != nil {
		panic(err)
	}

	validator := controller.NewBreedValidator()
	err = validator.Init()
	if err != nil {
		panic(err)
	}

	catsRepo := repos.NewGenericRepo[model.SpyCat](db)
	spyCatHandler := handlers.NewCatHandler(controller.NewSpyCatController(catsRepo, validator))

	targetRepo := repos.NewGenericRepo[model.Target](db)
	targetHandler := handlers.NewTargetHandler(controller.NewTargetController(targetRepo))

	missionsRepo := repos.NewGenericRepo[model.Mission](db)
	missionsHandler := handlers.NewMissionHandler(controller.NewMissionController(missionsRepo))

	s := NewServer(":8080", spyCatHandler, missionsHandler, targetHandler)
	go func() {
		err = s.Listen()
		if err != nil {
			t.Logf("Error during listening: %s\n", err)
			t.FailNow()
		}
	}()

	return s, "http://localhost:8080"
}

func generateRandomName() (string, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	randomName := base64.URLEncoding.EncodeToString(randomBytes)
	return randomName, nil
}
