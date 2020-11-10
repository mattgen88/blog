CREATE TABLE Category(
    Category Id Integer PRIMARY KEY NOT NULL,
    Name Text NOT NULL
);

CREATE TABLE Role (
    RoleId Integer PRIMARY KEY NOT NULL,
    Name TEXT NOT NULL
);

CREATE TABLE Users (
    UserId Integer PRIMARY KEY NOT NULL,
    Username Text NOT NULL,
    Hash TEXT NOT NULL,
    Created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    RealName Text,
    Email Text,
    Role Integer NULL REFERENCES Role(RoleID)
);

CREATE TABLE Articles (
    ArticleId Integer PRIMARY KEY NOT NULL,
    Title Text Not Null,
    Author Integer Not Null REFERENCES Users(UserId),
    Body Text NOT NULL,
    Date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    Slug Text NOT NULL,
    Category Integer NOT NULL DEFAULT 1 REFERENCES Category(CategoryId)
);



