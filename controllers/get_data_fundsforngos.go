package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/gocolly/colly/v2"
	"github.com/imroc/req/v3"
)

type Article struct {
	Title    string `json:"title"`
	Url      string `json:"url"`
	Shares   int    `json:"shares"`
	Deadline string `json:"deadline"`
	Text     string `json:"text"`
}

func (idb InDb) GetArticleData() {
	result := []map[string]interface{}{}
	columns := []string{"Index", "Title", "Url", "Shares", "Deadline", "Text"}
	colLen := len(columns)
	page := 600

	list_article := get_article_data(page)

	startTime := time.Now()
	xlsx := excelize.NewFile()

	sheetName := "All Student"
	xlsx.SetSheetName(xlsx.GetSheetName(1), sheetName)

	xlsx.SetCellValue(sheetName, "A1", "Funds For NGOs")
	xlsx.MergeCell(sheetName, "A1", "F1")
	xlsx.SetColWidth(sheetName, "A", "F", 20)

	letter := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	idb.setStyle(xlsx, sheetName, []string{
		fmt.Sprintf("%v%v", string(letter[colLen-1]), 1),
		fmt.Sprintf("%v%v", string(letter[colLen-1]), 2),
		fmt.Sprintf("%v%v", string(letter[colLen-1]), len(list_article)+2),
	})

	for key, val := range columns {
		xlsx.SetCellValue(sheetName, fmt.Sprintf("%v%v", string(letter[key]), "2"), val)
	}

	startRow := 3

	for key := range list_article {
		data := map[string]interface{}{
			"Index":    key + 1,
			"Title":    list_article[key].Title,
			"Url":      list_article[key].Url,
			"Shares":   list_article[key].Shares,
			"Deadline": list_article[key].Deadline,
			"Text":     list_article[key].Text,
		}

		xlsx.SetRowHeight(sheetName, startRow, 40)
		xlsx.SetCellValue(sheetName, fmt.Sprintf("%v%v", string(letter[0]), startRow+key), key+1)
		xlsx.SetCellValue(sheetName, fmt.Sprintf("%v%v", string(letter[1]), startRow+key), list_article[key].Title)
		xlsx.SetCellValue(sheetName, fmt.Sprintf("%v%v", string(letter[2]), startRow+key), list_article[key].Url)
		xlsx.SetCellValue(sheetName, fmt.Sprintf("%v%v", string(letter[3]), startRow+key), list_article[key].Shares)
		xlsx.SetCellValue(sheetName, fmt.Sprintf("%v%v", string(letter[4]), startRow+key), list_article[key].Deadline)
		xlsx.SetCellValue(sheetName, fmt.Sprintf("%v%v", string(letter[5]), startRow+key), list_article[key].Text)

		result = append(result, data)
	}

	elapsedTime := time.Since(startTime).Seconds()
	fmt.Printf("\n\nExecution time for app iteration %d: %.2f seconds", len(result), elapsedTime)

	err := xlsx.SaveAs("Funds For NGOs.xlsx")
	if err != nil {
		fmt.Println(err)
	}
}

func get_article_data(page int) []Article {
	result := []Article{}
	var listArticleIndonesia []Article
	var listArticleEducation []Article

	for i := 1; i <= page; i++ {
		list_article, jumlah_article := scrapper_fundsforngos("https://www2.fundsforngos.org/tag/indonesia/page", i)
		if jumlah_article == 0 {
			break
		}

		listArticleIndonesia = append(result, list_article...)
	}

	for i := 1; i <= page; i++ {
		list_article, jumlah_article := scrapper_fundsforngos("https://www2.fundsforngos.org/category/education/page", i)
		if jumlah_article == 0 {
			break
		}

		listArticleEducation = append(result, list_article...)
	}

	var mapEduArticles = make(map[string]*Article)

	listArticleIndonesia = getArticleDetail(listArticleIndonesia)
	listArticleEducation = getArticleDetail(listArticleEducation)

	for key := range listArticleEducation {
		article := &listArticleEducation[key]
		mapEduArticles[article.Title] = article
	}

	for key := range listArticleIndonesia {
		article := &listArticleIndonesia[key]
		included := false

		if mapEduArticles[article.Title] != nil {
			result = append(result, *article)
			included = true
		}

		fmt.Printf("\nArticle %v included : %v\n", article.Title, included)
	}

	return result
}

func scrapper_fundsforngos(url string, page int) (listArticle []Article, jumlah_article int) {
	// Instantiate default collector
	fakeChrome := req.DefaultClient().ImpersonateChrome()

	c := colly.NewCollector(
		colly.UserAgent(fakeChrome.Headers.Get("user-agent")),
	)

	c.SetRequestTimeout(60 * time.Second)
	c.SetClient(&http.Client{
		Transport: fakeChrome.Transport,
	})

	listArticle = []Article{}
	jumlah_article = 0

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Add("Connection", "keep-alive")
		r.Headers.Add("Accept", "*/*")

		fmt.Println("\nUrl : ", r.URL.String())
		fmt.Println("\nVisiting...")
	})

	c.OnHTML("#genesis-content", func(e *colly.HTMLElement) {
		e.ForEach("article", func(i int, wrapper *colly.HTMLElement) {
			fmt.Println("Title : ", wrapper.ChildText("h2.entry-title"))
			article := Article{
				Title: wrapper.ChildText("h2.entry-title"),
				Url:   wrapper.ChildAttr("h2.entry-title a", "href"),
			}

			listArticle = append(listArticle, article)
			jumlah_article += 1
		})
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("\nVisited")
	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("\nGot this error : ", e)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("\nFinished")
	})

	segment := ""
	if page != 0 {
		segment = fmt.Sprint(page)
	}

	c.Visit(fmt.Sprintf("%v/%v", url, segment))

	fmt.Println(len(listArticle))
	return
}

func getArticleDetail(listArticle []Article) []Article {
	for key := range listArticle {
		article := &listArticle[key]

		fakeChrome := req.DefaultClient().ImpersonateChrome()

		newCol := colly.NewCollector(
			colly.UserAgent(fakeChrome.Headers.Get("user-agent")),
		)

		newCol.SetRequestTimeout(60 * time.Second)
		newCol.SetClient(&http.Client{
			Transport: fakeChrome.Transport,
		})

		newCol.OnRequest(func(r *colly.Request) {
			fmt.Println("\nUrl detail page : ", r.URL.String())
			fmt.Println("\nVisiting detail page...")
		})

		newCol.OnHTML("#genesis-content > article:nth-child(2)", func(f *colly.HTMLElement) {
			shares := f.ChildText(".counts")
			if shares != "" {
				article.Shares, _ = strconv.Atoi(shares)
			}

			deadline := f.ChildText(".entry-content > p:nth-child(2)")
			if deadline != "" {
				deadline = strings.ReplaceAll(deadline, "Deadline: ", "")
				article.Deadline = deadline
			}

			f.ForEach(".entry-content", func(i int, entry_content *colly.HTMLElement) {
				entry_content.DOM.Find("aside").Remove()
				entry_content.DOM.Find("p").First().Remove()

				html, err := entry_content.DOM.Html()
				if err != nil {
					fmt.Println("Failed to get text", err.Error())
				}

				if html != "" {
					md := HTMLtoMd(html)

					// replace all new line
					md = strings.ReplaceAll(md, "\n", "\\n")
					article.Text = md
				}
			})
		})

		newCol.OnResponse(func(r *colly.Response) {
			fmt.Println("\nVisited detail page")
		})

		newCol.OnError(func(r *colly.Response, e error) {
			fmt.Println("\nGot this error : ", e)
		})

		newCol.OnScraped(func(r *colly.Response) {
			fmt.Println("\nFinished detail page")
		})

		newCol.Visit(article.Url)
	}

	return listArticle
}

func HTMLtoMd(html string) string {
	converter := md.NewConverter("", true, nil)

	markdown, err := converter.ConvertString(html)
	if err != nil {
		log.Fatal(err)
	}

	return markdown
}
