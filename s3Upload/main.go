package s3Upload

import (
	"github.com/kr/s3/s3util"
	"io"
	"log"
	"os"
)
func UploadImg(f io.Reader, l string) (int, string) {
  log.Println(" \033[42m INIT [s3Upload] UploadImage args[l:", l, "] \033[0m ")
  s3util.DefaultConfig.AccessKey = os.Getenv("ECHO_DEV_S3_ACCESS_KEY")
  s3util.DefaultConfig.SecretKey = os.Getenv("ECHO_DEV_S3_SECRET_KEY")
  // s3Url := os.Getenv("ECHO_DEV_S3_URL")
  //Create new s3 object
  s3w, _ := s3util.Create(l, nil, nil)
  //Copy file to s3 object
  _, err := io.Copy(s3w, f)
  //Log Upload Size
  // fmt.Println("Image Size: ", p, " bytes")
  //Error/Reponse handling
  if err != nil {
    panic(err)
		return 500, "failure"
  }
  s3w.Close()
  log.Println(" \033[41m RETURN [s3Upload] UploadImage successful \033[0m ")
  return 200, l
}