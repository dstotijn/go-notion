package notion

type Person struct {
	Email string `json:"email"`
}

type Bot struct{}

type User struct {
	ID        string  `json:"id"`
	Type      string  `json:"type"`
	Name      string  `json:"name"`
	AvatarURL *string `json:"avatar_url"`

	Person *Person `json:"person"`
	Bot    *Bot    `json:"bot"`
}
