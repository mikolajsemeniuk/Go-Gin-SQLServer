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

insert into users (username) values ('john doe'), ('mike mock'), ('lucy applegate'), ('sam taylor')
insert into posts (title, user_id) values ('lorem', 1), ('ipsum', 1), ('dolor', 2), ('sit', 2), ('amet', 3)