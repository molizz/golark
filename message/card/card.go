package card

type Header struct {
	Title *Tag `json:"title"`
}

type Config struct {
	WideScreenMode bool `json:"wide_screen_mode"`
}

type Link struct {
	URL        string `json:"url,omitempty"`         // url
	AndroidURL string `json:"android_url,omitempty"` // android url
	IOSURL     string `json:"ios_url,omitempty"`     // ios url
	PCURL      string `json:"pc_url,omitempty"`      // pc url
}

type Card struct {
	Header *Header `json:"header,omitempty"`
	Config *Config `json:"config,omitempty"`
	Link   *Link   `json:"card_link,omitempty"`

	Elements []*Tag `json:"elements"`
}

func (c *Card) AddElements(tags ...*Tag) {
	if len(tags) == 0 {
		return
	}
	c.Elements = append(c.Elements, tags...)
}

func NewSimpleCard(title, desc, url string) *Card {
	div := TagDiv(desc)
	div.SetExtra(TagButtonURL("详情", url))

	c := NewCard(title)
	c.AddElements(div)
	return c
}

func NewCard(title string) *Card {
	c := &Card{}
	c.Header = &Header{
		Title: &Tag{
			Tag:     TagLabelPlainText,
			Content: title,
		},
	}
	c.Config = &Config{WideScreenMode: false}
	return c
}
