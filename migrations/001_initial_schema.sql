-- Initial database schema for File Print Service
-- Compatible with PostgreSQL 12+ (Neon Database)
-- Created: 2025-10-29

-- ============================================
-- TABLE: folders
-- Stores user-created folders for organizing uploads
-- ============================================
CREATE TABLE IF NOT EXISTS folders (
    id VARCHAR(255) PRIMARY KEY,              -- UUID generated in application
    name VARCHAR(255) NOT NULL,               -- User-defined folder name
    created_at TIMESTAMP NOT NULL DEFAULT NOW(), -- When folder was created
    file_count INTEGER NOT NULL DEFAULT 0     -- Number of files in folder (cached)
);

-- Index for faster folder retrieval by creation date
CREATE INDEX IF NOT EXISTS idx_folders_created_at ON folders(created_at DESC);

-- ============================================
-- TABLE: uploaded_files
-- Stores metadata for all uploaded files
-- Physical files are stored on disk or cloud storage
-- ============================================
CREATE TABLE IF NOT EXISTS uploaded_files (
    id VARCHAR(255) PRIMARY KEY,              -- UUID generated in application
    folder_id VARCHAR(255) NOT NULL,          -- Reference to parent folder
    folder_name VARCHAR(255) NOT NULL,        -- Denormalized for quick display
    file_name VARCHAR(500) NOT NULL,          -- Original filename from user
    file_size BIGINT NOT NULL,                -- Size in bytes
    file_type VARCHAR(50) NOT NULL,           -- Extension (pdf, jpg, png, etc.)
    file_path TEXT NOT NULL,                  -- Path to physical file
    uploaded_at TIMESTAMP NOT NULL DEFAULT NOW(), -- Upload timestamp
    
    -- Foreign key constraint (optional - can remove if folders deleted independently)
    CONSTRAINT fk_folder FOREIGN KEY (folder_id) REFERENCES folders(id) ON DELETE CASCADE
);

-- Indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_uploaded_files_folder_id ON uploaded_files(folder_id);
CREATE INDEX IF NOT EXISTS idx_uploaded_files_uploaded_at ON uploaded_files(uploaded_at DESC);

-- ============================================
-- TABLE: admins
-- Stores admin user credentials
-- Passwords are bcrypt-hashed in application
-- ============================================
CREATE TABLE IF NOT EXISTS admins (
    username VARCHAR(100) PRIMARY KEY,        -- Admin username (unique)
    password_hash VARCHAR(255) NOT NULL,      -- Bcrypt hashed password
    created_at TIMESTAMP NOT NULL DEFAULT NOW() -- When admin was created
);

-- Insert default admin (password: changeme123)
-- IMPORTANT: Change this password immediately in production!
-- The hash below is for 'changeme123'
INSERT INTO admins (username, password_hash) 
VALUES ('admin', '$2a$10$rC5kqXGQ8VQqVK5VYVqKQeGZxhWDYJx8aYXN5KH5jHFQTYGFqKWYO')
ON CONFLICT (username) DO NOTHING;

-- ============================================
-- COMMENTS: Document table purposes
-- ============================================
COMMENT ON TABLE folders IS 'User-created folders for organizing uploaded files';
COMMENT ON TABLE uploaded_files IS 'Metadata for all uploaded files (physical files stored separately)';
COMMENT ON TABLE admins IS 'Admin user accounts for dashboard access';

-- ============================================
-- GRANTS: Ensure proper permissions (adjust for your Neon setup)
-- ============================================
-- If you have a specific database user, grant permissions:
-- GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO your_db_user;
-- GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO your_db_user;
