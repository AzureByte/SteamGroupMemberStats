package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "html/template"
    "strings"
    "math"
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
    pg := &Page{}
    t, _ := template.ParseFiles("index.html")
    t.Execute(r_writer, pg)
}

func currentlyPlayingHandler (r_writer http.ResponseWriter, req *http.Request) {
    communityid := req.URL.Path[len("/community/"):]

    // Get data from SteamAPI
    memberslistxml := getMembersListXML(communityid)
    memberslist := getSteamIdsFromXmlResponse(memberslistxml)
    SteamIDs := getSteamIdsFromXmlResponse(string(memberslist));

    }



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



// Get group details from XML API
func getMembersListXML(communityId string) string {

    // Check whether the community id is an integer and format the request accordingly.
    commIdAsInt, err := strconv.Atoi(communityId)

    if err != nil && commIdAsInt > 0 {
        resp, err := http.Get("https://steamcommunity.com/gid/" + commIdAsInt + "/memberslistxml/?xml=1") // Note when getting data via gid, despite requesting https, Steam will return data via http
    } else {
        resp, err := http.Get("https://steamcommunity.com/groups/" + communityId + "/memberslistxml/?xml=1")
    }
    defer resp.Body.Close()

    memberslistxmlbytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        panic(fmt.Println("Error :",err))
    }

    return string(memberslistxmlbytes)
}


// Hack to retrieve the steamIds via the XML
func getSteamIdsFromXmlResponse(xml string) []string {

    // Exlude the first chunk before the IDs
    steamIDs := strings.SplitAfter(xml, "<steamID64>")[1:]

    // Replace the xml tags in each of the elements
    for i := 0; i < len(steamIDs); i++ {
        steamIDs[i] = strings.Replace(steamIDs[i], "</steamID64>\r\n", "", -1)
        steamIDs[i] = strings.Replace(steamIDs[i], "<steamID64>", "", -1)
    }

    // Replace the last part of the xml Doc
    lastEle := steamIDs[len(steamIDs)-1]
    steamIDs[len(steamIDs)-1] = strings.SplitN(lastEle, "</members>", 2)[0]

    return steamIDs
}