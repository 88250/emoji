package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

type EmojiSheet struct {
	ID    string       `json:"id"`
	Title string       `json:"title"`
	Items []*EmojiItem `json:"items"`
}

type EmojiItem struct {
	Unicode         string `json:"unicode"`
	Description     string `json:"description"`
	DescriptionZhCN string `json:"description_zh_cn"`
	Keywords        string `json:"keywords"`
}

func main() {
	data, err := os.ReadFile("emoji.json")
	if nil != err {
		fmt.Println(err)
		return
	}

	var emojis []*EmojiSheet
	if err = json.Unmarshal(data, &emojis); nil != err {
		fmt.Println(err)
		return
	}

	for _, e := range emojis {
		fmt.Println(e.Title)
		for i, item := range e.Items {
			item.DescriptionZhCN, err = googleTranslate(item.Description)
			if nil != err {
				fmt.Printf("translate [%s] failed: %s", item.Description, err)
			}
			keywordsZhCN, err := googleTranslate(item.Keywords)
			if nil != err {
				fmt.Printf("translate [%s] failed: %s", item.Keywords, err)
			}
			item.Keywords += "," + strings.ReplaceAll(keywordsZhCN, "，", ",")
			fmt.Printf("%s|%s|%s|%s\n", item.Unicode, item.Keywords, item.Description, item.DescriptionZhCN)

			//if 5 < i {
			//	break
			//}
			_ = i
		}
	}

	data, err = json.MarshalIndent(emojis, "", "\t")
	if nil != err {
		fmt.Println(err)
		return
	}
	if err = os.WriteFile("final.json", data, 0644); nil != err {
		fmt.Println(err)
		return
	}
}

func googleTranslate(text string) (ret string, err error) {
	//dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:1080", nil, proxy.Direct)
	//if err != nil {
	//	return
	//}
	//
	//httpTransport := &http.Transport{Dial: dialer.Dial}
	//http.DefaultClient.Transport = httpTransport

	ctx := context.Background()
	client, err := translate.NewClient(ctx)
	if err != nil {
		return
	}

	translations, err := client.Translate(ctx, []string{text}, language.Chinese, &translate.Options{Source: language.English, Format: "text", Model: "base"})
	if nil == err {
		ret = translations[0].Text
	}
	return
}
