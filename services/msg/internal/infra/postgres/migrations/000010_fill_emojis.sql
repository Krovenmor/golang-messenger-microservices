-- +goose Up
insert into emojis(emoji)
values ('💩'), ('🤡'), ('🌈'), ('🔥'), ('🧔🏿‍♂️'), ('💯');

-- +goose Down
truncate table emojis cascade;