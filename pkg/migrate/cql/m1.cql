-- Comment

CREATE TABLE IF NOT EXISTS tags (id uuid, namespace text, resource text, key text, value text, PRIMARY KEY ((namespace), resource, key, value));

-- Call 1

CREATE INDEX IF NOT EXISTS ON tags(value);