# migrate migrate users from firestore collection to mongo db.
import datetime
from firebase_admin import firestore, credentials, initialize_app
from pymongo import MongoClient
import os
from dotenv import load_dotenv


load_dotenv()

userColName = "users"


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

    except Exception as e:
        raise e
