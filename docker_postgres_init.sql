CREATE TABLE category(
    categoryID Integer PRIMARY KEY NOT NULL,
    name Text NOT NULL
);

CREATE TABLE role (
    roleID Integer PRIMARY KEY NOT NULL,
    name TEXT NOT NULL
);

CREATE TABLE users (
    userId Integer PRIMARY KEY NOT NULL,
    username Text NOT NULL,
    hash TEXT NOT NULL,
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    realName Text,
    email Text,
    role Integer NULL REFERENCES role(roleID)
);

CREATE TABLE articles (
    articleID Integer PRIMARY KEY NOT NULL,
    title Text Not Null,
    author Integer Not Null REFERENCES users(userID),
    body Text NOT NULL,
    date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    slug Text NOT NULL,
    category Integer NOT NULL DEFAULT 1 REFERENCES category(categoryID)
);



