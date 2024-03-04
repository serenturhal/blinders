import datetime
import random
import string
# migrate migrate users from firestore collection to mongo db.
import datetime
from firebase_admin import firestore, credentials, initialize_app
from pymongo import MongoClient
from redis.client import Redis
import os
import dotenv

from blinders.explore_core.main import Explore, Embedder, MatchInfo

genders = ["male", "female"]

majors = [
    "software engineer",
    "graphic designer",
    "chef",
    "police officer",
    "accountant",
    "writer",
    "electrician",
    "nurse",
    "student",
    "teacher",
    "doctor",
    "driver",
    "solider",
    "athlete",
]
langs = [
    "mandarin",
    "arabic",
    "russian",
    "german",
    "japanese",
    "portuguese",
    "italian",
    "vietnamese",
    "chinese",
    "english",
    "spanish",
]
interests = [
    "reading",
    "traveling",
    "photography",
    "gardening",
    "cooking",
    "painting",
    "music",
    "fitness",
    "writing",
    "football",
    "swimming",
    "coding",
    "running",
]
countries = [
    "us",
    "ca",
    "gb",
    "de",
    "au",
    "br",
    "in",
    "za",
    "vn",
    "cn",
    "fr",
    "jp",
]
matchCol = "matches"
userCol = "users"


def random_user(idx: int) -> dict:
    time_now = datetime.datetime.now()
    return {
        "firebaseUID": "".join(
            random.choices(string.ascii_lowercase + string.digits, k=10) + [str(idx)]
        ),
        "imageURL": "".join(
            random.choices(string.ascii_lowercase + string.digits, k=10) + [str(idx)]
        ),
        "name": "".join(
            random.choices(string.ascii_lowercase + string.digits, k=10) + [str(idx)]
        ),
        "friends": [],
        "createdAt": time_now,
        "updatedAt": time_now,
    }


def random_match_profile(user_id: str, name: str, mongo_user_id: str) -> MatchInfo:
    return MatchInfo(
        user_id,
        name,
        genders[random.randint(0, 1)],
        majors[random.randint(0, len(majors) - 1)],
        langs[random.randint(0, len(langs) - 1)],
        countries[random.randint(0, len(countries) - 1)],
        random.sample(langs, k=random.randint(1, 5)),
        random.sample(interests, k=random.randint(1, 5)),
        mongo_user_id,
        random.randint(10, 50),
    )

userColName = "users"
matchColName ="matches"

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
        mongo_client = MongoClient(mongoURL)
        db = mongo_client[os.getenv("MONGO_DATABASE", "Default")]
        match_col = db[matchColName]
        embedder = Embedder()
        redis_client = Redis(host="localhost", port=6379)
        explore = Explore(redis_client, embedder, match_col)

        num_Mock = 10000
        for idx in range(num_Mock):
            doc = random_user(idx)
            mongoUser = db[userColName].insert_one(doc)
            info = random_match_profile(doc.get("firebaseUID"), doc.get("name"), mongoUser.inserted_id)
            match_col.insert_one(info.__dict__)
            explore.add_user_embed(info)
    except Exception as e:
        raise e
