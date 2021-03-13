select o.id, o.order_date, o.total_price, 
       os.id as status_id, os.status, os.status_change_date, os.status_description, 
	   coalesce(oi.id, -1) as order_item_id, coalesce(oi.product_id, -1), coalesce(p.name, '') as product_name, 
	   coalesce(oi.fruit_id, -1), coalesce(f.name, '') as fruit_name, 
	   coalesce(oi.quantity, 0), coalesce(oi.unit_price, 0.0)
from orders o
left outer join order_items oi on o.id = oi.order_id
left join (select os.* from (select order_id, max(status) max_status
		                     from order_status 
			                 group by order_id) mos
	                         inner join order_status os on os.order_id = mos.order_id and os.status = mos.max_status) os on os.order_id = o.id
left join products p on p.id = oi.product_id
left join fruits f on f.id = oi.fruit_id
where o.deleted = false
and (-1 = $1 or o.id = $2)
order by o.id desc