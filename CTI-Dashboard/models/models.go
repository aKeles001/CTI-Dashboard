package models

type Forum struct {
	ForumID          string `json:"forum_id"`
	ForumURL         string `json:"forum_url"`
	ForumName        string `json:"forum_name"`
	ForumDescription string `json:"forum_description"`
	LastScaned       string `json:"last_scaned"`
	ForumHTML        string `json:"forum_html"`
	ForumScreenshot  string `json:"forum_screenshot"`
	ForumEngine      string `json:"forum_engine"`
}

type Post struct {
	PostID      string `json:"post_id"`
	ForumID     string `json:"forum_id"`
	ThreadURL   string `json:"thread_url"`
	Status      string `json:"status"`
	Severity    string `json:"severity_level"`
	Title       string `json:"title"`
	PostContent string `json:"content"`
	PostAuthor  string `json:"author"`
	PostDate    string `json:"date"`
}

type Chart struct {
	ForumID    string `json:"forum_id"`
	ForumName  string `json:"forum_name"`
	ForumURL   string `json:"forum_url"`
	PostCount  int    `json:"count"`
	High       int    `json:"high"`
	Medium     int    `json:"medium"`
	Low        int    `json:"low"`
	Unassigned int    `json:"unassigned"`
	LastScaned string `json:"last_scaned"`
}

type Job struct {
	JobID     string `json:"job_id"`
	ThreadURL string `json:"thread_url"`
	ForumName string `json:"forum_name"`
}
