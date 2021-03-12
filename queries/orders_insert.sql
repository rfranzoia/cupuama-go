insert into orders (order_date, total_price) 
values (now(), $1)  
returning id