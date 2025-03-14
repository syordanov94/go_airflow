# Step 1: Importing Modules
# To initiate the DAG Object
from airflow import DAG
# Importing datetime and timedelta modules for scheduling the DAGs
from datetime import timedelta, datetime
# Importing operators 
from airflow.operators.dummy_operator import DummyOperator
from airflow.operators.python import PythonOperator
from airflow.models import Variable
import random

def calculate_player_scores(**kwargs):
    user_names_var = Variable.get("users")
    user_names = user_names_var.split(",")

    #Let's do some processing to simulate score calculation.
    scores = {}
    for user in user_names:
        # get the current time in epoch format and multiply it by a random number
        current_time = datetime.now().timestamp()
        random_number = random.randint(0, 100)

        # calculate the score
        score = current_time * random_number
        scores[user] = score

    # join the scores into a single string and store them in the variable value
    users = []
    for user, score in scores.items():
        users.append("#".join([user, str(score)]))
    scores_str = "|".join(users)
    Variable.set("player_scores", scores_str)

default_args = {
        'owner' : 'airflow',
        'start_date' : datetime(2022, 11, 12),
}

dag = DAG(dag_id='player_score_calculation_dag',
        default_args=default_args,
        schedule_interval='*/5 * * * *',
        catchup=False
    )

calculate_player_scores = PythonOperator(
    task_id='calculate_player_scores',
    python_callable=calculate_player_scores,
    dag=dag,
)

# Step 4: Creating task
# Creating first task
start = DummyOperator(task_id = 'start', dag = dag)
# Creating second task 
end = DummyOperator(task_id = 'end', dag = dag)

# Step 5: Setting up dependencies 
start >> calculate_player_scores >> end