package main

import (
	"bytes"
	"fmt"
	"github.com/russross/blackfriday/v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

const rootDirectory = "content"
const publicDirectory = "docs"
const templateDirectory = "templates"

type fileinfo struct {
	filename string
	title    string
}

// tag name => [] of filenames
var tags = make(map[string][]fileinfo)
var renderer = blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{})
var template []byte

func init() {
	tmp, err := ioutil.ReadFile(templateDirectory + "/page.html")
	failOnError(err, "unable to read file template")
	template = tmp
}

func main() {
	start := time.Now()
	log.Printf("--- Starting conversion")
	dir, err := ioutil.ReadDir(rootDirectory)
	failOnError(err, "unable to read content directory: "+rootDirectory)

	var wg sync.WaitGroup
	wg.Add(len(dir))
	for _, file := range dir {
		go func(file os.FileInfo) {
			process(file.Name())
			wg.Done()
		}(file)
	}
	wg.Wait()

	wg.Add(len(tags))
	for tag, files := range tags {
		go func(tag string, files []fileinfo) {
			processTag(tag, files)
			wg.Done()
		}(tag, files)
	}
	wg.Wait()
	log.Printf("--- Finished after %.1gs", time.Now().Sub(start).Seconds())
}

func processTag(tag string, files []fileinfo) {
	var sb strings.Builder

	title := "Articles tagged " + tag
	sb.WriteString("<h1>" + title + "</h1>")
	sb.WriteString(`<ul>`)
	for _, file := range files {
		sb.WriteString(fmt.Sprintf(`<li><a href="%s">%s</a></li>`, file.filename, file.title))
		sb.WriteString("\n")
	}
	sb.WriteString(`</ul>`)

	bs := bytes.ReplaceAll(template, []byte("${content}"), []byte(sb.String()))
	bs = bytes.ReplaceAll(bs, []byte("${title}"), []byte(title))

	tagFilename := "tag-" + tag + ".html"
	ioutil.WriteFile(publicDirectory+"/"+tagFilename, bs, 0644)
}

func process(filename string) {
	log.Printf("Processing <%v>", filename)
	bs, err := ioutil.ReadFile(rootDirectory + "/" + filename)
	failOnError(err, "unable to read filename:"+filename)

	var obs []byte
	var ofilename string

	switch filepath.Ext(filename) {
	case ".md":
		bs = convertWikiLinks(bs)
		bs = convertImages(bs)
		bs = convertTags(bs)
		bs = convertDate(bs)
		mbs := blackfriday.Run(bs, blackfriday.WithRenderer(renderer))
		obs = bytes.ReplaceAll(template, []byte("${content}"), mbs)
		obs = bytes.ReplaceAll(obs, []byte("${title}"), []byte(getTitle(bs)))
		ofilename = strings.ReplaceAll(filename, ".md", ".html")
		ofilename = strings.ReplaceAll(ofilename, " ", "-")
		updateTags(bs, ofilename)
	default:
		ofilename = filename
		obs = bs
	}

	log.Printf("Writing output file <%v>", ofilename)
	ioutil.WriteFile(publicDirectory+"/"+ofilename, obs, 0644)
}

func convertDate(bs []byte) []byte {
	// 202010151100
	regex := regexp.MustCompile(`\d{12}`)
	submatches := regex.FindAllStringSubmatch(string(bs), -1)
	for _, matches := range submatches {
		date := matches[0]
		html := fmt.Sprintf(`<span class="date pull-right">%s-%s-%s %s:%s</span>`, date[:4], date[4:6], date[6:8], date[8:10], date[10:12])
		bs = bytes.ReplaceAll(bs, []byte(date), []byte(html))
	}

	return bs
}

func updateTags(bs []byte, filename string) {
	title := getTitle(bs)
	linkTags := getTags(bs)
	for _, tag := range linkTags {
		tag = tag[1:]
		tags[tag] = append(tags[tag], fileinfo{
			filename: filename,
			title:    title,
		})
	}
}

func getTitle(bs []byte) string {
	if len(bs) == 0 {
		log.Printf("!!! File empty")
		return "Hello, world"
	}
	lines := bytes.Split(bs, []byte("\n"))
	line := string(lines[0])
	if strings.HasPrefix(line, "#") {
		return line[2:]
	} else {
		return line
	}
}

func convertImages(markdown []byte) []byte {
	regex := regexp.MustCompile(`!\[\]\((.*?) (.*?)\)`)
	submatches := regex.FindAllStringSubmatch(string(markdown), -1)
	for _, matches := range submatches {
		text := matches[0]
		image := matches[1]
		width := matches[2]
		html := fmt.Sprintf(`<img src="%s" width="%s"/>`, image, width)
		markdown = bytes.ReplaceAll(markdown, []byte(text), []byte(html))
	}

	return markdown
}

func convertWikiLinks(markdown []byte) []byte {
	regex := regexp.MustCompile(`\[\[(.*?)\]\]`)
	submatches := regex.FindAllStringSubmatch(string(markdown), -1)

	for _, matches := range submatches {
		if len(matches) < 2 {
			continue
		}
		fileLinkName := matches[1]
		displayName := fileLinkName
		wikiLink := matches[0]
		if !strings.Contains(fileLinkName, ".") {
			fileLinkName = fileLinkName + ".html"
		}
		fileLinkName = strings.ReplaceAll(fileLinkName, " ", "-")
		fileLinkName = strings.ToLower(fileLinkName)
		markdownLink := fmt.Sprintf(`[%s](/%s)`, displayName, fileLinkName)
		markdown = bytes.ReplaceAll(markdown, []byte(wikiLink), []byte(markdownLink))
	}

	return markdown
}

func convertTags(markdown []byte) []byte {
	linkTags := getTags(markdown)

	for _, tag := range linkTags {
		link := fmt.Sprintf(`<a href="/tag-%s.html" class="tag">%s </a>`, tag[1:], tag)
		markdown = bytes.ReplaceAll(markdown, []byte(tag), []byte(link))
	}

	return markdown
}

func getTags(data []byte) []string {
	markdown := string(data)
	regex := regexp.MustCompile(` *(#\w+)`)
	matches := regex.FindAllString(markdown, -1)

	tags := []string{}

	for _, tag := range matches {
		if tag != "" {
			tag = strings.Trim(tag, " \n\r\t")
			tags = append(tags, tag)
		}
	}

	return tags
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
