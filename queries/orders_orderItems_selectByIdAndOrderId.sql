select oi.id
from order_items oi
where oi.order_id = $1 and oi.id = $2
