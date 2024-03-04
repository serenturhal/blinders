from blinders.explore_core.main import Explore
from blinders.explore_core.embedder import Embedder
from redis.client import Redis
import pymongo
import os
from dotenv import load_dotenv
import json


load_dotenv()

matchColName = "matches"
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
    mongoClient = pymongo.MongoClient(mongoURL)
    db = mongoClient[os.getenv("MONGO_DATABASE", "Default")]
    match_col = db[matchColName]
    userCol = db[userColName]

    embedder = Embedder()
    redisClient = Redis(host="localhost", port=6379)
    explore = Explore(redisClient, embedder, match_col)
    usersCur = userCol.find({})
    matchesCur = match_col.find({})
    users = []
    matches = []

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
            curr = matchesCur.next()
        except Exception as e:
            print(e)
            break
        assert isinstance(curr, dict)
        matches.append(curr)
    with open("mock/matches.json", "w") as f:
        json.dump(matches, f, default=str)
