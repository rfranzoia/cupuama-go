update orders
set date_updated = now()
where id = $2
