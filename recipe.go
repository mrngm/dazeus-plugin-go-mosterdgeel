package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"strings"
	//	"net"

	"github.com/dazeus/dazeus-go"
)

type Recipe struct {
	ID      int           `json:"id"`
	Link    string        `json:"link"`
	Slug    string        `json:"slug"`
	Title   RecipeTitle   `json:"title"`
	Content RecipeContent `json:"content"`
}

type RecipeTitle struct {
	Rendered  string `json:"rendered"`
	Protected bool   `json:"protected"`
}

type RecipeContent struct {
	Rendered  string `json:"rendered"`
	Protected bool   `json:"protected"`
}

type RenderedContent struct {
	GeschiktVoor    []string `xml:"ul>li"`
	Lijsten         []Lists  `xml:"div>div>ul"`
	Bereidingswijze []string `xml:"div>ol>li"`
}
type Lists struct {
	Class string   `xml:"class,attr"`
	Elems []string `xml:"li"`
}

func getTitles(recipes []Recipe) string {
	var titles []string
	for _, r := range recipes {
		titles = append(titles, html.UnescapeString(r.Title.Rendered))
	}
	return strings.Join(titles, ", ")
}

func GetPossibleRecipes(ev dazeus.Event) (int, error) {
	rID := 0

	if len(ev.Params) <= 1 {
		return 0, fmt.Errorf("Da's een zoektocht waar ik niet aan wil beginnen.")
	}

	// Volledige zoekstring zit in ev.Params[0], losse woorden in ev.Params[1:]
	res, err := SearchRecipes(ev.Params[0])
	if err != nil {
		return 0, err
	}

	switch len(res) {
	case 0:
		return 0, fmt.Errorf("Geen recepten gevonden :-(")
	case 1:
		return res[0].ID, nil
	default:
		return 0, fmt.Errorf("Meer dan één recept gevonden: %v", getTitles(res))
	}

	return rID, nil
}

func GetRecipe(rID int) (Recipe, error) {
	var r Recipe

	url := fmt.Sprintf("https://mosterdgeel.nl/wp-json/wp/v2/posts/%d", rID)
	resp, err := http.Get(url)
	if err != nil {
		return r, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return r, err
	}

	if err := json.Unmarshal(body, &r); err != nil {
		return r, err
	}

	return r, nil
}

func SearchRecipes(needle string) ([]Recipe, error) {
	var r []Recipe

	url := fmt.Sprintf("https://mosterdgeel.nl/wp-json/wp/v2/posts?search=%s", needle)
	resp, err := http.Get(url)
	if err != nil {
		return r, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return r, err
	}

	if err := json.Unmarshal(body, &r); err != nil {
		return r, err
	}

	return r, nil
}

func (rc *RecipeContent) String() string {
	var r RenderedContent
	data := "<recept>" + rc.Rendered + "</recept>"
	err := xml.Unmarshal([]byte(data), &r)
	if err != nil {
		fmt.Printf("error: %v", err)
		return ""
	}
	return fmt.Sprintf("%+v", r.Bereidingswijze)
}
