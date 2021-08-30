DROP TABLE Singers;
DROP TABLE Albums;
DROP TABLE Concerts;
DROP TABLE Tags;

CREATE TABLE Singers
(
    SingerId        INT64     NOT NULL,
    FirstName       STRING(1024),
    LastName        STRING(1024),
    TestTime        TIMESTAMP,
    TestSpannerTime TIMESTAMP,
    TagIds ARRAY<STRING(36)>,
    Numbers ARRAY<INT64>,
    CreateTime      TIMESTAMP NOT NULL,
    UpdateTime      TIMESTAMP NOT NULL OPTIONS(allow_commit_timestamp= true),
    DeleteTime      TIMESTAMP
) PRIMARY KEY (SingerId);

CREATE TABLE Albums
(
    SingerId INT64 NOT NULL,
    AlbumId  INT64 NOT NULL,
    Title    STRING(1024),
) PRIMARY KEY (SingerId);


CREATE TABLE Concerts
(
    SingerId  INT64 NOT NULL,
    ConcertId INT64 NOT NULL,
    Price     INT64 NOT NULL,
) PRIMARY KEY (SingerId);

CREATE TABLE Tags (
    TagId      STRING(36) NOT NULL,
    Name       STRING(256) NOT NULL,
    DeleteTime TIMESTAMP,
    CreateTime TIMESTAMP NOT NULL,
    UpdateTime TIMESTAMP NOT NULL OPTIONS(allow_commit_timestamp=true)
) PRIMARY KEY(TagId);


INSERT
Singers (SingerId, FirstName, LastName, TagIds, Numbers, UpdateTime, CreateTime) VALUES (12, 'Melissa', 'Garcia', ["a3eb54bd-0138-4c22-b858-41bbefc5c050", "a3eb54bd-0138-4c22-b858-41bbefc5c051"], [1, 2], CURRENT_TIMESTAMP() ,CURRENT_TIMESTAMP());
INSERT
Singers (SingerId, FirstName, LastName, TagIds, Numbers, UpdateTime, CreateTime) VALUES (13, 'Russell', 'Morales', ["a3eb54bd-0138-4c22-b858-41bbefc5c050", "a3eb54bd-0138-4c22-b858-41bbefc5c051"], [1, 2], CURRENT_TIMESTAMP() ,CURRENT_TIMESTAMP());
INSERT
Singers (SingerId, FirstName, LastName, TagIds, Numbers, UpdateTime, CreateTime) VALUES (14, 'Jacqueline', 'Long', ["a3eb54bd-0138-4c22-b858-41bbefc5c050", "a3eb54bd-0138-4c22-b858-41bbefc5c051"], [1, 2], CURRENT_TIMESTAMP() ,CURRENT_TIMESTAMP());
INSERT
Singers (SingerId, FirstName, LastName, TagIds, Numbers, UpdateTime, CreateTime) VALUES (15, 'Dylan', 'Shaw', ["a3eb54bd-0138-4c22-b858-41bbefc5c050", "a3eb54bd-0138-4c22-b858-41bbefc5c051"], [1, 2], CURRENT_TIMESTAMP() ,CURRENT_TIMESTAMP());


INSERT
Albums (SingerId, AlbumId, Title) VALUES (12, 1, 'Garcia');
INSERT
Albums (SingerId, AlbumId, Title) VALUES (13, 2, 'Morales');
INSERT
Albums (SingerId, AlbumId, Title) VALUES (14, 3, 'Long');
INSERT
Albums (SingerId, AlbumId, Title) VALUES (15, 4, 'Shaw');

INSERT
Concerts (SingerId, ConcertId, Price) VALUES (12, 1, 100);
INSERT
Concerts (SingerId, ConcertId, Price) VALUES (13, 2, 200);
INSERT
Concerts (SingerId, ConcertId, Price) VALUES (14, 3, 300);
INSERT
Concerts (SingerId, ConcertId, Price) VALUES (15, 4, 400);

INSERT
Tags (TagId, Name, CreateTime, UpdateTime) VALUES ("a3eb54bd-0138-4c22-b858-41bbefc5c050", "Rock", CURRENT_TIMESTAMP() ,CURRENT_TIMESTAMP());
INSERT
Tags (TagId, Name, CreateTime, UpdateTime) VALUES ("a3eb54bd-0138-4c22-b858-41bbefc5c051", "Pop", CURRENT_TIMESTAMP() ,CURRENT_TIMESTAMP());
INSERT
Tags (TagId, Name, CreateTime, UpdateTime) VALUES ("a3eb54bd-0138-4c22-b858-41bbefc5c052", "Anime", CURRENT_TIMESTAMP() ,CURRENT_TIMESTAMP());
INSERT
Tags (TagId, Name, CreateTime, UpdateTime) VALUES ("a3eb54bd-0138-4c22-b858-41bbefc5c053", "Dance", CURRENT_TIMESTAMP() ,CURRENT_TIMESTAMP());