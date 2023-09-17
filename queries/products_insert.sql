insert into products (name, unit, deleted, date_created)
values ($1, $2, false, now()) returning id
