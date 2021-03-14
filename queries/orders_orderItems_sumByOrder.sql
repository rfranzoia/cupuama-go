select sum(oi.unit_price * oi.quantity) as total_price
from order_items oi
where oi.order_id = $1
