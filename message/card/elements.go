package card

import (
	"fmt"
)

const (
	TagLabelHr           = "hr"
	TagLabelDiv          = "div"
	TagLabelPlainText    = "plain_text"
	TagLabelSelectStatic = "select_static"
	TagLabelNote         = "note"
	TagLabelButton       = "button"
	TagLabelLarkMarkdown = "lark_md"
)

const (
	ButtonTypeDefault = "default"
	ButtonTypePrimary = "primary"
	ButtonTypeDanger  = "danger"
)

type Tag struct {
	Tag      string      `json:"tag"`               // tag label
	Content  string      `json:"content,omitempty"` //
	Type     string      `json:"type,omitempty"`
	Text     *Tag        `json:"text,omitempty"`
	Title    *Tag        `json:"title,omitempty"`
	ImgKey   string      `json:"img_key,omitempty"`
	Alt      *Tag        `json:"alt,omitempty"`
	Actions  []*Tag      `json:"actions,omitempty"`
	Extra    *Tag        `json:"extra,omitempty"`
	Options  []*Tag      `json:"options,omitempty"`
	Value    interface{} `json:"value,omitempty"`    // option value
	Elements []*Tag      `json:"elements,omitempty"` //
	URL      string      `json:"url,omitempty"`
}

func (t *Tag) AddElement(tags ...*Tag) *Tag {
	t.Elements = append(t.Elements, tags...)
	return t
}

func (t *Tag) SetExtra(tag *Tag) *Tag {
	t.Extra = tag
	return t
}

func TagDiv(title string) *Tag {
	return &Tag{
		Tag: TagLabelDiv,
		Text: &Tag{
			Tag:     TagLabelPlainText,
			Content: title,
		},
	}
}

func TagHr() *Tag {
	return &Tag{
		Tag: TagLabelHr,
	}
}

func TagImg(imgKey, imgTitle string) *Tag {
	return &Tag{
		Tag: "img",
		Title: &Tag{
			Tag:     TagLabelPlainText,
			Content: imgTitle,
		},
		ImgKey: imgKey,
		Alt: &Tag{
			Tag:     TagLabelPlainText,
			Content: imgTitle,
		},
	}
}

func TagPlainText(content string) *Tag {
	return &Tag{
		Tag:     TagLabelPlainText,
		Content: content,
	}
}

func TagNote(note string) *Tag {
	noteTag := &Tag{
		Tag: TagLabelNote,
	}
	if len(note) > 0 {
		noteTag.AddElement(TagPlainText(note))
	}
	return noteTag
}

func TagButtonURL(title, url string) *Tag {
	return &Tag{
		Tag: TagLabelButton,
		Text: &Tag{
			Tag:     TagLabelPlainText,
			Content: title,
		},
		URL:  url,
		Type: ButtonTypePrimary,
	}
}

func TagURL(title, url string) *Tag {
	return &Tag{
		Tag:     TagLabelLarkMarkdown,
		Content: fmt.Sprintf("[%s](%s)", title, url),
	}
}
