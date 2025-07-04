-- Stories table
CREATE TABLE IF NOT EXISTS stories (
    id INTEGER PRIMARY KEY,
    type VARCHAR(10) DEFAULT 'Story' NOT NULL,
    title TEXT NOT NULL,
    url TEXT,
    score INTEGER DEFAULT 0 CHECK (score >= 0),
    author VARCHAR(255) NOT NULL,
    created_at BIGINT NOT NULL,
    comments_count INTEGER DEFAULT 0 CHECK (comments_count >= 0)
);

-- Asks table
CREATE TABLE IF NOT EXISTS asks (
    id INTEGER PRIMARY KEY,
    type VARCHAR(10) DEFAULT 'Ask' NOT NULL,
    title TEXT NOT NULL,
    text TEXT,
    score INTEGER DEFAULT 0 CHECK (score >= 0),
    author VARCHAR(255) NOT NULL,
    reply_ids INTEGER[] DEFAULT '{}',
    replies_count INTEGER DEFAULT 0 CHECK (replies_count >= 0),
    created_at BIGINT NOT NULL
);

-- Jobs table
CREATE TABLE IF NOT EXISTS jobs (
    id INTEGER PRIMARY KEY,
    type VARCHAR(10) DEFAULT 'Job' NOT NULL,
    title TEXT NOT NULL,
    text TEXT,
    url TEXT,
    score INTEGER DEFAULT 0 CHECK (score >= 0),
    author VARCHAR(255) NOT NULL,
    created_at BIGINT NOT NULL
);

-- Comments table (FIXED - matches updated model)
CREATE TABLE IF NOT EXISTS comments (
    story_id INTEGER NOT NULL,
    id INTEGER PRIMARY KEY,
    type VARCHAR(10) DEFAULT 'Comment' NOT NULL,
    text TEXT NOT NULL,
    author VARCHAR(255) NOT NULL,
    created_at BIGINT NOT NULL,
    parent_id INTEGER,
    reply_ids INTEGER[] DEFAULT '{}',
    FOREIGN KEY (story_id) REFERENCES stories(id) ON DELETE CASCADE
);

-- Polls table (FIXED - matches updated model)
CREATE TABLE IF NOT EXISTS polls (
    id INTEGER PRIMARY KEY,
    type VARCHAR(10) DEFAULT 'Poll' NOT NULL,
    title TEXT NOT NULL,
    score INTEGER DEFAULT 0 CHECK (score >= 0),
    author VARCHAR(255) NOT NULL,
    poll_options TEXT[] DEFAULT '{}',
    reply_ids INTEGER[] DEFAULT '{}',
    created_at BIGINT NOT NULL
);
