INSERT INTO USERS (ID, LOGIN, PASSWORD)
VALUES (
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11'::UUID,
    'ODvornikova',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy' 
);

INSERT INTO DOCUMENTS (
    ID, 
    USER_ID, 
    FILENAME, 
    PATH, 
    CREATED_AT,
    MIME,
    HAS_FILE,
    IS_PUBLIC
)
VALUES (
     'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11'::UUID,
    (SELECT ID FROM USERS WHERE LOGIN = 'ODvornikova'), 
    'photo.jpg',
    '/storage/b0/b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11_photo.jpg', 
    '2018-12-24 10:30:56'::TIMESTAMPTZ,
    'image/jpg',
    true,
    false
);

