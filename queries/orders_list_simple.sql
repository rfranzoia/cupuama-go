select o.id, to_char(o.order_date, 'yyyy-mm-dd') as order_date, o.total_price,
       os.id as status_id, os.status, to_char(os.status_change_date, 'yyyy-mm-dd hh24:mi') as status_change_date, os.status_description, 
       o.deleted, to_char(o.date_created, 'yyyy-mm-dd hh24:mi') as date_created, coalesce(to_char(o.date_updated, 'yyyy-mm-dd hh24:mi'),'') as date_updated
from orders o
inner join (select os.* from (select order_id, max(status) max_status
		                     from order_status 
			                 group by order_id) mos
	                         inner join order_status os on os.order_id = mos.order_id and os.status = mos.max_status) os on os.order_id = o.id
where deleted = false
order by order_date desc, o.id asc
