//package main
//
//import (
//	"bytes"
//	"encoding/json"
//	"fmt"
//	"golang.org/x/text/encoding/charmap"
//	"golang.org/x/text/transform"
//	"io/ioutil"
//	"net/http"
//)
//
//type (
//	translate struct {
//		Text   string `json:"q"`
//		Source string `json:"source"`
//		Target string `json:"target"`
//	}
//	AnswerTranslate struct {
//		Data struct {
//			Translations struct {
//				TranslatedText string `json:"translatedText"`
//			} `json:"translations"`
//		} `json:"data"`
//	}
//	AnswerTranslate2 struct {
//		Data struct {
//			Translations map[string]string `json:"translations"`
//		} `json:"data"`
//	}
//)
//
//const (
//	urlLang      = "https://deep-translate1.p.rapidapi.com/language/translate/v2/languages"
//	urlTranslate = "https://deep-translate1.p.rapidapi.com/language/translate/v2"
//)
//
//func main() {
//	status, ans := GetLang()
//	fmt.Println(status)
//	fmt.Println(string(ans))
//	status, body := Translate(translate{Text: "board", Source: "en", Target: "uk"})
//	if status != http.StatusOK {
//		fmt.Println("main.go -> main -> Translate: err: ", status)
//		return
//	}
//
//	data := AnswerTranslate2{}
//	err := json.Unmarshal(body, &data)
//	if err != nil {
//		fmt.Println("main.go -> Request -> response ", err, "\n body", string(body))
//		return
//	}
//	//dec := charmap.Windows1252.NewDecoder()
//	utfData := body[43:len(body)-4]
//	//reader := dec.Reader(bytes.NewReader(utfData))
//	reader := transform.NewReader(bytes.NewReader(utfData),charmap.Windows1252.NewDecoder())
//	text, _ := ioutil.ReadAll(reader)
//	fmt.Printf("%#v",body[43:len(body)-4])
//	fmt.Println()
//	//fmt.Println(string(text))
//	fmt.Println(string(body[43:len(body)-4]))
//	//fmt.Println(string(body))
//	fmt.Println(string(text))
//	//fmt.Println("data ",data)
//	//fmt.Println(data.Data.Translations.TranslatedText)
//
//}
//
//func GetLang() (int, []byte) {
//	return Request(http.MethodGet, urlLang, []byte{})
//}
//func Translate(translate translate) (int, []byte) {
//
//	data, _ := json.Marshal(translate)
//	return Request(http.MethodPost, urlTranslate, data)
//}
//
//func Request(method, url string, body []byte) (int, []byte) {
//	req, _ := http.NewRequest(method, url, bytes.NewReader(body))
//
//	if method == http.MethodPost {
//		req.Header.Add("content-type", "application/json")
//	}
//	req.Header.Add("x-rapidapi-key", "3cb9a869c6msha09e6ad712fe86dp10a9f0jsnb6dd2d6ae7f6")
//	req.Header.Add("x-rapidapi-host", "deep-translate1.p.rapidapi.com")
//
//	res, _ := http.DefaultClient.Do(req)
//
//	defer res.Body.Close()
//	bodyAns, _ := ioutil.ReadAll(res.Body)
//	if res.StatusCode != http.StatusOK {
//		fmt.Println("main.go -> Request -> response ", res, "\n body", string(body))
//	}
//	return res.StatusCode, bodyAns
//}


