# Go-Airflow Lab

In this repository, we will explore how to interact with Apache Airflow using Go. We will create a simple API using Gin-Go that will interact with an local instance of Apache Aifrflow.

This repository presents an example of a API written in Go that interacts with Apache Airflow and updates variables and DAGs to achieve desired
results.

## Pre-requisites

    - Docker installed
    - Golang 1.22 or higher
    - A knowledge of Apache Airflow DAG creation and management. This includes some knowledge of Python.

## Project Description

We will create an API that creates (and updates) users names for a made up video game. This API will publish those names to Apache Airflow in the form of a variable. This variable will be consumed by an example DAG that will process the scores of the users and update those in some storage system (maybe another variable or a database entry) periodically.

The API will have the following endpoints:

- *POST /users*: This endpoint will receive a JSON payload with the user name and will create a new user in the system. This user will be added to the Airflow variable. The JSON payload will look like this:

```json
    {
        "users": [
            {
                "name": "Fede"
            },
            {
                "name": "Vini"
            },
            {
                "name": "Jude"
            }
        ]
    }
```

- *GET /users/{userName}/score*: This endpoint will return the score of the user with the given name. The response will have the following format

```json
    {
        "users": [
            {
                "name": "Fede"
            },
            {
                "name": "Vini"
            },
            {
                "name": "Jude"
            }
        ]
    }
```

- *POST /users/score*: This endpoint will manually trigger the DAG that will update the scores of the users.

## How to run the project

Before starting the API, we need to start the local Airflow instance. To do this we need to run the following command:

```bash
    docker-compose up
```

Once that is done, we can start the API by running the following command:

```bash
    go run cmd/main.go
```

The API will be available at <http://localhost:8080>.

## Helpful Tips

- You can visualize the DAG execution in the apache airflow UI by going to <http://localhost:8080> and logging in with the credentials `admin:admin`.

## Helpful Links

- [Apache Airflow DAGs](https://airflow.apache.org/docs/apache-airflow/stable/core-concepts/dags.html)
- [Running Apache Airflow Locally](https://airflow.apache.org/docs/apache-airflow/stable/howto/docker-compose/index.html)
