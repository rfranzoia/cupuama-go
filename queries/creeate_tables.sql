create table fruits (
    id serial primary key,
    name varchar(255) not null,
    initials varchar(4) not null,
    harvest varchar(255),
    deleted boolean not null,
    date_created timestamp not null,
    date_updated timestamp
);

insert into fruits (name, initials, harvest, date_created, deleted) values ('Cupuacu', 'CUPU', 'Agosto', now(), false);
insert into fruits (name, initials, harvest, date_created, deleted) values ('Maracuja', 'MARA', 'Ano Inteiro', now(), false);’
