CREATE TABLE users(
    user_id uuid DEFAULT uuid_generate_v4(),
    nickname varchar(55) not null,
    password_hash varchar(255) not null,
    email varchar(255) not null unique,
    age int,
    image_url text,
    PRIMARY KEY(user_id)
);

CREATE TABLE friend_list(
    user_id uuid not null,
    friend_id uuid not null,
    CONSTRAINT fk_user_id
        FOREIGN KEY(user_id)
            REFERENCES users(user_id)
);

