-- Create "posts" table
CREATE TABLE "public"."posts" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "user_id" bigint NOT NULL,
  "title" text NOT NULL,
  "content" text NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_posts_deleted_at" to table: "posts"
CREATE INDEX "idx_posts_deleted_at" ON "public"."posts" ("deleted_at");
-- Create "profiles" table
CREATE TABLE "public"."profiles" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "user_id" bigint NOT NULL,
  "name" text NULL,
  "bio" text NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_profiles_deleted_at" to table: "profiles"
CREATE INDEX "idx_profiles_deleted_at" ON "public"."profiles" ("deleted_at");
-- Create index "idx_profiles_user_id" to table: "profiles"
CREATE UNIQUE INDEX "idx_profiles_user_id" ON "public"."profiles" ("user_id");
-- Create "users" table
CREATE TABLE "public"."users" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "email" text NOT NULL,
  "password" text NOT NULL,
  "role" bigint NOT NULL DEFAULT 2,
  PRIMARY KEY ("id")
);
-- Create index "idx_email_unique" to table: "users"
CREATE UNIQUE INDEX "idx_email_unique" ON "public"."users" ("email") WHERE (deleted_at IS NULL);
-- Create index "idx_users_deleted_at" to table: "users"
CREATE INDEX "idx_users_deleted_at" ON "public"."users" ("deleted_at");
