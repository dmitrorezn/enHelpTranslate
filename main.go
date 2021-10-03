package main

import (
	"encoding/json"
	"os"
	"time"

	"bufio"
	"fmt"
	xls "github.com/xuri/excelize"
	"io/ioutil"
	"net/http"
	"strings"

	"gopkg.in/mgo.v2"
)

var (
	coll = getSession().DB("counter").C("counterRequestTranslate")
	countV2 = getCountV2()
	privateKey = getKey()
)

func getSession() *mgo.Session  {
	db, err := mgo.Dial("mongodb://localhost:27017")
	if err !=nil {
		fmt.Println("dial ",err)
		return nil
	}
	return db
}
func getCountV2() int  {
	c := struct {
		Id uint64 `bson:"_id"`
		Count int `bson:"count"`
	}{}
	if err := coll.FindId(2).One(&c); err != nil {
		fmt.Println("getCount err",err)
	}
	return c.Count
}
func save()  {
	c := struct {
		Count int `bson:"count"`
	}{Count: countV2}

	coll.UpdateId(2,c)
}

func ReadExcel(name string) {
	f, err := xls.OpenFile(name)
	if err != nil {
		fmt.Println("OpenFile", err, name)
		return
	}
	rows, err := f.Rows("Entries")
	if err != nil {
		fmt.Println("Rows", err, name)
		return
	}

	row := 0
	//ct := 0
	for rows.Next() {

		columns, err := rows.Columns()
		if err != nil {
			fmt.Println("Columns", err, name)
			return
		}
		if len(columns) == 0 {
			break
		}
		if row >= 1 && row <= 756 {
			//if strings.Contains(columns[4],",") {
				//ct++
				//fmt.Println(" columns[4]", columns[4], "count ", ct)
			//}
			//fmt.Printf("colums 1 4 %s %s\n", columns[0], columns[4])
			translated, errStr := TranslateV2(columns[0])
			if errStr != "" {
				// Save spreadsheet by the given path.
				if err := f.SaveAs(name); err != nil {
					fmt.Println("SaveAs err =", err)
				}
				fmt.Println(" shotdowt error", errStr)
				return
			}
			fmt.Println("Translate", translated)
			fmt.Println()
			if err := f.SetCellStr("Entries", "E"+fmt.Sprint(row+1), translated); err != nil {
				fmt.Printf("SetCellStr: err %s, row %d \n", err, row)
			}
			time.Sleep(time.Second * 3/2)
		}
		row++
	}
	// Save spreadsheet by the given path.
	if err := f.SaveAs(name); err != nil {
		fmt.Println("SaveAs err =", err)
	}

	if err := f.Close(); err != nil {
		fmt.Println("Close err =", err)
	}
}
func incRequestsV2(){
	countV2++
}


func main() {
	fmt.Println("count start",countV2)
	//w, err := TranslateV2("car rental")
	//fmt.Println(w, err)
	//ReadExcel("translatefile.xlsx")
	fmt.Println("count end",countV2)
	save()
	return
	//fmt.Println("lol", lol)
	//file, _ := os.Create("lol.txt")
	//writtern, err := file.Write([]byte(lol))
	//fmt.Println("n ",writtern ,"err", err)
	//
	//file.Close()
	//answerfile := ""
	////var err error
	//file, err = os.Open("lol.txt")
	//if err != nil {
	//	fmt.Println("open", err)
	//}
	//defer file.Close()
	//scanner := bufio.NewScanner(file)
	//scanner.Split(bufio.ScanWords)
	//for scanner.Scan() {
	//	fmt.Println(scanner.Text())
	//	answerfile = scanner.Text()
	//}
	//if err := scanner.Err(); err != nil {
	//	fmt.Println("Err",err)
	//}
	//fmt.Println("n ",answerfile)

}

func Translate(word string) (string, error) {
	url := "https://google-translate1.p.rapidapi.com/language/translate/v2"

	payload := strings.NewReader(fmt.Sprintf("q=%s&target=ru&source=en", word))

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("accept-encoding", "application/gzip")
	req.Header.Add("x-rapidapi-key", "")
	req.Header.Add("x-rapidapi-host", "google-translate1.p.rapidapi.com")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode != 200 {
		return "", fmt.Errorf("erroor translate")
	}

	return string(body[44 : len(body)-5]), nil
}

func TranslateV1(word string) (string, string) {
	url := "https://deep-translate1.p.rapidapi.com/language/translate/v2"

	payload := strings.NewReader(fmt.Sprintln("{\r\"q\": \"+" + word + "\",\r\"source\": \"en\",\r\"target\": \"uk\"\r}"))

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/json")
	req.Header.Add("x-rapidapi-key", "")
	req.Header.Add("x-rapidapi-host", "deep-translate1.p.rapidapi.com")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
	if res.StatusCode != 200 {
		return "", string(body)
	}

	t := struct {
		Data struct {
			Tr struct {
				Trt string `json:"translatedText"`
			} `json:"translations"`
		} `json:"data"`
	}{}
	err := json.Unmarshal(body, &t)
	if err != nil {
		fmt.Println("TranslateV1 err", err)
		return t.Data.Tr.Trt, err.Error()
	}
	//fmt.Printf("body %s\n", body)
	t.Data.Tr.Trt =strings.Trim(strings.Trim(t.Data.Tr.Trt,"+"),"+ ")
	//fmt.Printf("str %s\n", t.Data.Tr.Trt)
	//f, err := os.OpenFile("l.txt", os.O_RDWR|os.O_CREATE, 0755)
	//f.WriteString(t.Data.Tr.Trt)
	//f.Close()
	return t.Data.Tr.Trt, ""
}

func Answer(code int,body []byte)  (string, string) {
	fmt.Println(string(body))
	if code != 200 {
		return "", string(body)
	}

	a := struct {
		TranslatedText string `json:"translated_text"`
	}{}
	err := json.Unmarshal(body, &a)
	if err != nil {
		fmt.Println("TranslateV1 err", err)
		return a.TranslatedText, err.Error()
	}
	//fmt.Printf("body %s\n", body)
	a.TranslatedText =strings.Trim(strings.Trim(a.TranslatedText,"+"),"+ ")
	fmt.Printf("str %s\n", a.TranslatedText)
	f, err := os.OpenFile("l.txt", os.O_RDWR|os.O_CREATE, 0755)
	f.WriteString(a.TranslatedText)
	f.Close()
	return a.TranslatedText, ""
}


func TranslateV2(word string) (string, string) {
	if countV2 > 1000 {
		return "", "too much requests for today"
	}
	incRequestsV2()
	fmt.Println(strings.Replace(word," ","%20",-1))
	w := strings.Replace(word," ","%20",-1)
	url := "https://translo.p.rapidapi.com/translate?text="+w+"&from=en&to=uk"

	payload := strings.NewReader("{\r\n    \"text\": \""+word+"\",\r\n    \"to\": \"ja\"\r\n}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/json")
	req.Header.Add("x-rapidapi-key", privateKey)
	req.Header.Add("x-rapidapi-host", "translo.p.rapidapi.com")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	//fmt.Println(res)
	fmt.Println(string(body))
	fmt.Println(" total requests counts TranslateV2", countV2)
	return Answer(res.StatusCode,body)
}

func getKey() string {
	f, _:= os.OpenFile("l.txt",os.O_RDWR|os.O_CREATE, 0755)
	b := bufio.NewReader(f)
	key, _ := b.ReadString(byte('\n'))
	fmt.Println(key)
	return key
}
