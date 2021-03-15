select oi.id, oi.product_id, p.name as product_name,
	oi.fruit_id, f.name as fruit_name, oi.quantity, oi.unit_price
from order_items oi
inner join products p on p.id = oi.product_id
inner join fruits f on f.id = oi.fruit_id
where oi.order_id = $1
