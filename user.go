package notion

type UserType string

const (
	UserTypePerson UserType = "person"
	UserTypeBot    UserType = "bot"
)

type Person struct {
	Email string `json:"email"`
}

type Bot struct {
	Owner BotOwner `json:"owner"`
}

type BotOwnerType string

const (
	BotOwnerTypeWorkspace BotOwnerType = "workspace"
	BotOwnerTypeUser      BotOwnerType = "user"
)

type BotOwner struct {
	Type      BotOwnerType `json:"type"`
	Workspace bool         `json:"workspace"`
	User      *User        `json:"user"`
}

// BaseUser contains the fields that are always returned for user objects.
// See: https://developers.notion.com/reference/user#where-user-objects-appear-in-the-api
type BaseUser struct {
	ID string `json:"id"`
}

type User struct {
	BaseUser

	Type      UserType `json:"type"`
	Name      string   `json:"name"`
	AvatarURL string   `json:"avatar_url"`

	Person *Person `json:"person"`
	Bot    *Bot    `json:"bot"`
}

// ListUsersResponse contains results (users) and pagination data returned from a list request.
type ListUsersResponse struct {
	Results    []User  `json:"results"`
	HasMore    bool    `json:"has_more"`
	NextCursor *string `json:"next_cursor"`
}
