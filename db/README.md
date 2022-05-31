## Requirements
- aws cli
- mysql

## Create and configure the db
Create the instance using aws cli
```bash
DB_USERNAME=<your_user_name> DB_PASSWORD=<you_password> make create-instance
```
Go to the aws console and get the db hostname and port.

Be sure to add your ip address to the allowed traffic

Run the next command to create the table
```bash
DB_USERNAME=<your_user_name> DB_PASSWORD=<you_password> DB_HOST=<db_host> DB_PORT=<db_port> make run-migrations
```