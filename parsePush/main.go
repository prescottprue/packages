package parsePush

import (
  "net/http"
  "fmt"
  "os"
  "github.com/melvinmt/firebase"
  "encoding/json"
  "bytes"
  "io/ioutil"
)

type PushData struct {
  Alert string `json:"alert"`
}
type PushWhere struct {
  Type string `json:"deviceType"`
  Id string `json:"installationId"`
}
type PushMessage struct {
  Data   PushData `json:"data"`
  Where  PushWhere `json:"where"`
}
//Get Push Id from firebase then notify
func NotifyUser(uid string, m string) int {
  fmt.Println("NotifyUser called to notify ", uid, " with ", m);
  //[TODO] Replace this with a local SQL Database 
  //Get pushID from firebase
  fbUrl := os.Getenv("ECHO_DEV_FB_URL")
  fbSecret := os.Getenv("ECHO_DEV_FB_SECRET")
  parId := os.Getenv("ECHO_DEV_PARSE_ID")
  parKey := os.Getenv("ECHO_DEV_PARSE_KEY")

  uUrl := fbUrl + "/users/" + uid +"/pushId"
  fmt.Println("uUrl:", uUrl)
  mainRef := firebase.NewReference(uUrl).Auth(fbSecret).Export(false)
  var err error
  var pid string

  if err = mainRef.Value(&pid); err != nil {
    panic(err)
  }
  fmt.Println("PushId of ", pid, " retreived from firbase for ", uid)
  //Create Push Message with message and pushId
  pd := PushData{m}
  pw := PushWhere{"ios", pid}
  pm := PushMessage{pd, pw}

  jsonMsg, _ := json.Marshal(pm)
  fmt.Printf("Pushing Message: %s\n", jsonMsg)
  contentReader := bytes.NewReader(jsonMsg)
  req, _ := http.NewRequest("POST", "https://api.parse.com/1/push", contentReader)
  req.Header.Set("Content-Type", "application/json")
  req.Header.Set("X-Parse-Application-Id", parId)
  req.Header.Set("X-Parse-REST-API-Key", parKey)
  client := &http.Client{}
  resp, _ := client.Do(req)
  defer resp.Body.Close()
  body, _ := ioutil.ReadAll(resp.Body)
  fmt.Printf("Push completed %s\n", body)
  //body contains {"result": true}
  return resp.StatusCode
}