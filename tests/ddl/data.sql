CREATE TABLE Singers (
                         SingerId   INT64 NOT NULL,
                         FirstName  STRING(1024),
                         LastName   STRING(1024),
) PRIMARY KEY (SingerId);

INSERT Singers (SingerId, FirstName, LastName) VALUES (12, 'Melissa', 'Garcia');
INSERT Singers (SingerId, FirstName, LastName) VALUES (13, 'Russell', 'Morales');
INSERT Singers (SingerId, FirstName, LastName) VALUES (14, 'Jacqueline', 'Long');
INSERT Singers (SingerId, FirstName, LastName) VALUES (15, 'Dylan', 'Shaw');