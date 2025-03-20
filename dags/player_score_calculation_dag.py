from airflow import DAG
from datetime import datetime
import random
# Importing operators 
from airflow.operators.dummy_operator import DummyOperator
from airflow.operators.python import PythonOperator
from airflow.models import Variable

def calculate_player_scores(**kwargs):
    # get the user names from the variable
    user_names_var = Variable.get("users")
    user_names = user_names_var.split(",")

    #Let's do some processing to simulate score calculation.
    scores = calculate_scores(user_names)

    # join the scores into a single string and store them in the variable value
    users = []
    for user, score in scores.items():
        users.append("#".join([user, str(score)]))
    scores_str = "|".join(users)
    Variable.set("player_scores", scores_str)


# This function is currently very simple but you could add more complex logic here if needed that might require a lot of time to process.
# In this case, Airflow execution will be quite handy since it can be schedule it and run it in the Apache environment 
# and your API will not suffer any performance issues and degradation.
def calculate_scores(players):
    #Let's do some processing to simulate score calculation.
    scores = {}
    for user in players:
        # get the current time in epoch format and multiply it by a random number
        current_time = datetime.now().timestamp()
        random_number = random.randint(0, 100)

        # calculate the score
        score = current_time * random_number
        scores[user] = score

    return scores


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

# Let's create two dummy operators in between the main one just to play around a bit more
start = DummyOperator(task_id = 'start', dag = dag)
end = DummyOperator(task_id = 'end', dag = dag)

start >> calculate_player_scores >> end