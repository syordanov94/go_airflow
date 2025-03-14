# Step 1: Importing Modules
# To initiate the DAG Object
from airflow import DAG
# Importing datetime and timedelta modules for scheduling the DAGs
from datetime import timedelta, datetime
# Importing operators 
from airflow.operators.dummy_operator import DummyOperator
from airflow.operators.python import PythonOperator

def write_simple_file(**kwargs):
    file_path = 'hello_airflow.txt'
    with open(file_path, 'w') as file:
        file.write("Hello from Airflow DAG!\n")

# Step 2: Initiating the default_args
default_args = {
        'owner' : 'airflow',
        'start_date' : datetime(2022, 11, 12),
}

# Step 3: Creating DAG Object
dag = DAG(dag_id='example_dag',
        default_args=default_args,
        schedule_interval='*/5 * * * *',
        catchup=False
    )

write_file_task = PythonOperator(
    task_id='write_file_task',
    python_callable=write_simple_file,
    dag=dag,
)

# Step 4: Creating task
# Creating first task
start = DummyOperator(task_id = 'start', dag = dag)
# Creating second task 
end = DummyOperator(task_id = 'end', dag = dag)

# Step 5: Setting up dependencies 
start >> write_file_task >> end