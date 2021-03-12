select order_id, max(status) max_status
from order_status 
group by order_id
having order_id = $1