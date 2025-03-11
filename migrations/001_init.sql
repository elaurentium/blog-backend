-- migrations/001_init.sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    hashed_password TEXT NOT NULL,
    salt TEXT NOT NULL,
    full_name VARCHAR(100),
    bio TEXT,
    avatar_url TEXT,
    role VARCHAR(50) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    last_login TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE TABLE subs (
    id UUID PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    rules TEXT[],
    creator_id UUID REFERENCES users(id),
    is_private BOOLEAN NOT NULL DEFAULT FALSE,
    banner_url TEXT,
    icon_url TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE TABLE posts (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    user_id UUID REFERENCES users(id),
    sub_id UUID REFERENCES subs(id),
    upvotes INT NOT NULL DEFAULT 0,
    downvotes INT NOT NULL DEFAULT 0,
    is_locked BOOLEAN NOT NULL DEFAULT FALSE,
    is_pinned BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE TABLE comments (
    id UUID PRIMARY KEY,
    content TEXT NOT NULL,
    user_id UUID REFERENCES users(id),
    post_id UUID REFERENCES posts(id),
    parent_id UUID REFERENCES comments(id),
    upvotes INT NOT NULL DEFAULT 0,
    downvotes INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE TYPE vote_type AS ENUM ('upvote', 'downvote');

CREATE TABLE votes (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    post_id UUID REFERENCES posts(id),
    comment_id UUID REFERENCES comments(id),
    type vote_type NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT votes_post_or_comment_check CHECK (
        (post_id IS NOT NULL AND comment_id IS NULL) OR 
        (post_id IS NULL AND comment_id IS NOT NULL)
    ),
    UNIQUE(user_id, post_id, comment_id)
);

CREATE TABLE sub_members (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    sub_id UUID REFERENCES subs(id),
    role VARCHAR(50) NOT NULL DEFAULT 'member', -- member, moderator, admin
    joined_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, sub_id)
);

CREATE TABLE user_notifications (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    type VARCHAR(50) NOT NULL, -- comment_reply, post_reply, vote, mention
    content TEXT NOT NULL,
    related_post_id UUID REFERENCES posts(id),
    related_comment_id UUID REFERENCES comments(id),
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE user_subscriptions (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    sub_id UUID REFERENCES subs(id),
    notification_preference VARCHAR(50) NOT NULL DEFAULT 'default', -- default, all, none
    subscribed_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, sub_id)
);

-- Indexes for performance optimization
CREATE INDEX idx_posts_user_id ON posts(user_id);
CREATE INDEX idx_posts_sub_id ON posts(sub_id);
CREATE INDEX idx_comments_post_id ON comments(post_id);
CREATE INDEX idx_comments_user_id ON comments(user_id);
CREATE INDEX idx_votes_user_id ON votes(user_id);
CREATE INDEX idx_votes_post_id ON votes(post_id);
CREATE INDEX idx_votes_comment_id ON votes(comment_id);
CREATE INDEX idx_sub_members_user_id ON sub_members(user_id);
CREATE INDEX idx_sub_members_sub_id ON sub_members(sub_id);
CREATE INDEX idx_user_notifications_user_id ON user_notifications(user_id);

-- Additional constraints
ALTER TABLE posts 
ADD CONSTRAINT fk_posts_user_id 
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE posts 
ADD CONSTRAINT fk_posts_sub_id 
FOREIGN KEY (sub_id) REFERENCES subs(id) ON DELETE CASCADE;

ALTER TABLE comments 
ADD CONSTRAINT fk_comments_user_id 
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE comments 
ADD CONSTRAINT fk_comments_post_id 
FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE;

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Triggers to automatically update updated_at
CREATE TRIGGER update_users_modtime
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();

CREATE TRIGGER update_posts_modtime
BEFORE UPDATE ON posts
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();

CREATE TRIGGER update_comments_modtime
BEFORE UPDATE ON comments
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();

CREATE TRIGGER update_subs_modtime
BEFORE UPDATE ON subs
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();