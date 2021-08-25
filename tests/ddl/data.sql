CREATE TABLE Singers (
                         SingerId   INT64 NOT NULL,
                         FirstName  STRING(1024),
                         LastName   STRING(1024),
) PRIMARY KEY (SingerId);

CREATE TABLE Albums (
                        SingerId   INT64 NOT NULL,
                        AlbumId    INT64 NOT NULL,
                        Title      STRING(1024),
) PRIMARY KEY (SingerId);


CREATE TABLE Concerts (
                          SingerId   INT64 NOT NULL,
                          ConcertId  INT64 NOT NULL,
                          Price      INT64 NOT NULL,
) PRIMARY KEY (SingerId);



INSERT Singers (SingerId, FirstName, LastName) VALUES (12, 'Melissa', 'Garcia');
INSERT Singers (SingerId, FirstName, LastName) VALUES (13, 'Russell', 'Morales');
INSERT Singers (SingerId, FirstName, LastName) VALUES (14, 'Jacqueline', 'Long');
INSERT Singers (SingerId, FirstName, LastName) VALUES (15, 'Dylan', 'Shaw');


INSERT Albums (SingerId, AlbumId, Title) VALUES (12, 1, 'Garcia');
INSERT Albums (SingerId, AlbumId, Title) VALUES (13, 2, 'Morales');
INSERT Albums (SingerId, AlbumId, Title) VALUES (14, 3, 'Long');
INSERT Albums (SingerId, AlbumId, Title) VALUES (15, 4, 'Shaw');

INSERT Concerts (SingerId, ConcertId, Price) VALUES (12, 1, 100);
INSERT Concerts (SingerId, ConcertId, Price) VALUES (13, 2, 200);
INSERT Concerts (SingerId, ConcertId, Price) VALUES (14, 3, 300);
INSERT Concerts (SingerId, ConcertId, Price) VALUES (15, 4, 400);