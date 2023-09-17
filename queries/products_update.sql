update products
    set name = $1,
    unit = $2,
    deleted = $3,
    date_updated = now()
where id = $4
