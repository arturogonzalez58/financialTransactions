run-migrations:
	mysql -h $(DB_HOST) -P $(DB_PORT) -u $(DB_USERNAME) -p < migration.sql

create-instance:
	aws rds create-db-instance --db-name transactions --engine MySQL \
	--db-instance-identifier transactions --backup-retention-period 3 \
	--db-instance-class db.t2.micro --allocated-storage 5 --publicly-accessible \
	--master-username $(DB_USERNAME) --master-user-password $(DB_PASSWORD)