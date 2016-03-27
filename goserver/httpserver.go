package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "html/template"
    //"encoding/json"
)

type Page struct {
    Title string
    Body []byte
    Data string
}

type JsonFile struct {
    Name string
    Data []byte
}

func (j *JsonFile) save() error {
    filename := j.Name + ".json"
    return ioutil.WriteFile(filename, j.Data, 0600)
}

func handler(r_writer http.ResponseWriter, req *http.Request) {
    fmt.Fprintf(r_writer, "This should load the home page!", "")
}

func currentlyPlayingHandler (r_writer http.ResponseWriter, req *http.Request) {
    communityid := req.URL.Path[len("/community/"):]

    //Get data from SteamAPI here.


    const data = `{ "topgames" : [
    { "name" : "Dota 2",                           "gameid" : "123321", "numberOfPlayers": 932474},
    { "name" : "Counter-Strike: Global Offensive", "gameid" : "123322", "numberOfPlayers": 297220},
    { "name" : "Team Fortress 2",                  "gameid" : "123323", "numberOfPlayers": 54002},
    { "name" : "Tom Clancy's The Division",        "gameid" : "123324", "numberOfPlayers": 50007},
    { "name" : "Grand Theft Auto V",               "gameid" : "123325", "numberOfPlayers": 42798},
    { "name" : "Warframe",                         "gameid" : "123326", "numberOfPlayers": 31640},
    { "name" : "Sid Meier's Civilization V",       "gameid" : "123327", "numberOfPlayers": 29243},
    { "name" : "Football Manager 2016",            "gameid" : "123328", "numberOfPlayers": 28598},
    { "name" : "Fallout 4",                        "gameid" : "123329", "numberOfPlayers": 27485},
    { "name" : "ARK: Survival Evolved",            "gameid" : "123330", "numberOfPlayers": 27370}]}`;

    fmt.Println(data);

    pg, err := loadPage(communityid)
    if err != nil {
        pg = &Page{Title: communityid, Data: data}
    }
    t, _ := template.ParseFiles("community.html")
    t.Execute(r_writer, pg)
}

func loadPage(title string) (*Page, error) {
    filename := title + ".txt"
    body, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &Page{Title: title, Body: body}, nil
}

func main() {
    http.HandleFunc("/", handler)
    http.HandleFunc("/community/", currentlyPlayingHandler)
    http.ListenAndServe(":8888", nil)
}