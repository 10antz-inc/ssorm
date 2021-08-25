CREATE TABLE Singers (
                         SingerId   INT64 NOT NULL,
                         FirstName  STRING(1024),
                         LastName   STRING(1024),
                         DeleteTime TIMESTAMP,
                         CreateTime TIMESTAMP NOT NULL,
                         UpdateTime TIMESTAMP NOT NULL OPTIONS(allow_commit_timestamp=true)
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



INSERT Singers (SingerId, FirstName, LastName,UpdateTime,CreateTime) VALUES (12, 'Melissa', 'Garcia',CURRENT_TIMESTAMP() ,CURRENT_TIMESTAMP() );
INSERT Singers (SingerId, FirstName, LastName,UpdateTime,CreateTime) VALUES (13, 'Russell', 'Morales',CURRENT_TIMESTAMP() ,CURRENT_TIMESTAMP() );
INSERT Singers (SingerId, FirstName, LastName,UpdateTime,CreateTime) VALUES (14, 'Jacqueline', 'Long',CURRENT_TIMESTAMP() ,CURRENT_TIMESTAMP() );
INSERT Singers (SingerId, FirstName, LastName,UpdateTime,CreateTime) VALUES (15, 'Dylan', 'Shaw',CURRENT_TIMESTAMP() ,CURRENT_TIMESTAMP() );


INSERT Albums (SingerId, AlbumId, Title) VALUES (12, 1, 'Garcia');
INSERT Albums (SingerId, AlbumId, Title) VALUES (13, 2, 'Morales');
INSERT Albums (SingerId, AlbumId, Title) VALUES (14, 3, 'Long');
INSERT Albums (SingerId, AlbumId, Title) VALUES (15, 4, 'Shaw');

INSERT Concerts (SingerId, ConcertId, Price) VALUES (12, 1, 100);
INSERT Concerts (SingerId, ConcertId, Price) VALUES (13, 2, 200);
INSERT Concerts (SingerId, ConcertId, Price) VALUES (14, 3, 300);
INSERT Concerts (SingerId, ConcertId, Price) VALUES (15, 4, 400);