package postgres

const createUserQuery = `
INSERT INTO "Users"(id, first_name, last_name, age, profile) 
VALUES ($1, $2, $3, $4, $5)`

const getUserQuery = `
SELECT id, first_name, last_name, age, profile FROM "Users" WHERE id = $1`

const deleteUserQuery = `
DELETE FROM "Users" WHERE id = $1`

const updateUserQuery = `
UPDATE "Users" SET first_name = $2, last_name = $3, age = $4, profile = $5 WHERE id = $1`
