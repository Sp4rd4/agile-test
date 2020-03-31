package fuzzyelem

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

const (
	// MaxDepthFromRoot sets how deep will recursive search go
	MaxDepthFromRoot = 4
	// RootDepth decides order of ancestor node to start search from
	RootDepth = 2

	PossibleSimilarityTreshold = 3
	CertainSimilarityTreshold  = 6
)

// searchResult encapsulates data to evaluate nodes against
type searchData struct {
	tag   string
	id    string
	text  string
	class string
	href  string
	title string
}

// searchResult encapsulates node path and its score for a search
type searchResult struct {
	path  string
	score int
}

// Search looks for htlm element in targetPath html file similar to element from sourcePath file.
// Target element from sourcePath is selected by id attribute
func Search(id, sourcePath, targetPath string) (path string, err error) {
	if id == "" {
		return "", errors.New("empty id")
	}

	var sourceFile, targetFile *os.File
	if sourceFile, err = os.Open(sourcePath); err != nil {
		return "", errors.New("source: " + err.Error())
	}
	defer sourceFile.Close()
	if targetFile, err = os.Open(targetPath); err != nil {
		return "", errors.New("target: " + err.Error())
	}
	defer targetFile.Close()

	return search(id, sourceFile, targetFile)
}

// search function prepares searchData and starts recursive search
func search(id string, source, target io.Reader) (path string, err error) {
	doc, err := goquery.NewDocumentFromReader(source)
	if err != nil {
		return "", errors.New("source: " + err.Error())
	}

	docTarget, err := goquery.NewDocumentFromReader(target)
	if err != nil {
		return "", errors.New("target: " + err.Error())
	}

	selection := doc.Find("#" + id)
	if selection.Length() < 1 {
		return "", fmt.Errorf("element with id '%v' not found in source", id)
	}

	class, _ := selection.Attr("class")
	title, _ := selection.Attr("title")
	href, _ := selection.Attr("href")
	sd := &searchData{
		tag:   selection.Get(0).Data,
		id:    id,
		text:  strings.TrimSpace(selection.Text()),
		class: class,
		title: title,
		href:  href,
	}

	out := make(chan searchResult)
	result := searchResult{}
	wg := &sync.WaitGroup{}

	for i := 0; i <= RootDepth; i++ {
		if selection.Parent().Length() > 0 {
			selection = selection.Parent()
		}
	}
	root := docTarget.Find(fullPath(selection, false))
	if root.Length() < 1 {
		return "", nil
	}

	wg.Add(1)
	go inspect(wg, out, 0, root, sd)
	go func() {
		wg.Wait()
		close(out)
	}()

	for res := range out {
		if res.score > result.score {
			result = res
		}
	}
	return result.path, nil
}

// recursive search function
func inspect(wg *sync.WaitGroup, out chan<- searchResult, depth int, sel *goquery.Selection, sd *searchData) {
	defer wg.Done()

	score := score(sel, sd)
	if score > PossibleSimilarityTreshold {
		out <- searchResult{score: score, path: fullPath(sel, true)}
	}
	if score > CertainSimilarityTreshold || depth >= MaxDepthFromRoot {
		return
	}

	sel.Children().Each(func(_ int, sel *goquery.Selection) {
		wg.Add(1)
		go inspect(wg, out, depth+1, sel, sd)
	})
}

func score(sel *goquery.Selection, sd *searchData) int {
	if sel.Get(0).Data != sd.tag {
		return 0
	}
	if v, _ := sel.Attr("id"); v != "" && v == sd.id {
		return 10
	}
	score := 0
	if strings.TrimSpace(sel.Text()) == sd.text {
		score += 5
	}
	if v, _ := sel.Attr("class"); sd.class != "" && v == sd.class {
		score += 3
	}
	if v, _ := sel.Attr("href"); sd.href != "" && v == sd.href {
		score += 1
	}
	if v, _ := sel.Attr("title"); sd.title != "" && v == sd.title {
		score += 1
	}
	return score

}

func elemSelector(s *goquery.Selection, index bool) string {
	path := s.Get(0).Data

	if v, ok := s.Attr("id"); ok {
		path += "#" + v
	}
	if v, ok := s.Attr("class"); ok {
		path += "." + strings.ReplaceAll(v, " ", ".")
	}
	if !index {
		return path
	}
	filtered := s.SiblingsFiltered(path)
	if filtered.Length() > 0 {
		path += "[" + strconv.Itoa(filtered.Index()) + "]"
	}
	return path
}

func fullPath(s *goquery.Selection, indices bool) string {
	path := elemSelector(s, indices)
	for p := s.Parent(); p.Length() > 0; p = p.Parent() {
		path = elemSelector(p, indices) + " > " + path
	}
	return path
}
