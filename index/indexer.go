package index

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"log"
	"net/http"
)

type Indexer interface {
	IndexerAction(settings *Configuration, uuid *string) error
}

type AddItem struct {}

type DeleteItem struct {}


func (axn AddItem) IndexerAction(settings *Configuration, uuid *string) error {

	manifestJson, err := getManifest(settings.DSpaceHost, *uuid)
	if err != nil {
		return err
	}
	manifest, err := unMarshallManifest(manifestJson)
	if err != nil {
		return err
	}
	annotations, err := unMarshallAnnotationList(getAnnotationList(manifest.SeeAlso.Id))
	if err != nil {
		return err
	}
	annotationsMap := createAnnotationMap(annotations.Resources)
	altoFiles, err := getAltoFiles(annotationsMap)
	if err != nil {
		return err
	}
	err = indexFiles(*uuid, annotationsMap, altoFiles, manifest.Id, *settings)
	if err != nil {
		return err
	}
	return nil

}

func (axn DeleteItem) IndexerAction(settings *Configuration, uuid *string) error {

	return nil

}

func getAltoFiles(annotationsMap map[string]string) ([]string, error) {
	metsReader, err := getMetsXml(annotationsMap["mets.xml"])
	if err != nil {
		return nil, err
	}
	altoFiles := getOcrFileNames(metsReader)
	return altoFiles, nil
}


func getOcrFileNames(metsReader io.Reader) []string {
	var fileNames = make([]string, 50)
	parser := xml.NewDecoder(metsReader)
	ocrFileElement := false
	altoCounter := 0
	for {
		token, err := parser.Token()
		if err != nil {
			break
		}
		switch t := token.(type) {
		case xml.StartElement:
			element := xml.StartElement(t)
			name := element.Name.Local
			if name == "file" {
				for i := 0; i < len(element.Attr); i++ {
					if element.Attr[i].Value == "ocr" {
						ocrFileElement = true
					}
				}
			}
			if name == "FLocat" && ocrFileElement == true {
				for i := 0; i < len(element.Attr); i++ {
					if element.Attr[i].Name.Local == "href" {
						// Allocate more capacity.
						if altoCounter == cap(fileNames)  {
							newFileNames := make([]string, 2*cap(fileNames))
							copy(newFileNames, fileNames)
							fileNames = newFileNames
						}
						fileNames[altoCounter] = element.Attr[i].Value
						altoCounter++
					}
				}
			}

		case xml.EndElement:
			element := xml.EndElement(t)
			name := element.Name.Local
			if name == "file" {
				ocrFileElement = false
			}
		}
	}

	return fileNames

}

func createAnnotationMap(annotations []ResourceAnnotation) map[string]string {

	annotationMap := make(map[string]string)
	for i := 0; i < len(annotations); i++ {
		value := annotations[i].Resource.Id
		key := annotations[i].Resource.Label
		annotationMap[key] = value
	}
	return annotationMap
}

func unMarshallManifest(bytes []byte) (Manifest, error) {
	var manifest Manifest
	if err := json.Unmarshal(bytes, &manifest); err != nil {
		errorMessage := errors.New("could not unmarshal manifest: " + err.Error())
		return manifest, errorMessage
	}
	return manifest, nil
}

func unMarshallAnnotationList(bytes []byte) (ResourceAnnotationList, error) {
	var annotations ResourceAnnotationList
	if err := json.Unmarshal(bytes, &annotations); err != nil {
		errorMessage := errors.New("could not unmarshal annotationlist: " + err.Error())
		return annotations, errorMessage
	}
	return annotations, nil
}

func getManifest(host string, uuid string) ([]byte, error) {
	endpoint := getApiEndpoint(host, uuid, "manifest")
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, err
}

func getAnnotationList(id string) []byte {
	resp, err := http.Get(id)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	return body

}

func getApiEndpoint(host string, uuid string, iiiftype string) string {
	return host + "/iiif/" + uuid + "/" + iiiftype
}
