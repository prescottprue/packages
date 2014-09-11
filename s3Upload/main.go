package s3Upload

import (
	"github.com/kr/s3/s3util"
	"io"
	"log"
	"os"
)
func UploadImg(f io.Reader, l string) (int, string) {
  log.Println("UploadImage run with l:", l)
  s3util.DefaultConfig.AccessKey = os.Getenv("ECHO_DEV_S3_ACCESS_KEY")
  s3util.DefaultConfig.SecretKey = os.Getenv("ECHO_DEV_S3_SECRET_KEY")
  // s3Url := os.Getenv("ECHO_DEV_S3_URL")
  //Create new s3 object
  s3w, _ := s3util.Create(l, nil, nil)
  //Copy file to s3 object
  p, err := io.Copy(s3w, f)
  //Log Upload Size
  log.Println("Image Size: ", p, " bytes")
  //Error/Reponse handling
  if err != nil {
    log.Fatal("Something went wrong")
		return 500, "failure"
  }
  s3w.Close()
  log.Println("UploadImage successful")
  return 200, l
}