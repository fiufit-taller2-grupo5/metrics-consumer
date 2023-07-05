import requests
import random
from concurrent.futures import ThreadPoolExecutor, as_completed
import time
import random

def trigger_metric(metric_name):
    res = requests.post("https://metrics-queue-nodeport-prod2-szwtomas.cloud.okteto.net/api/metrics/system", json={"metric_name": metric_name})
    if res.status_code != 202:
        print(f"Something went wrong, received {res.status_code} trying to enqueue {metric_name}")
    else:
        print("Queued metric: ", metric_name)

def main():
    print("Starting farm!")
    while True:
        metrics = ["user_created", "user_blocked", "user_updated", "training_created", "training_updated", "training_finished"]
        m = random.choice(metrics)
        trigger_metric(m)
        time.sleep(0.05)
            
        



main()
