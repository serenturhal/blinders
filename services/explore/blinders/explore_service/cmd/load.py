# load mock/json file to mongodb
import datetime
import json
import os

import dotenv
import pymongo
from bson.objectid import ObjectId
from redis.client import Redis

from blinders.explore_core.embedder import Embedder
from blinders.explore_core.main import Explore, MatchInfo

matchColName = "matches"
userColName = "users"
genders = ["male", "female"]

if __name__ == "__main__":
    dotenv.load_dotenv()
    mongoURL = "mongodb://{}:{}@{}:{}/{}".format(
        os.getenv("MONGO_USERNAME"),
        os.getenv("MONGO_PASSWORD"),
        os.getenv("MONGO_HOST"),
        os.getenv("MONGO_PORT"),
        os.getenv("MONGO_DATABASE"),
    )
    mongoClient = pymongo.MongoClient(mongoURL)
    db = mongoClient.get_database(os.getenv("MONGO_DATABASE", "Default"))
    match_col = db.get_collection(matchColName)
    user_col = db.get_collection(userColName)

    embedder = Embedder()
    redis_client = Redis(host="localhost", port=6379)
    explore = Explore(redis_client, embedder, match_col)

    with open("mock/users.json", "r") as f:
        users = json.load(f)
        for user in users:
            now = datetime.datetime.now()
            user["createdAt"] = now
            user["updatedAt"] = now
            user["_id"] = ObjectId(user["_id"])
            user_col.insert_one(user)

    with open("mock/matches.json", "r") as f:
        matches = json.load(f)
        for match in matches:
            match["userID"] = ObjectId(match["userID"])
            match["_id"] = ObjectId(match["_id"])
            match_col.insert_one(match)
            try:
                matchInfo = MatchInfo(
                    str(match.get("userID")),
                    match.get("name"),
                    match.get("gender"),
                    match.get("major"),
                    match.get("native"),
                    match.get("country"),
                    match.get("learnings"),
                    match.get("interests"),
                    match.get("age"),
                )
                explore.add_user_embed(matchInfo)
            except Exception as e:
                print(e)
                break
