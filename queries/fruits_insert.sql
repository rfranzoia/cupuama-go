insert into fruits (name, harvest, initials, deleted, date_created)
values ($1, $2, $3, false, now()) returning id
