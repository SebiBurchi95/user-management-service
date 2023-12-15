This is a service that exposes a REST API that performs CRUD operations on a users table in an SQLite database using the ent library.

Please make sure you have installed go 1.21 on your setup.

For running the app as a standalone service, run: make all (the server will start on port 8080)

For testing the app using a REST API client, run: ./api_test.sh (please make sure that you have righs on the script)

For running the unite tests with coverage, run: make test
