package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
	"golang.org/x/net/html"
)

const webDir = "./assets/web"
const varNamePrefix = "asset_"

var exportTo = []string{"src/proxy/web_data.go", "src/daemon/web_data.go"}

func main() {
	wa, _ := listWebAssets()
	out := make([]byte, 0)
	out = append(out, []byte("package main\n\n")...)
	for _, name := range wa {
		if strings.HasSuffix(name, ".html") {
			data, err := webCompile(name)
			if err != nil {
				panic(err)
			}
			varName := strcase.ToLowerCamel(varNamePrefix + strings.TrimSuffix(name, ".html"))
			line := fmt.Sprintf("const %s = \"%s\"\n", varName, data)
			out = append(out, []byte(line)...)
		} else if strings.HasSuffix(name, ".ico") {
			dataReader, err := getWebAsset(name)
			if err != nil {
				panic(err)
			}
			data, err := ioutil.ReadAll(dataReader)
			if err != nil {
				panic(err)
			}
			varName := strcase.ToLowerCamel(varNamePrefix + strings.TrimSuffix(name, ".ico"))
			line := fmt.Sprintf(
				"const %s = \"%s\"\n", varName,
				base64.StdEncoding.EncodeToString(data),
			)
			out = append(out, []byte(line)...)
		}
	}
	for _, exportPath := range exportTo {
		if err := ioutil.WriteFile(exportPath, out, 0755); err != nil {
			panic(err)
		}
	}
}

func listWebAssets() ([]string, error) {
	fis, err := ioutil.ReadDir(webDir)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0)
	for _, fi := range fis {
		if fi.IsDir() {
			continue
		}
		out = append(out, fi.Name())
	}
	return out, nil
}

func webCompile(name string) (string, error) {
	// read asset and parse html
	data, err := getWebAsset(name)
	if err != nil {
		return "", err
	}
	doc, err := html.Parse(data)
	if err != nil {
		return "", err
	}
	var ittDoc func(*html.Node)
	ittDoc = func(n *html.Node) {
		// find all assets and replace with data uris
		if n.Type == html.ElementNode {
			for i := range n.Attr {
				switch n.Attr[i].Key {
				case "href":
					{
						// read asset
						hrefData, err := getWebAsset(strings.Trim(n.Attr[i].Val, "/"))
						if err != nil {
							log.Printf("[WARN] %s", err.Error())
							return
						}
						// get mime type
						mime := ""
						for _, attr := range n.Attr {
							if attr.Key == "type" {
								mime = attr.Val
							}
						}
						// generate data uri and replace href
						dataUri, err := generateDataUri(hrefData, mime)
						if err != nil {
							log.Printf("[WARN] %s", err.Error())
							return
						}
						n.Attr[i].Val = dataUri
						break
					}
				case "src":
					{
						log.Println(n.Data)
						break
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			ittDoc(c)
		}
	}
	ittDoc(doc)
	buf := bytes.Buffer{}
	if err := html.Render(&buf, doc); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func getWebAsset(name string) (io.Reader, error) {
	pathTo := filepath.Join(webDir, name)
	rawData, err := ioutil.ReadFile(pathTo)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(rawData), nil
}

func generateDataUri(r io.Reader, mime string) (string, error) {
	byteData, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}
	return "data:" + mime + ";base64," + base64.StdEncoding.EncodeToString(byteData), nil
}
