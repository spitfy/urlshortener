BEGIN;
CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
ALTER TABLE urls ADD COLUMN IF NOT EXISTS user_id INT;
COMMIT;
-- todo может быть и не нужна табл users... из за внешнего ключа тест не проходит по inc10
-- DO $$
--     BEGIN
--         IF NOT EXISTS (
--             SELECT 1
--             FROM information_schema.table_constraints
--             WHERE table_name = 'urls'
--               AND constraint_type = 'FOREIGN KEY'
--               AND constraint_name = 'fk_user_id'
--         ) THEN
--             ALTER TABLE urls
--                 ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id);
--         END IF;
--     END
-- $$;