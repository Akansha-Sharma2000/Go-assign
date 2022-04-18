package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/gorilla/mux"
)

type Characters struct{
	Name string `json:"name"`
	MaxPower int `json:"max_power"`
}

type Avengers struct{
	Name string `json:"name"`
	Character []Characters
}

  
type netList struct{
  Name string `json:"name"`
  MaxPower int `json:"max_power"`
  Count int
}

var finalList []netList
var flagCheck =0

func handleRequests() {
  myRouter := mux.NewRouter().StrictSlash(true)
  myRouter.HandleFunc("/", getCharactersArray).Methods("GET")
  myRouter.HandleFunc("/{name}", userCharacter).Methods("GET")
  myRouter.HandleFunc("/post",addCharacters).Methods("POST")
  log.Fatal(http.ListenAndServe("localhost:8080", myRouter))
}

//Getting all the characters info from the urls using GET method
func getCharacters() []netList{
    if flagCheck==0 {
      urls := []string{
        "http://www.mocky.io/v2/5ecfd5dc3200006200e3d64b",
        "http://www.mocky.io/v2/5ecfd630320000f1aee3d64d",
        "http://www.mocky.io/v2/5ecfd6473200009dc1e3d64e",
      }
      for i:=0; i<len(urls);i++ {
        response, err := http.Get(urls[i])
        if err != nil {
            fmt.Print(err.Error())
            os.Exit(1)
        }
    
        responseData, err := ioutil.ReadAll(response.Body)
        if err != nil {
            log.Fatal(err)
        }
    
        var responseObject Avengers
        json.Unmarshal(responseData, &responseObject)
        
        for i := 0; i < len(responseObject.Character); i++ {
          results := []netList{{Name: strings.ToLower(responseObject.Character[i].Name), MaxPower: responseObject.Character[i].MaxPower}}
          finalList=append(finalList, results...)
        }
      }
      flagCheck=1
    } 
    return finalList
}

//For displaying the information in 8080 port
func getCharactersArray(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Content-Type","application/json")
  json.NewEncoder(w).Encode(finalList)
}

//Removing extra characters from the list
func removeCharacters(finalList []netList) []netList{
  if len(finalList)>10 {
        for len(finalList)!=10 {
            finalList=finalList[:len(finalList)-1]
        }
    }
    return finalList
}

//Sorting
func sorting(finalList []netList){
  sort.Slice(finalList, func(i, j int) bool {
    if finalList[i].Count > finalList[j].Count {
        return true
    }
    if finalList[i].Count < finalList[j].Count {
        return false
    }
    return finalList[i].MaxPower > finalList[j].MaxPower
})
}

//Getting specific characters MaxPower by it's name
func userCharacter(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    var flag=0
    key := vars["name"]
    key = strings.ToLower(key)
    finalList := getCharacters()
    for i:=0;i<len(finalList);i++{
        if key==finalList[i].Name {
            flag=1
            jsonString, _ := json.Marshal(finalList[i].MaxPower)
            finalList[i].Count+=1
            fmt.Fprintln(w, string(jsonString))
            break
        }
    }
    if(flag==0){
        fmt.Fprintf(w, "Not available")
    }
    sorting(finalList)
    finalList=removeCharacters(finalList)
}

func addCharacters(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type","application/json")
  var newCharacter netList
  _=json.NewDecoder(r.Body).Decode(&newCharacter)
  finalList=append(finalList, newCharacter)
  json.NewEncoder(w).Encode(newCharacter)
  sorting(finalList)
  finalList=removeCharacters(finalList)
}

func main() {
  finalList=getCharacters()
  handleRequests()
}




