package main

import (
	"bytes"
	"fmt"
	"github.com/OpenPrinting/goipp"
	"net/http"
)

func makeAttrCollection(name string,
	member1 goipp.Attribute, members ...goipp.Attribute) goipp.Attribute {

	col := make(goipp.Collection, len(members)+1)
	col[0] = member1
	copy(col[1:], members)

	return goipp.MakeAttribute(name, goipp.TagBeginCollection, col)
}

// Build IPP OpGetPrinterAttributes request
func getPrinterAttributes(uri string) ([]byte, error) {
	m := goipp.NewRequest(goipp.DefaultVersion, goipp.OpGetPrinterAttributes, 1)
	m.Operation.Add(goipp.MakeAttribute("attributes-charset",
		goipp.TagCharset, goipp.String("utf-8")))
	m.Operation.Add(goipp.MakeAttribute("attributes-natural-language",
		goipp.TagLanguage, goipp.String("en-US")))
	m.Operation.Add(goipp.MakeAttribute("printer-uri",
		goipp.TagURI, goipp.String(uri)))
	m.Operation.Add(goipp.MakeAttribute("requested-attributes",
		goipp.TagKeyword, goipp.String("all")))

	return m.EncodeBytes()
}

func getMediaReadyAttr(uri string) (string, error) {
	request, err := getPrinterAttributes(uri)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(uri, goipp.ContentType, bytes.NewBuffer(request))
	if err != nil {
		return "", err
	}

	var respMsg goipp.Message

	err = respMsg.Decode(resp.Body)
	if err != nil {
		return "", err
	}

	// Access the "media-ready" attribute
	for _, group := range respMsg.Groups {
		for _, attr := range group.Attrs {
			if attr.Name == "media-ready" {
				for _, value := range attr.Values {
					v := fmt.Sprintf("%v", value.V)
					return v, nil
				}
			}
		}
	}
	return "", nil
}
