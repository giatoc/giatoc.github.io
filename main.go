package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type Person struct {
	ID   string `json:"id"`
	Data struct {
		Gender    string `json:"gender"`
		FirstName string `json:"first name"`
		LastName  string `json:"last name"`
		Birthday  string `json:"birthday"`
		Avatar    string `json:"avatar"`
	} `json:"data"`
	Rels struct {
		Children []string `json:"children"`
		Spouses  []string `json:"spouses"`
		Father   string   `json:"father"`
		Mother   string   `json:"mother"`
	} `json:"rels"`
}

func main() {
	jsonData, _ := os.ReadFile("src/components/data.json")
	persons := []Person{}
	json.Unmarshal([]byte(jsonData), &persons)

	mID := map[string]string{}
	mRelID := map[string]string{}
	changeID := false
	for _, p := range persons {
		newID := p.ID
		if changeID {
			newID = convertID(p.Data.FirstName + p.Data.LastName)
			if newID == "" || newID == "-" {
				newID = p.ID
			}
			fmt.Printf("%s -> %s\n", p.ID, newID)

			c := 0
			sf := ""
			for {
				if _, ok := mID[newID+sf]; !ok {
					mID[p.ID] = newID + sf
					break
				} else {
					c += 1
					sf = strconv.Itoa(c)
				}
			}
		}

		mID[p.ID] = newID
		mRelID[p.Rels.Father] = p.Rels.Father
		mRelID[p.Rels.Mother] = p.Rels.Mother
		for _, s := range p.Rels.Spouses {
			mRelID[s] = s
		}
		for _, c := range p.Rels.Children {
			mRelID[c] = c
		}
	}

	if len(mID) != len(mRelID) {
		for k := range mID {
			if _, ok := mRelID[k]; !ok {
				fmt.Println("Missing ID", k)
			}
		}
		for k := range mRelID {
			if _, ok := mID[k]; !ok {
				fmt.Println("Missing Rel ID", k)
			}
		}
	}

	for i := range persons {
		persons[i].ID = mID[persons[i].ID]
		persons[i].Rels.Father = mID[persons[i].Rels.Father]
		persons[i].Rels.Mother = mID[persons[i].Rels.Mother]

		for j := range persons[i].Rels.Spouses {
			persons[i].Rels.Spouses[j] = mID[persons[i].Rels.Spouses[j]]
		}

		for j := range persons[i].Rels.Children {
			persons[i].Rels.Children[j] = mID[persons[i].Rels.Children[j]]
		}
	}

	if changeID {
		sort.Slice(persons, func(i, j int) bool {
			if persons[i].ID == "luuvoduc" {
				return true
			}
			return persons[i].Data.LastName > persons[j].Data.LastName
		})
		jsonData, _ = json.Marshal(persons)
		os.WriteFile("src/components/data.json", jsonData, 0644)
	}

}

func convertID(name string) string {
	name = strings.TrimSpace(name)

	rMapping := func(r rune) rune {
		sortedSpecialRunes := []rune{'Đ', 'đ', 'Ł'}
		replacedByRunes := []rune{'D', 'd', 'L'}

		pos := sort.Search(len(sortedSpecialRunes), func(i int) bool { return sortedSpecialRunes[i] >= r })
		if pos != -1 && r == sortedSpecialRunes[pos] {
			return replacedByRunes[pos]
		}
		return r
	}

	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC, runes.Map(rMapping))
	strTransform, _, _ := transform.String(t, name)
	return strings.ToLower(strings.Replace(strTransform, " ", "", -1))
}
