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
    forum_engine TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS posts (
    post_id TEXT PRIMARY KEY,
    forum_id TEXT,
    thread_url TEXT UNIQUE,
    status TEXT DEFAULT 'pending', -- pending, scraped, analyzed, failed
    severity_level NUMERIC DEFAULT 0,
    title TEXT,
    content TEXT,
    author TEXT,
    date TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(forum_id) REFERENCES forums(forum_id)
);