package model

type Manifest struct {
	Context     string     `json:"@context"`
	Type        string     `json:"@type"`
	Id          string     `json:"@id"`
	Label       string     `json:"label"`
	Metadata    []Metadata `json:"-"`
	Service     Service    `json:"-"`
	SeeAlso     SeeAlso    `json:"SeeAlso,omitempty"`
	Sequences   []Sequence `json:"-"`
	Thumbnail   Thumbnail  `json:"-"`
	ViewingHint string     `json:"-"`
	Related     Related    `json:"-"`
}
type Metadata struct {
	Label string   `json:"label"`
	Value []string `json:"value"`
}

type Service struct {
	Context  string `json:"@context"`
	Id       string `json:"@id"`
	Profile  string `json:"profile"`
	Protocol string `json:"protocol"`
}

type SeeAlso struct {
	Id    string `json:"@id"`
	Type  string `json:"@type"`
	Label string `json:"label"`
}

type Sequence struct {
	Id       string   `json:"@id"`
	Type     string   `json:"@type"`
	Canvases []Canvas `json:"canvases"`
}

type Canvas struct {
	Id        string    `json:"@id"`
	Type      string    `json:"@type"`
	Label     string    `json:"label"`
	Thumbnail Thumbnail `json:"thumbnail"`
	Images    []Image   `json:"images"`
	Width     int       `json:"width,int"`
	Height    int       `json:"height,int"`
}

type Thumbnail struct {
	Id      string  `json:"@id"`
	Type    string  `json:"@type"`
	Label   string  `json:"label"`
	Service Service `json:"service"`
	Format  string  `json:"format"`
}

type Image struct {
	Type       string   `json:"@type"`
	Motivation string   `json:"motivation"`
	Resource   Resource `json:"resource"`
	On         string   `json:"on"`
}

type Resource struct {
	Id      string  `json:"@id"`
	Type    string  `json:"@type"`
	Service Service `json:"service"`
	Format  string  `json:"format"`
}

type Related struct {
	Id     string `json:"@id"`
	Label  string `json:"label"`
	Format string `json:"format"`
}

type ResourceAnnotationList struct {
	Context   string               `json:"@context"`
	Id        string               `json:"@id"`
	Type      string               `json:"@type"`
	Resources []ResourceAnnotation `json:"resources"`
}
type ResourceAnnotation struct {
	Id         string                     `json:"@id"`
	Type       string                     `json:"@type"`
	Motivation string                     `json:"motivation"`
	Resource   ResourceAnnotationResource `json:"resource"`
}
type ResourceAnnotationResource struct {
	Id     string `json:"@id"`
	Label  string `json:"label"`
	Format string `json:"format"`
}
