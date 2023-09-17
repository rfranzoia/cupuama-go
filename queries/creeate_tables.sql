create table fruits (
    id serial primary key,
    name varchar(255) not null,
    initials varchar(4) not null,
    harvest varchar(255),
    deleted boolean not null,
    date_created timestamp not null,
    date_updated timestamp
);

alter table fruits owner to cupuama;

insert into fruits (name, initials, harvest, date_created, deleted) values ('Cupuacu', 'CUPU', 'Agosto', now(), false);
insert into fruits (name, initials, harvest, date_created, deleted) values ('Maracuja', 'MARA', 'Ano Inteiro', now(), false);’

create table products (
    id serial primary key,
    name varchar(255) not null,
    unit varchar(4) not null,
    deleted boolean not null,
    date_created timestamp not null,
    date_updated timestamp
);

alter table products owner to cupuama;

insert into products (name, unit, date_created, deleted) values ('Polpa 500g', 'PCT', now(), false);
insert into products (name, unit, date_created, deleted) values ('Polpa 10kg', 'LATA',now(), false);
