update fruits
    set name = $1,
    harvest = $2,
    initials = $3,
    deleted = $4,
    date_updated = now()
where id = $5
