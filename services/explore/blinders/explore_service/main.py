from blinders.explore_core.main import Explore
from blinders.explore_core.embedder import Embedder
import pymongo
from redis.client import Redis
import os
from dotenv import load_dotenv

from blinders.explore_service.core.main import ServiceWorker


matchColName = "matchs"
if __name__ == "__main__":
    load_dotenv()
    try:
        mongoURL = "mongodb://{}:{}@{}:{}/{}".format(
            os.getenv("MONGO_USERNAME"),
            os.getenv("MONGO_PASSWORD"),
            os.getenv("MONGO_HOST"),
            os.getenv("MONGO_PORT"),
            os.getenv("MONGO_DATABASE"),
        )
        print(mongoURL)
        mongoClient = pymongo.MongoClient(mongoURL)
        db = mongoClient[os.getenv("MONGO_DATABASE", "Default")]
        matchCol = db[matchColName]

        embbeder = Embedder()
        redisClient = Redis(host="localhost", port=6379, decode_responses=True)
        explore = Explore(redisClient, embbeder, matchCol)
    except Exception as e:
        print(e)
    service_core = ServiceWorker(redisClient=redisClient, exploreCore=explore)
    try:
        service_core.loop()
    except Exception as e:
        print(e)
