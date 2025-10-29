package models

type Item struct {
	ItemTitle string `json:"title"`
	URL       string `json:"url"`
}

func (i Item) FilterValue() string {
	return i.ItemTitle
}

func (i Item) Title() string {
	return i.ItemTitle
}

func (i Item) Description() string {
	return i.URL
}
