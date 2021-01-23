package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"sync"

	"github.com/bushiyama/ghe-api/pkg/env"

	cli "github.com/urfave/cli/v2"
)

const (
	limitPage = 10
)

type obj struct {
	Name   string
	Events int
}
type byEvents []obj

func (a byEvents) Len() int           { return len(a) }
func (a byEvents) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byEvents) Less(i, j int) bool { return a[i].Events > a[j].Events }

var (
	envVars *env.Variables
)

func main() {
	app := &cli.App{
		Action: func(c *cli.Context) error {
			// load vars
			var err error
			envVars, err = env.LoadVariables()
			if err != nil {
				log.Println("[ERROR]load env vars error: " + err.Error())
				return err
			}

			// とりあえず10P全部 新しい順なので時系列見たい時は注意
			wg := new(sync.WaitGroup)
			resChan := make(chan []map[string]interface{}, limitPage)
			for i := 1; i <= limitPage; i++ {
				wg.Add(1)
				go getEvents(wg, resChan, i)
			}
			wg.Wait()
			close(resChan)

			// countのためmapへ
			analyzeMap := map[string]int{}
			for v := range resChan {
				for _, v := range v {
					n := v["actor"].(map[string]interface{})["login"].(string)
					analyzeMap[n] = analyzeMap[n] + 1
				}
			}
			// sorting
			objSlices := make([]obj, len(analyzeMap))
			cnt := 0
			for k, v := range analyzeMap {
				objSlices[cnt] = obj{Name: k, Events: v}
				cnt++
			}
			sort.Sort(byEvents(objSlices))

			// print top 3
			fmt.Println(envVars.Domain + "/" + envVars.User + "/" + envVars.Repo + " Autor of Latest300events TOP3!")
			for i, v := range objSlices {
				if i < 3 {
					fmt.Println(v)
				}
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func getEvents(wg *sync.WaitGroup, resChan chan []map[string]interface{}, p int) {
	defer wg.Done()

	url := "https://" + envVars.Domain + "/api/v3/repos/" + envVars.User + "/" + envVars.Repo + "/events"

	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("Authorization", "token "+envVars.Token)
	params := request.URL.Query()
	params.Add("page", strconv.Itoa(p))
	request.URL.RawQuery = params.Encode()

	client := new(http.Client)
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	var body []map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		log.Fatal(err)
	}
	resChan <- body
}
