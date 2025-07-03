
-- Stories table (matches story.go)
CREATE TABLE IF NOT EXISTS stories (
    id INTEGER PRIMARY KEY,           -- Will store the HN story ID
    type VARCHAR(10) DEFAULT 'Story', -- Always 'Story' for this table
    title TEXT NOT NULL,              -- Story title
    url TEXT,                         -- Link to external article
    score INTEGER DEFAULT 0,          -- Upvotes
    author VARCHAR(255) NOT NULL,     -- Username who posted
    created_at BIGINT NOT NULL,       -- Unix timestamp
    comments_count INTEGER DEFAULT 0  -- Number of comments
);

-- Asks table (matches ask.go)
CREATE TABLE IF NOT EXISTS asks (
    id INTEGER PRIMARY KEY,           -- Will store the HN ask ID
    type VARCHAR(10) DEFAULT 'Ask',   -- Always 'Ask' for this table
    title TEXT NOT NULL,              -- Ask title
    text TEXT,                        -- Ask text
    score INTEGER DEFAULT 0,          -- Upvotes
    author VARCHAR(255) NOT NULL,     -- Username who posted
    reply_ids INTEGER[],              -- Array of IDs for replies to this ask
    replies_count INTEGER DEFAULT 0   -- Number of replies to this ask 
    created_at BIGINT NOT NULL,       -- Unix timestamp        
);

-- Jobs table (matches job.go)
CREATE TABLE IF NOT EXISTS jobs (
    id INTEGER PRIMARY KEY,
    type VARCHAR(10) DEFAULT 'Job',
    title TEXT NOT NULL,
    text TEXT,
    url TEXT,                         
    score INTEGER DEFAULT 0,
    author VARCHAR(255) NOT NULL,
    created_at BIGINT NOT NULL      
);

-- Comments table (matches comment.go)
CREATE TABLE IF NOT EXISTS comments (
    story_id INTEGER NOT NULL,          -- (FK) ID of the story this comment belongs to
    id INTEGER PRIMARY KEY,             -- Will store the HN comment ID
    type VARCHAR(10) DEFAULT 'Comment', -- Always 'Comment' for this table
    text TEXT NOT NULL,                 -- Comment text
    author VARCHAR(255) NOT NULL,       -- Username who posted
    created_at BIGINT NOT NULL,         -- Unix timestamp
    parent_id INTEGER,                  -- ID of the parent comment (if any)
    reply_ids INTEGER[]                 -- Array of IDs for replies to this comment
    FOREIGN KEY (story_id) REFERENCES stories(id) ON DELETE CASCADE
);

-- Polls table (matches poll.go)
CREATE TABLE IF NOT EXISTS polls (
    id INTEGER PRIMARY KEY,             -- Will store the HN poll ID
    type VARCHAR(10) DEFAULT 'Poll',    -- Always 'Poll' for this table
    title TEXT NOT NULL,                -- Poll title
    score INTEGER DEFAULT 0,            -- Upvotes
    author VARCHAR(255) NOT NULL,       -- Username who posted
    poll_options TEXT[],                -- Array of poll options
    reply_ids INTEGER[],                -- Array of IDs for replies to this poll
    created_at BIGINT NOT NULL          -- Unix timestamp
);
