# Parse Push Notifications

### Go package for simplified targeting of push notifications using [Parse](http://www.parse.com) and [Firebase](http://www.firebase.com)

### Functions
`NotifyId(installationId string, message string)`

* Sends message directly to device matching the installation id provided

`NotifyUser(uid string, message string)`
	
* Looks up "pushId" parameter from Firebase user corresponding to the UID provided
* Pushes to the looked up installation Id

### Environment Variables
You must set environment variables listed below:
	1. ECHO_DEV_FB_URL --> Firebase Instance Url
  1. ECHO_DEV_FB_SECRET --> Firebase Secret (Admin Panel)
  1. ECHO_DEV_PARSE_ID --> Parse App Id
  1. ECHO_DEV_PARSE_KEY --> Parse Master Key

#### Note:
##### Notify User Function assumes that your firebase follows this common layout:
	
	Users:{
		//Using Google
		google:123098234098:{
			pushId: "999fffff-999f-9999-99f9-f99999f99999" --> Parse InstallationId
			//Other User Data
		},
		//Using Simple Login
		simplelogin:1:{
			pushId: "999fffff-999f-9999-99f9-f99999f99999" --> Parse InstallationId
			//Other User Data
		},
		//Using Facebook
		facebook:12390182390:{
			pushId: "999fffff-999f-9999-99f9-f99999f99999" --> Parse InstallationId
			//Other User Data
		}
	}
