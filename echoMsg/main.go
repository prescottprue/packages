package echoMsg

import (
	"github.com/melvinmt/firebase"
	"github.com/bessolabs/packages/parsePush"
	"github.com/bessolabs/packages/s3Upload"
	"log"
	"os"
  "io"
	"fmt"
)
type User struct {
    Uid string `json:"uid"`
    DisplayName string `json:"displayName"`
}
type Image struct {
  Url string `json:"url"`
}
type Message struct {
    Title string `json:"title"`
    CreatedAt string `json:"createdAt"`
    Recipients []string `json:"recipients"`
    User User `json:"user"`
    Image Image `json:"image"`
    Id string `json:"id"`
}
type Response struct {
	CreatedAt string `json:"createdAt"`
	User User `json:"user"`
	Image Image `json:"image"`
}
type ResponseInfo struct {
  Image io.Reader `json:"image"`
  User User `json:"user"`
  Mid string  `json:"mid"`
  CreatedAt string `json:"createdAt"`
  Id string `json:"id"`
}
type BookmarkRequest struct {
  User User `json:"user"`
  Message Message `json:"message"`
}
func SendMsg(m *Message) int {
	log.Println("SendMsg called with", m)
	
	if us := UpdateImgUrl(m); us != 200 {
		fmt.Println("Url Update failed:", us)
	}
	if rs := PushMessageToRecipients(m); rs != 200 {
		fmt.Println("Recipients send failed:", rs)
	}
	return 200
}
func PushMessageToRecipients(m *Message) int {
	log.Println("RecipientsSend called")
	fbUrl := os.Getenv("ECHO_DEV_FB_URL")
  fbSecret := os.Getenv("ECHO_DEV_FB_SECRET")
  //recipient url
  var rUrl string
  var ref *firebase.Reference
  //For loop for recipients
  for ind, uid := range m.Recipients {
    fmt.Println("Recipient:", ind)
    //Send To Each Recipient
    rUrl = fbUrl + "/users/" + uid + "/messages/received/" + m.Id
    fmt.Println("rUrl:", rUrl)

    ref = firebase.NewReference(rUrl).Auth(fbSecret).Export(false)
    var err error
    if err = ref.Write(&m); err != nil {
        panic(err)
    }
    //Notify Recipient
    parsePush.NotifyUser(m.User.Uid, "New Echo from " + m.User.DisplayName)
  }
	return 200
}
func UpdateImgUrl(m *Message) int {
	//Add imgUrl to original msg and author's sent folder
  fmt.Println("UpdateImgUrl Called for", m)
  fbUrl := os.Getenv("ECHO_DEV_FB_URL")
  fbSecret := os.Getenv("ECHO_DEV_FB_SECRET")
  var err error
  

  mUrl := fbUrl + "/messages/" + m.Id
  fmt.Println("mUrl:", mUrl)

  // update image url
  imageRef := firebase.NewReference(mUrl+"/image").Auth(fbSecret).Export(false)
  if err = imageRef.Write(&m.Image); err != nil {
      panic(err)
  }

  // update message id 
  msgRef := firebase.NewReference(mUrl+"/id").Auth(fbSecret).Export(false)
  if err = msgRef.Write(&m.Id); err != nil {
      panic(err)
  }

  //Update Author Message
  aUrl := fbUrl + "/users/" + m.User.Uid + "/messages/sent/"+ m.Id + "/image"
  fmt.Println("aUrl:", aUrl)
  authRef := firebase.NewReference(aUrl).Auth(fbSecret).Export(false)
  var authErr error
  if authErr = authRef.Write(&m.Image); err != nil {
      panic(authErr)
  }

  return 200
}
//Send Response To Author and Recipients
func SendResponse(ri *ResponseInfo) int {
	l := "userData/"+ ri.User.Uid + "/"+ ri.Mid + "/file.jpg"
	us, url := s3Upload.UploadImg(ri.Image,l)
  if us != 200 {
		fmt.Println("Error Uploading Image")
	}
	//Get response authors info from message
	var r *Response
	r.Image.Url = url
	r.User.Uid = ri.User.Uid
	r.User.DisplayName = ri.User.DisplayName
	r.CreatedAt = "69696969696"
	var m *Message
	//Get original message object
	m = GetMessage(ri.Mid)

  fmt.Println("SendResponse called for:", r)
  if as := AuthorSendResponse(m, r); as != 200 {
  	fmt.Println("Error Sending Response To Author:", as)
  }
  if rs := RecipientsSendResponse(m, r); rs != 200 {
  	fmt.Println("Error Sending Response To Recipients:", rs)
  }
  return 200
}
//Send Response to all recipients (including response author's received)
func RecipientsSendResponse(m *Message, r *Response) int {
  fmt.Println("RecipientsSendResponse Called with", r)
	
	fbUrl := os.Getenv("ECHO_DEV_FB_URL")
  fbSecret := os.Getenv("ECHO_DEV_FB_SECRET")
  //recipient url
  var rUrl string
  var res int
  var n string
  var ref *firebase.Reference
  //For loop for recipients
  for ind, uid := range m.Recipients {
    fmt.Println("Recipient:", ind)
    //Send To Each Recipient
    rUrl = fbUrl + "/users/" + uid + "/messages/received/" + m.Id + "/responses"
    fmt.Println("rUrl:", rUrl)

    ref = firebase.NewReference(rUrl).Auth(fbSecret).Export(false)
    var err error
    if err = ref.Push(&r); err != nil {
        panic(err)
    }
    //[TODO] Don't notify author of response
    //Notify Recipient
    n = r.User.DisplayName + " responded to " + m.User.DisplayName + "'s echo"
    if res = parsePush.NotifyUser(uid, n); res != 200 {
    	fmt.Println("Error Notifying Recipient " + uid, res)
    }
  }
	return 200
}
func AuthorSendResponse(m *Message, r *Response) int {
  fmt.Println("AuthorSendResponse Called with", r)
  fbUrl := os.Getenv("ECHO_DEV_FB_URL")
  fbSecret := os.Getenv("ECHO_DEV_FB_SECRET")
  //recipient url
  var aUrl string
  var ref *firebase.Reference
    //Send To Each Recipient
    aUrl = fbUrl + "/users/" + m.User.Uid + "/messages/sent/" + m.Id + "/responses"
    fmt.Println("aUrl:", aUrl)

    ref = firebase.NewReference(aUrl).Auth(fbSecret).Export(false)
    var err error
    if err = ref.Push(&r); err != nil {
        panic(err)
    }
    //Notify Author
    n := r.User.DisplayName + " responded to " + m.Title
    if res := parsePush.NotifyUser(m.User.Uid, n); res != 200 {
    	fmt.Println("Error Notifying Author:", res)
    }
	return 200
}
func GetMessage(mid string) *Message {
	fmt.Println("GetMessage Called with", mid)
  fbUrl := os.Getenv("ECHO_DEV_FB_URL")
  fbSecret := os.Getenv("ECHO_DEV_FB_SECRET")
  //recipient url
  var mUrl string
  var ref *firebase.Reference
  var msg *Message
    //Send To Each Recipient
    mUrl = fbUrl + "/messages/"+ mid
    fmt.Println("mUrl:", mUrl)

    ref = firebase.NewReference(mUrl).Auth(fbSecret).Export(false)
    var err error
    if err = ref.Value(&msg); err != nil {
        panic(err)
    }
  return msg
}
func SendBookmark(r *BookmarkRequest) int {
  fmt.Println("SendBookmark Called with", r)

  fbUrl := os.Getenv("ECHO_DEV_FB_URL")
  fbSecret := os.Getenv("ECHO_DEV_FB_SECRET")
  mUrl := fbUrl + "/messages/"+ r.Message.Id +"/bookmarkedBy/"+ r.User.Uid
  //Add bookmarker's uid to BookmarkedBy of main message
  mesRef := firebase.NewReference(mUrl).Auth(fbSecret).Export(false)
  if err := mesRef.Write(r.User); err != nil {
    panic(err)
  }
  fmt.Println("Message with id:", r.Message.Id, " has been bookmarked by: ", r.User)
  //notify Author (maybe recipients)
  aMsg := r.User.DisplayName + " bookmarked " + r.Message.Title
  parsePush.NotifyUser(r.Message.User.Uid, aMsg)
  fmt.Println(r.Message.User.DisplayName + " was sent a push notification about the bookmark")

  return 200
}