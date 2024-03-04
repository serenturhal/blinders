# migrate migrate users from firestore collection to mongo db.
import datetime
from firebase_admin import firestore, credentials, initialize_app
from pymongo import MongoClient
import os
import dotenv

userColName = "users"

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

        creds = credentials.Certificate("./firebase.admin.development.json")
        app = initialize_app(creds)
        firestoreClient = firestore.client(app)
        userDocs = firestoreClient.collection("Users").stream()

        for user_doc in userDocs:
            doc = user_doc.to_dict()
            if doc is None:
                userDocs.close()
                break

            userID = doc.get("firebaseUid")
            name = doc.get("name")
            image_url = doc.get("imageUrl")
            friends_firebase_uid = doc.get("friends")

            if (
                userID is None
                or name is None
                or image_url is None
                or friends_firebase_uid is None
            ):
                continue

            now = datetime.datetime.now()
            mongoUser = db[userColName].insert_one(
                {
                    "firebaseUID": userID,
                    "imageURL": image_url,
                    "friends": friends_firebase_uid,
                    "createdAt": now,
                    "updatedAt": now,
                }
            )

    except Exception as e:
        raise e
