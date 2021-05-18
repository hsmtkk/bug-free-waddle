package datasrc

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Selector interface {
	SelectDataSrc(html string) ([]string, error)
}

type selectorImpl struct{}

func New() Selector {
	return &selectorImpl{}
}

func (s *selectorImpl) SelectDataSrc(html string) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML; %w", err)
	}
	paths := []string{}
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		dataSrc, ok := s.Attr("data-src")
		if ok {
			paths = append(paths, dataSrc)
		}
	})
	return paths, nil
}
