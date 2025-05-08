package entities

type Image struct {
	ID    int
	Title string
	Path  string
	Owner int
	Data  []byte
}

type ImageTransform struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Rotate int    `json:"rotate"`
	Resize int    `json:"resize"`
}
