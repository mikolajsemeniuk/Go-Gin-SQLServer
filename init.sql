-- https://www.mssqltips.com/sqlservertip/2733/solving-the-sql-server-multiple-cascade-path-issue-with-a-trigger/
-- https://www.sqlservertutorial.net/sql-server-triggers/sql-server-create-trigger/
use master
drop database if exists db
create database db
use db

go

create table users (
    user_id int primary key IDENTITY(1,1), -- identity tells sql server to start `user_id` from 1 and increment it by 1
    username varchar(max)
)

create table posts (
    post_id INT PRIMARY KEY IDENTITY(1,1),
    title varchar(max),
    user_id int foreign key references users(user_id) --on delete cascade --useless while we use trigger to delete posts while user created
)

create table user_likes (
    followed_id int not null,
    follower_id int not null,
    constraint user_like_id primary key (followed_id, follower_id),
    constraint fk_followed foreign key (followed_id) references users (user_id),
    constraint fk_follower foreign key (follower_id) references users (user_id)
)

create table post_likes (
    user_id int not null,
    post_id int not null,
    constraint post_like_id primary key (user_id, post_id),
    constraint fk_user_id foreign key (user_id) references users (user_id),
    constraint fk_post_id foreign key (post_id) references posts (post_id)
)

go

create trigger [delete_user]
   on users
   instead of delete
as
begin
 set nocount on;
 delete from [post_likes] where user_id in (select user_id from deleted)
 delete from [user_likes] where (followed_id in (select user_id from deleted)) or (follower_id in (select user_id from deleted))
 delete from [posts] where user_id in (select user_id from deleted)
 delete from users where user_id in (select user_id from deleted)
end

go

create trigger [delete_post]
    on posts
    instead of delete
as
begin
    set nocount on;
    delete from [post_likes] where post_id in (select post_id from deleted)
    delete from [posts] where post_id in (select post_id from deleted)
end

go

insert into users (username) values ('john doe'), ('mike mock'), ('lucy applegate'), ('sam taylor');
insert into posts (title, user_id) values ('lorem', 1), ('ipsum', 1), ('dolor', 2), ('sit', 2), ('amet', 3);
insert into user_likes (followed_id, follower_id) values (1, 3), (1, 2), (2, 4), (3, 1), (3, 2), (3, 4);
insert into post_likes (user_id, post_id) values (1, 3), (1, 4), (1, 5), (2, 2), (2, 5), (4, 5), (4, 4);

--4, 5, 6, 7
go
select * from users;
select * from posts;
select * from user_likes;
select * from post_likes;
-- go
-- delete from users where user_id = 1;
-- delete from posts where post_id = 5;
-- go
-- select * from users;
-- select * from posts;
-- select * from user_likes;
-- select * from post_likes;