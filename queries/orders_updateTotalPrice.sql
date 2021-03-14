update orders
set date_updated = now(),
    total_price = $1
where id = $2
