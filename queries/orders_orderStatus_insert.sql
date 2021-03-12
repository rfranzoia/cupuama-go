insert into order_status (order_id, status, status_description) 
values ($1, $2, $3) 
returning id