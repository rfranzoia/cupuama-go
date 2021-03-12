select o.id, o.order_date, o.total_price, 
		os.id as status_id, os.status, os.status_change_date, os.status_description, 
		oi.id as order_item_id, oi.product_id, p.name as product_name, oi.fruit_id, f.name as fruit_name, 
		oi.quantity, oi.unit_price 
from orders o 
inner join order_items oi on oi.order_id = o.id 
inner join (select os.* from (select order_id, max(status) max_status 
							from order_status 
							group by order_id) mos 
							inner join order_status os on os.order_id = mos.order_id and os.status = mos.max_status) os on os.order_id = o.id 
inner join products p on p.id = oi.product_id 
inner join fruits f on f.id = oi.fruit_id 
where o.deleted = false
and (-1 = $1 or o.id = $2)
order by o.id desc