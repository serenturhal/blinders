from blinders.explore_core.main import Explore
from blinders.explore_core.embedder import Embedder
from redis.client import Redis
from pymongo import MongoClient
import os
from dotenv import load_dotenv
import json


load_dotenv()

matchColName = "matchs"
userColName = "users"
genders = ["male", "female"]

if __name__ == "__main__":
    mongoURL = "mongodb://{}:{}@{}:{}/{}".format(
        os.getenv("MONGO_USERNAME"),
        os.getenv("MONGO_PASSWORD"),
        os.getenv("MONGO_HOST"),
        os.getenv("MONGO_PORT"),
        os.getenv("MONGO_DATABASE"),
    )
    print(mongoURL)
    mongoClient = MongoClient(mongoURL)
    db = mongoClient[os.getenv("MONGO_DATABASE", "Default")]
    matchCol = db[matchColName]
    userCol = db[userColName]

    embbeder = Embedder()
    redisClient = Redis(host="localhost", port=6379)
    explore = Explore(redisClient, embbeder, matchCol)
    usersCur = userCol.find({})
    matchsCur = matchCol.find({})
    print(usersCur.max)
    users = []
    matchs = []

    while True:
        try:
            curr = usersCur.next()
        except Exception as e:
            print(e)
            break
        assert isinstance(curr, dict)
        users.append(curr)

    with open("mock/users.json", "w") as f:
        json.dump(users, f, default=str)
    while True:
        try:
            curr = matchsCur.next()
        except Exception as e:
            print(e)
            break
        assert isinstance(curr, dict)
        matchs.append(curr)
    with open("mock/matchs.json", "w") as f:
        json.dump(matchs, f, default=str)
