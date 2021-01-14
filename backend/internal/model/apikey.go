package model

// APIkeyTableSQL is the SQL statement for creating a table corresponding to the APIkey model
var APIkeyTableSQL = `CREATE TABLE IF NOT EXISTS apikey (
	key_hash CHARACTER(64) PRIMARY KEY,
	account_id CHARACTER(36));`
