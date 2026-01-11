package src

type forum struct {
	ForumID          string `json:"forum_id"`
	ForumURL         string `json:"forum_url"`
	ForumName        string `json:"forum_name"`
	ForumDescription string `json:"forum_description"`
	LastScaned       string `json:"last_scaned"`
}
