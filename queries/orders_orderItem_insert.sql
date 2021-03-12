insert into order_items (order_id, product_id, fruit_id, quantity, unit_price) 
values ($1, $2, $3, $4, $5) 
returning id