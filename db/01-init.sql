CREATE TABLE IF NOT EXISTS "expenses" (
    "id" SERIAL PRIMARY KEY,
    "title" VARCHAR(255) NOT NULL,
    "amount" INT NOT NULL,
    "note" TEXT NOT NULL,
    "tags" VARCHAR(255) [] NOT NULL
);
INSERT INTO "expenses" ("title", "amount", "note", "tags")
VALUES (
        'Coffee',
        100,
        'Coffee from Starbucks',
        ARRAY ['food', 'coffee']
    );