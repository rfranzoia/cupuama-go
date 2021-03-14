select login, password, first_name, last_name, date_of_birth, deleted, 
	   to_char(date_created, 'yyyy-mm-dd hh24:mi') as date_created, 
	   coalesce(to_char(date_updated, 'yyyy-mm-dd hh24:mi'), '') as date_updated
from users 
where deleted = false