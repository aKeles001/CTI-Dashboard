PRAGMA FOREIGN_KEYS = ON;


--Forums Table
CREATE TABLE IF NOT EXISTS forums (
    forum_id TEXT PRIMARY KEY,
    forum_name TEXT NOT NULL,
    forum_url TEXT NOT NULL,
    forum_html TEXT,
    forum_description TEXT,
    forum_screenshot TEXT,
    last_scaned DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);


--Posts Table
CREATE TABLE IF NOT EXISTS posts (
    post_id TEXT PRIMARY KEY,
    forum_id TEXT,
    author TEXT,
    content TEXT,
    severity REAL DEFAULT 0.0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (forum_id) REFERENCES forums(forum_id) ON DELETE CASCADE
);