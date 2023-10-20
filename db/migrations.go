package db

var Migrations = `
CREATE TABLE IF NOT EXISTS dev_challenge (
id INTEGER PRIMARY KEY AUTOINCREMENT,
sheet_id VARCHAR(255),
cell_id VARCHAR(255),
cell_value TEXT,
cell_result DECIMAL
);
CREATE INDEX IF NOT EXISTS sheet_id_idx ON dev_challenge (sheet_id);
CREATE INDEX IF NOT EXISTS cell_id_idx ON dev_challenge (cell_id);
CREATE UNIQUE INDEX IF NOT EXISTS unique_sheet_cell_idx ON dev_challenge(sheet_id, cell_id);

CREATE TABLE IF NOT EXISTS string_array (
id INTEGER PRIMARY KEY AUTOINCREMENT,
dev_challenge_id INTEGER,
string_value VARCHAR(255),
FOREIGN KEY(dev_challenge_id) REFERENCES dev_challenge(id)
);
INSERT INTO dev_challenge (sheet_id, cell_id, cell_value, cell_result) VALUES ('0','0','0',0) ON CONFLICT DO NOTHING;`
