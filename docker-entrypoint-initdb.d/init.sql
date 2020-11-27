CREATE TABLE category(
    categoryID SERIAL PRIMARY KEY,
    name Text NOT NULL
);

CREATE TABLE role (
    roleID SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE users (
    userId SERIAL PRIMARY KEY,
    username Text NOT NULL,
    hash TEXT NOT NULL,
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    realName Text,
    email Text,
    role Integer NULL REFERENCES role(roleID)
);

CREATE TABLE articles (
    articleID SERIAL PRIMARY KEY,
    title Text NOT Null,
    author Integer NOT Null REFERENCES users(userID),
    body Text NOT NULL,
    date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    slug Text NOT NULL,
    category Integer NOT NULL DEFAULT 1 REFERENCES category(categoryID)
);

INSERT INTO public."role" ("name") VALUES('admin');
