# load mock/json file to mongodb
from blinders.explore_core.main import Explore, MatchInfo
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
    mongoClient = MongoClient(mongoURL)
    db = mongoClient.get_database(os.getenv("MONGO_DATABASE", "Default"))
    matchCol = db.get_collection(matchColName)
    userCol = db.get_collection(userColName)

    embbeder = Embedder()
    redisClient = Redis(host="localhost", port=6379)
    explore = Explore(redisClient, embbeder, matchCol)
    usersCur = userCol.find({})
    matchsCur = matchCol.find({})

    with open("mock/users.json", "r") as f:
        users = json.load(f)
        for user in users:
            userCol.insert_one(user)

    with open("mock/matchs.json", "r") as f:
        matchs = json.load(f)
        for match in matchs:
            matchCol.insert_one(match)
            try:
                matchInfo = MatchInfo(
                    match.get("firebaseUID"),
                    match.get("name"),
                    match.get("gender"),
                    match.get("major"),
                    match.get("native"),
                    match.get("country"),
                    match.get("learnings"),
                    match.get("interests"),
                    match.get("userID"),
                    match.get("age"),
                )
                explore.AddUserEmbed(matchInfo)
            except Exception as e:
                print(e)
                break
