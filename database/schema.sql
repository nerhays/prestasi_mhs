-- database/schema.sql

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- roles
CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- users
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    role_id UUID REFERENCES roles(id),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- permissions
CREATE TABLE IF NOT EXISTS permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) UNIQUE NOT NULL,
    resource VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,
    description TEXT
);

-- pivot
CREATE TABLE IF NOT EXISTS role_permissions (
    role_id UUID REFERENCES roles(id),
    permission_id UUID REFERENCES permissions(id),
    PRIMARY KEY (role_id, permission_id)
);

-- lecturers
CREATE TABLE IF NOT EXISTS lecturers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    lecturer_id VARCHAR(20) UNIQUE NOT NULL,
    department VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW()
);

-- students
CREATE TABLE IF NOT EXISTS students (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    student_id VARCHAR(20) UNIQUE NOT NULL,
    program_study VARCHAR(100),
    academic_year VARCHAR(10),
    advisor_id UUID REFERENCES lecturers(id),
    created_at TIMESTAMP DEFAULT NOW()
);

-- enum for achievement status
DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'achievement_status') THEN
        CREATE TYPE achievement_status AS ENUM ('draft','submitted','verified','rejected');
    END IF;
END$$;

-- achievement_references
CREATE TABLE IF NOT EXISTS achievement_references (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id UUID REFERENCES students(id),
    mongo_achievement_id VARCHAR(24) NOT NULL,
    status achievement_status,
    submitted_at TIMESTAMP,
    verified_at TIMESTAMP,
    verified_by UUID REFERENCES users(id),
    rejection_note TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Optional: index untuk query paling umum
CREATE INDEX IF NOT EXISTS idx_achievement_ref_student ON achievement_references(student_id);
CREATE INDEX IF NOT EXISTS idx_achievement_ref_status ON achievement_references(status);

-- Seed roles
INSERT INTO roles (name, description)
VALUES 
 ('Admin','Pengelola sistem'),
 ('Mahasiswa','Mahasiswa pelapor prestasi'),
 ('Dosen Wali','Dosen pembimbing akademik')
ON CONFLICT (name) DO NOTHING;

-- Seed permissions
INSERT INTO permissions (name, resource, action, description) VALUES
 ('achievement:create','achievement','create','Buat prestasi'),
 ('achievement:read','achievement','read','Lihat prestasi'),
 ('achievement:update','achievement','update','Update prestasi'),
 ('achievement:delete','achievement','delete','Hapus prestasi'),
 ('achievement:verify','achievement','verify','Verifikasi prestasi'),
 ('user:manage','user','manage','Kelola user')
ON CONFLICT (name) DO NOTHING;

-- (Optional) map some basic permissions to roles (example)
-- You can adjust role-permission mapping later via seeder code or admin UI
-- Example: Admin gets all permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'Admin'
ON CONFLICT DO NOTHING;
