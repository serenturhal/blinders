import os

import dotenv
import pymongo
import redis

from blinders.explore_core.embedder import Embedder
from blinders.explore_core.main import Explore
from blinders.explore_service.core.main import ServiceWorker

matchColName = "matches"
if __name__ == "__main__":
    dotenv.load_dotenv()
    try:
        mongoURL = "mongodb://{}:{}@{}:{}/{}".format(
            os.getenv("MONGO_USERNAME"),
            os.getenv("MONGO_PASSWORD"),
            os.getenv("MONGO_HOST"),
            os.getenv("MONGO_PORT"),
            os.getenv("MONGO_DATABASE"),
        )
        mongoClient = pymongo.MongoClient(mongoURL)
        db = mongoClient[os.getenv("MONGO_DATABASE", "Default")]
        matchCol = db[matchColName]

        embedder = Embedder()
        redis_client = redis.client.Redis(host="localhost", port=6379, decode_responses=True)
        explore = Explore(redis_client, embedder, matchCol)
        service_core = ServiceWorker(redis_client=redis_client, explore_core=explore)
        service_core.loop()

    except Exception as e:
        print("exception: ", e)
