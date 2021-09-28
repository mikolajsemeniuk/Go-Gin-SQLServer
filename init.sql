use master
drop database if exists db
create database db
use db

create table users (
    user_id INT PRIMARY KEY IDENTITY(1,1), -- identity tells sql server to start `user_id` from 1 and increment it by 1
    username varchar(max)
)

create table posts (
    post_id INT PRIMARY KEY IDENTITY(1,1),
    title varchar(max),
    user_id int foreign key references users(user_id) on delete cascade
)

create table user_likes (
    followed_id INT NOT NULL,
    follower_id INT NOT NULL,
    constraint user_like_id primary key (followed_id, follower_id),
    constraint fk_followed foreign key (followed_id) references users (user_id),
    constraint fk_follower foreign key (follower_id) references users (user_id)
)

insert into users (username) values ('john doe'), ('mike mock'), ('lucy applegate'), ('sam taylor')
insert into posts (title, user_id) values ('lorem', 1), ('ipsum', 1), ('dolor', 2), ('sit', 2), ('amet', 3)
insert into user_likes (followed_id, follower_id) values (1, 3), (1, 2), (2, 4), (3, 1), (3, 2), (3, 4)

delete from users where user_id = 3
select * from users
select * from user_likes