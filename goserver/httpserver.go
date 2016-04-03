package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "html/template"
    "strings"
    "strconv"
    "math"
    "encoding/json"
    "sort"
)


/**************************

Steam Related

**************************/

var (
    profileStates = map[int]string{
        0: "Offline",
        1: "Online",
        2: "Busy",
        3: "Away",
        4: "Snooze",
        5: "looking to trade",
        6: "looking to play",
    }
    steamAPIKey = ""
)

type JsonGetPlayerSummariesResponse struct {
    Response Response
}

type Response struct {
    Players []MemberInfo
}

type MemberInfo struct {
    Steamid string
    Personaname string
    Personastate int // 0 - Offline, 1 - Online, 2 - Busy, 3 - Away, 4 - Snooze, 5 - looking to trade, 6 - looking to play.
    Gameid string
    Gameextrainfo string
    Communityvisibilitystate int // 1 - the profile is not visible to you (Private, Friends Only, etc), 3 - the profile is "Public"
}

type AllGames struct {
    TopGames []Game
}

type Game struct {
    Name string
    Gameid string
    NumberOfPlayers int
}


/*************************/



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
    gamecounts := make(map[string]int)
    gamenames := make(map[string]string)

    // Get data from SteamAPI
    memberslistxml := getMembersListXML(communityid)
    SteamIDs := getSteamIdsFromXmlResponse(string(memberslistxml));

    NumberOfIDs := len(SteamIDs)
    NumberOfRequests := int(math.Ceil(float64(NumberOfIDs)/100))

    var r JsonGetPlayerSummariesResponse
    // Create a string for each set of max of n members
    n := 100
    i := 0
    for i < NumberOfRequests - 1 {
        s := strings.Join(SteamIDs[i*n : i*n + n], ",")

        r = getPlayerSummariesJSON(s, steamAPIKey)
        for _,v := range r.Response.Players {
            if v.Communityvisibilitystate == 3 && v.Gameid != "" {
                gamecounts[v.Gameid] += 1
                gamenames[v.Gameid] = v.Gameextrainfo
            }
        }

        i++
    }
    //s := strings.Join(SteamIDs[i:], ",")

    // r = getPlayerSummariesJSON(s, steamAPIKey)
    // for _,v := range r.Response.Players {
    //     if v.Communityvisibilitystate == 3 && v.Gameid != "" {
    //         gamecounts[v.Gameid] += 1
    //         gamenames[v.Gameid] = v.Gameextrainfo
    //     }
    // }



    var topg AllGames
    sorted := getSortedKeys(gamecounts)

    for _, v := range sorted[0:int(math.Min(float64(len(sorted)), 10))] {
        var g Game
        g.Name = gamenames[v]
        g.Gameid = v
        g.NumberOfPlayers = gamecounts[v]
        topg.TopGames = append(topg.TopGames, g)
    }

    jbytes, err := json.Marshal(topg)
    if err != nil {
        panic(err)
    }

    data := string(jbytes)

    fmt.Println(data)

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



func getPlayerSummariesJSON(csvSteamIDs string, APIKey string) JsonGetPlayerSummariesResponse {

    var r JsonGetPlayerSummariesResponse
    url := "http://api.steampowered.com/ISteamUser/GetPlayerSummaries/v0002/?key=" + APIKey + "&steamids=" + csvSteamIDs
    fmt.Println(url)
    resp, err_Get := http.Get(url)
    if err_Get != nil {
        panic(err_Get)
    }

    playersummariesjsonbytes, err_ReadAll := ioutil.ReadAll(resp.Body)
    if err_ReadAll != nil {
        panic(err_ReadAll)
    }

    err_Unmarshal := json.Unmarshal(playersummariesjsonbytes,&r)
    if err_Unmarshal != nil {
        panic(err_Unmarshal)
    }

    return r
}

// Get group details from XML API
func getMembersListXML(communityId string) string {

    // Check whether the community id is an integer and format the request accordingly.
    commIdAsInt, err := strconv.Atoi(communityId)
    var url string

    if err != nil && commIdAsInt > 0 {
        url = "https://steamcommunity.com/gid/" + communityId + "/memberslistxml/?xml=1" // Note when getting data via gid, despite requesting https, Steam will return data via http
    } else {
        url = "https://steamcommunity.com/groups/" + communityId + "/memberslistxml/?xml=1"
    }

    resp, err := http.Get(url)
    defer resp.Body.Close()

    memberslistxmlbytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        panic(err)
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

/**************************

Sorting Logic

**************************/


// Using the logic found here https://gist.github.com/ikbear/4038654
type sortedMap struct {
    orig map[string]int //original
    srt []string //sorted
}
func (sm *sortedMap) Len() int {
    return len(sm.orig)
}
func (sm *sortedMap) Less(i, j int) bool {
    return sm.orig[sm.srt[i]] > sm.orig[sm.srt[j]]
}
func (sm *sortedMap) Swap(i, j int) {
    sm.srt[i], sm.srt[j] = sm.srt[j], sm.srt[i]
}
func getSortedKeys(unsortedMap map[string]int) []string {
    sm := new(sortedMap)
    sm.orig = unsortedMap
    sm.srt = make([]string, len(unsortedMap))
    i := 0
    for key, _ := range unsortedMap {
        sm.srt[i] = key
        i++
    }
    sort.Sort(sm)
    return sm.srt
}
/*************************/
