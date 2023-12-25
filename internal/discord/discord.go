package discord

type Thumbnail struct {
	URL string `json:"url"`
}

type Field struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Embed struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Thumbnail   Thumbnail `json:"thumbnail"`
	Color       int32     `json:"color"`
	Fields      []Field   `json:"fields"`
}

const (
	ActionRowType = 1
	ButtonType    = 2
)

const (
	ButtonStyleLink = 5
)

type Component struct {
	Type  int8   `json:"type"`
	Style int8   `json:"style"`
	Label string `json:"label"`
	Url   string `json:"url"`
}

type ComponentRow struct {
	Type       int         `json:"type"`
	Components []Component `json:"components"`
}

type Webhook struct {
	Username   string         `json:"username"`
	AvatarURL  string         `json:"avatar_url"`
	Embeds     []Embed        `json:"embeds"`
	Components []ComponentRow `json:"components"`
}
