import datetime
from blinders.explore_core.main import Explore
from blinders.explore_core.embedder import Embedder
from redis.client import Redis
from firebase_admin import firestore, credentials, initialize_app
from pymongo import MongoClient
import os
from dotenv import load_dotenv


load_dotenv()

matchColName = "matchs"
userColName = "users"
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

if __name__ == "__main__":
    try:
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

        embbeder = Embedder()
        redisClient = Redis(host="localhost", port=6379)
        explore = Explore(redisClient, embbeder, matchCol)
        explore.initIndex()

        creds = credentials.Certificate("./firebase.admin.development.json")
        app = initialize_app(creds)
        firestoreClient = firestore.client(app)
        userDocs = firestoreClient.collection("Users").stream()

        for userDoc in userDocs:
            doc = userDoc.to_dict()
            if doc is None:
                userDocs.close()
                break

            userID = doc.get("firebaseUid")
            name = doc.get("name")
            imageURL = doc.get("imageUrl")
            friendsFirebaseUID = doc.get("friends")

            if (
                userID is None
                or name is None
                or imageURL is None
                or friendsFirebaseUID is None
            ):
                continue

            # cur = matchCol.find_one(filter={"firebaseUID": userID})
            # if cur is None:
            #     continue

            now = datetime.datetime.now()
            mongoUser = db[userColName].insert_one(
                {
                    "firebaseUID": userID,
                    "imageURL": imageURL,
                    "friends": friendsFirebaseUID,
                    "createdAt": now,
                    "updatedAt": now,
                }
            )

            # matchUser = MatchInfo(
            #     userID,
            #     name,
            #     genders[random.randint(0, 1)],
            #     majors[random.randint(0, len(majors) - 1)],
            #     langs[random.randint(0, len(langs) - 1)],
            #     countries[random.randint(0, len(countries) - 1)],
            #     random.sample(langs, k=random.randint(1, 5)),
            #     random.sample(interests, k=random.randint(1, 5)),
            #     mongoUser.inserted_id,
            #     random.randint(10, 50),
            # )
            # explore.AddUserMatch(matchUser)

    except Exception as e:
        raise e
