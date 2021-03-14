select id, name, unit, deleted, 
       to_char(date_created, 'yyyy-mm-dd hh24:mi') as date_created, 
       coalesce(to_char(date_updated, 'yyyy-mm-dd hh24:mi'), '') as date_updated
from products
where deleted = false and id = $1