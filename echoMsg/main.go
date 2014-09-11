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
    // CreatedAt string `json:"createdAt"`
    Recipients []string `json:"recipients"`
    User User `json:"user"`
    Image Image `json:"image"`
    Id string `json:"id"`
}
type Response struct {
	// CreatedAt string `json:"createdAt"`
	User User `json:"user"`
	Image Image `json:"image"`
  Id string `json:"id"`
}
// When using io.Reader (s3 package)
// type ResponseInfo struct {
//   Image io.Reader `json:"image"`
//   User User `json:"user"`
//   Mid string  `json:"mid"`
//   Id string `json:"id"`
// }
type ResponseInfo struct {
  Image Image `json:"image"`
  User User `json:"user"`
  Id string `json:"id"`
}
type BookmarkRequest struct {
  User User `json:"user"`
  Message Message `json:"message"`
}
//Yells
type Yell struct {
    Title string `json:"title"`
    // CreatedAt string `json:"createdAt"`
    User User `json:"user"`
    Image Image `json:"image"`
    Yid string `json:"yid"`
}
type YellInfo struct {
    Title string `json:"title"`
    // CreatedAt string `json:"createdAt"`
    User User `json:"user"`
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
    //[TODO] Return Status Int
  return msg
}
//-----------Echos---------------
func SendMessage(f io.Reader, mid string) int {
  //Get Message/Image Info from firebase
  msg := GetMessage(mid)
  l := msg.Image.Url
  //Upload Image to S3
  us, _ := s3Upload.UploadImg(f, l)
  if us != 200 {
    panic("Error uploading image")
  }
  //Send Message To Recipients
  rs := PushMessageToRecipients(msg)
  if rs != 200 {
    panic("Error Pushing to recipients")
  }
  return rs
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
    parsePush.NotifyUser(uid, "New Echo from " + m.User.DisplayName)
  }
 return 200
}

//Send Response with ri having io.Reader To Author and Recipients
// func SendResponse(ri *ResponseInfo) int {
// 	l := "userData/"+ ri.User.Uid + "/"+ ri.Mid + "/file.jpg"
// 	us, url := s3Upload.UploadImg(ri.Image,l)
//   if us != 200 {
// 		fmt.Println("Error Uploading Image")
// 	}
// 	//Get response authors info from message
// 	var r *Response
// 	r.Image.Url = url
// 	r.User.Uid = ri.User.Uid
// 	r.User.DisplayName = ri.User.DisplayName
// 	r.CreatedAt = "69696969696"
// 	var m *Message
// 	//Get original message object
// 	m = GetMessage(ri.Mid)

//   fmt.Println("SendResponse called for:", r)
//   if as := AuthorSendResponse(m, r); as != 200 {
//   	fmt.Println("Error Sending Response To Author:", as)
//   }
//   if rs := RecipientsSendResponse(m, r); rs != 200 {
//   	fmt.Println("Error Sending Response To Recipients:", rs)
//   }
//   return 200
// }
//If responseInfo(ri) includes img url
func SendResponse(r *Response) int {
 //Get response authors info from message

//Fill with example data for now
 var m *Message
 //Get original message object
 m = GetMessage(r.Id)

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
    rUrl = fbUrl + "/users/" + uid + "/messages/received/" + m.Id + "/response"
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
  var ref *firebase.Reference
    var err error

    //Send To main message
    mUrl := fbUrl + "/messages/" + m.Id + "/responses"
    fmt.Println("mUrl:", mUrl)
    ref = firebase.NewReference(mUrl).Auth(fbSecret).Export(false)
    if err = ref.Push(&r); err != nil {
        panic(err)
    }

    //Send To original author
    aUrl := fbUrl + "/users/" + m.User.Uid + "/messages/sent/" + m.Id + "/responses"
    fmt.Println("aUrl:", aUrl)
    ref = firebase.NewReference(aUrl).Auth(fbSecret).Export(false)
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

func SendBookmark(r *BookmarkRequest) int {
  fmt.Println("SendBookmark Called with", r)

  fbUrl := os.Getenv("ECHO_DEV_FB_URL")
  fbSecret := os.Getenv("ECHO_DEV_FB_SECRET")
  
  var err error

  //Add bookmarker's uid to BookmarkedBy of main message
  mUrl := fbUrl + "/messages/"+ r.Message.Id +"/bookmarkedBy/"+ r.User.Uid
  mesRef := firebase.NewReference(mUrl).Auth(fbSecret).Export(false)
  if err = mesRef.Write(r.User); err != nil {
    panic(err)
  }

  //Add bookmarker's uid to BookmarkedBy of sender's message
  sUrl := fbUrl + "/users/" + r.Message.User.Uid + "/messages/sent/"+ r.Message.Id +"/bookmarkedBy/"+ r.User.Uid
  sRef := firebase.NewReference(sUrl).Auth(fbSecret).Export(false)
  if err = sRef.Write(r.User); err != nil {
    panic(err)
  }

  fmt.Println("Message with id:", r.Message.Id, " has been bookmarked by: ", r.User)
  //notify Author (maybe recipients)
  aMsg := r.User.DisplayName + " bookmarked " + r.Message.Title
  parsePush.NotifyUser(r.Message.User.Uid, aMsg)
  fmt.Println(r.Message.User.DisplayName + " was sent a push notification about the bookmark")

  return 200
}
//---------------Yell Action Functions-----------//
// Handled on client using /api/upload
// func SendYell(f *io.Reader, yd *YellInfo) int {
//   fmt.Println("SendYell Called with", y)
//   fbUrl := os.Getenv("ECHO_DEV_FB_URL")
//   fbSecret := os.Getenv("ECHO_DEV_FB_SECRET")
//   //Create original yell object 
//   oUrl := fbUrl +"/yells"
//   oRef := firebase.NewReference(oUrl).Auth(fbSecret).Export(false)
//   var oerr error
//   if oerr = oRef.Push(yd); err != nil {
//     panic(err)
//   }
//   //update main yell object with image url
//   mUrl := fbUrl + "/yells/"+ y.Id +"/responses"
//   mesRef := firebase.NewReference(mUrl).Auth(fbSecret).Export(false)
//   var err error
//   if err = mesRef.Push(yr); err != nil {
//     panic(err)
//   }
//   fmt.Println("SendYellResponse posted response to main yell object")

// }

// type Response struct {
//   // CreatedAt string `json:"createdAt"`
//   User User `json:"user"`
//   Image Image `json:"image"`
//   Id string `json:"id"`
// }
func SendYellResponse(yr *Response) int {
  fmt.Println("SendYellResponse Called with", yr)
  fbUrl := os.Getenv("ECHO_DEV_FB_URL")
  fbSecret := os.Getenv("ECHO_DEV_FB_SECRET")
  //Post Response to main yell object
  mUrl := fbUrl + "/yells/"+ yr.Id +"/responses"
  mesRef := firebase.NewReference(mUrl).Auth(fbSecret).Export(false)
  var err error
  if err = mesRef.Push(yr); err != nil {
    panic(err)
  }
  fmt.Println("SendYellResponse posted response to main yell object")

  //Update ResponseCounter on authors yell object
  //[TODO] Get authors uid from firebase
  // aUrl := fbUrl + "/users/"+ yr.Id +"/responses"
  // mesRef := firebase.NewReference(mUrl).Auth(fbSecret).Export(false)
  // var err error
  // if err = mesRef.Push(yr); err != nil {
  //   panic(err)
  // }
  // fmt.Println("SendYellResponse updated author's counter")

  //Notify Author Of Response
  //[TODO] Get author uid
  // parsePush.NotifyUser(yr.User.Uid, aMsg)
  return 200
}