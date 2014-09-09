# EchoMsg Package

### Go package to handle message actions for [echo-mobile](http://www.github.com/bessolabs/echo-mobile)'s server [echo-go](http://www.github.com/bessolabs/echo-go)
### Types
User {
    Uid string 
    DisplayName string 
}
Image {
  Url string
}
Message {
    Title string 
    CreatedAt string
    Recipients []string
    User User 
    Image Image 
    Id string
}
Response {
	CreatedAt 
	User User 
	Image Image 
}
ResponseInfo {
  Image io.Reader 
  User User 
  Mid string  
  CreatedAt string 
  Id string
}

### Functions
#### Initial Message
`SendMsg(m *Message) int`

> Sends message to recipients specified in message object
> Calls UpdateImgUrl() then PushMessageToRecipients()

`PushMessageToRecipients(m *Message) int`

> Sends message to all recipients in recipients object.

`UpdateImgUrl(m *Message) int`

> Add Image Url to main message object and author's sent message object (Both created on client)

#### Response

`SendResponse(ri *ResponseInfo) int`

> Send out response to recipients and author
> Calls `RecipientsSendResponse()` and `AuthorSendResponse()`

`RecipientsSendResponse(m *Message, r *Response) int`

> Send response to all recipients (response author to be excluded)

`AuthorSendResponse(m *Message, r *Response) int`

> Send response to original message's author

#### Object Loading
`GetMessage(mid string) *Message`

> Get message object corresponding to message id (mid) provided

### Environment Variables
You must set environment variables listed below:
	
	1. ECHO_DEV_FB_URL --> Firebase Instance Url
  1. ECHO_DEV_FB_SECRET --> Firebase Secret (Admin Panel)
  1. ECHO_DEV_PARSE_ID --> Parse App Id
  1. ECHO_DEV_PARSE_KEY --> Parse Master Key