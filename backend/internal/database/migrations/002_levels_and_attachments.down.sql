ALTER TABLE ops_requests ALTER COLUMN request_date TYPE DATE USING request_date::date;

ALTER TABLE attachments
    DROP COLUMN IF EXISTS checksum,
    DROP COLUMN IF EXISTS mime_type,
    DROP COLUMN IF EXISTS file_ext,
    DROP COLUMN IF EXISTS file_name,
    DROP COLUMN IF EXISTS file_size;

ALTER TABLE request_types DROP COLUMN IF EXISTS required_level_rank;

DROP TABLE IF EXISTS user_levels CASCADE;
DROP TABLE IF EXISTS levels CASCADE;