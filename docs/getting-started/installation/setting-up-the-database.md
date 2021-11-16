# Setting up the database

WatchOps requires a database. If you’re just experimenting and learning WatchOps, you can stick with the default SQLite option. If you don’t want to use SQLite, then take a look at [Set up a Database Backend](https://airflow.apache.org/docs/apache-airflow/stable/howto/set-up-database.html) to set up a different database.

Usually, you need to run `make setup` in order to create the database schema that WatchOps can use.

Similarly, upgrading WatchOps usually requires an extra step of upgrading the database. This is done with `make migrate` CLI command. You should make sure that WatchOps components are not running while the upgrade is being executed.

In some deployments, such as [Helm Chart for WatchOps](helm.md), both initializing and running the database migration are executed automatically when WatchOps is upgraded.
