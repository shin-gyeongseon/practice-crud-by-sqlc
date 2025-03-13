-- CREATE TABLE "accounts" (
--   "id" bigserial PRIMARY KEY,
--   "owner" VARCHAR NOT NULL,
--   "balance" bigint NOT NULL,
--   "currency" varchar NOT NULL,
--   "created_at" timestamptz NOT NULL DEFAULT (now())
  
-- );

-- CREATE TABLE "entries" (
--   "id" bigserial PRIMARY KEY,
--   "account_id" bigint,
--   "amount" bigint NOT NULL,
--   "created_at" timestamptz NOT NULL DEFAULT (now())
-- );

-- CREATE TABLE "transfers" (
--   "id" bigserial PRIMARY KEY,
--   "from_account_id" bigint,
--   "to_account_id" bigint,
--   "amount" bigint NOT NULL,
--   "created_at" timestamptz NOT NULL DEFAULT (now())
-- );

-- CREATE INDEX ON "transfers" ("from_account_id");

-- CREATE INDEX ON "transfers" ("to_account_id");

-- CREATE INDEX ON "transfers" ("from_account_id", "to_account_id");

-- COMMENT ON COLUMN "entries"."amount" IS 'can be negative our?';

-- COMMENT ON COLUMN "transfers"."amount" IS 'must be positive';

-- ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

-- ALTER TABLE "transfers" ADD FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");

-- ALTER TABLE "transfers" ADD FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id");


--- 

CREATE TABLE "accounts" (
  "id" bigint PRIMARY KEY,
  "owner" varvhar NOT NULL,
  "balance" bigint NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestampz NOT NULL DEFAULT (now())
);

CREATE TABLE "entries" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigint,
  "amount" bigint NOT NULL,
  "created_at" timestampz NOT NULL DEFAULT (now())
);

CREATE TABLE "transfers" (
  "id" bigint PRIMARY KEY,
  "from_account_id" bigint,
  "to_account_id" bigint,
  "amount" bigint NOT NULL,
  "created_at" timestampz NOT NULL DEFAULT (now())
);

CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "hashed_password" varchar NOT NULL,
  "full_name" varvhar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password_changed_at" timestampz NOT NULL DEFAULT '0001-01-01 00:00:00+09',
  "created_at" timestampz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "entries" ("account_id");

CREATE INDEX ON "transfers" ("from_account_id");

CREATE INDEX ON "transfers" ("to_account_id");

CREATE INDEX ON "transfers" ("from_account_id", "to_account_id");

COMMENT ON COLUMN "entries"."amount" IS 'can be negative our?';

COMMENT ON COLUMN "transfers"."amount" IS 'must be positive';

ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id");
