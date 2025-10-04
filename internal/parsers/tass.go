package parser

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type EconomyArticle struct {
	Title string `json:"title"`
	URL   string `json:"url"`
	Text  string `json:"text"`
}

// ParseEconomyWithArticles возвращает список статей с заголовком, URL и текстом
func ParseEconomyWithArticlesFromTass(limit int) ([]EconomyArticle, error) {
	listURL := "https://tass.ru/ekonomika"

	req, err := http.NewRequest(http.MethodGet, listURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "news-parser/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch economy page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("bad status fetching economy page: %s", resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse economy page HTML: %w", err)
	}

	// Найдём ссылки
	var articles []EconomyArticle
	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if !ok {
			return
		}
		if !strings.HasPrefix(href, "/ekonomika/") {
			return
		}
		title := strings.TrimSpace(s.Text())
		if title == "" {
			title = strings.TrimSpace(s.Find("h2, h3, span").Text())
		}
		if len(title) < 3 {
			return
		}
		fullURL := "https://tass.ru" + href
		articles = append(articles, EconomyArticle{
			Title: title,
			URL:   fullURL,
			Text:  "", // пока без текста
		})
	})

	// Удалим дубли
	seen := map[string]bool{}
	var unique []EconomyArticle
	for _, a := range articles {
		if !seen[a.URL] {
			seen[a.URL] = true
			unique = append(unique, a)
		}
	}

	// Ограничим количество (если задано)
	if limit > 0 && len(unique) > limit {
		unique = unique[:limit]
	}

	// Теперь пройдемся по каждому и вытянем текст
	var wg sync.WaitGroup
	for idx, art := range unique {
		wg.Add(1)
		go func() {
			defer wg.Done()
			text, err := fetchArticleText(art.URL)
			if err != nil {
				// если не удалось, просто пропустить текст (оставить пустым) или залогировать
				//fmt.Printf("error fetching article %s: %v\n", art.URL, err)
				return
			}
			unique[idx].Text = text
		}()
	}
	wg.Wait()
	return unique, nil
}

// fetchArticleText загружает страницу статьи и пытается вытащить основной текст
func fetchArticleText(url string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("User-Agent", "news-parser/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("fetch article page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("bad status on article: %s", resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("parse article HTML: %w", err)
	}

	var paragraphs []string

	// Стратегия: искать <p> внутри контейнера статьи.
	// В TASS часто текст внутри <div class="journal_article__text"> или <div class="text-block">
	// Примеры селекторов (надо по ситуации подправить):
	doc.Find("div.text-block p, div.journal_article__text p").Each(func(i int, s *goquery.Selection) {
		para := strings.TrimSpace(s.Text())
		if para != "" {
			paragraphs = append(paragraphs, para)
		}
	})

	// Если ничего не нашлось, fallback: взять все <p>
	if len(paragraphs) == 0 {
		doc.Find("p").Each(func(i int, s *goquery.Selection) {
			para := strings.TrimSpace(s.Text())
			if para != "" {
				paragraphs = append(paragraphs, para)
			}
		})
	}

	text := strings.Join(paragraphs, "\n\n")
	return text, nil
}
