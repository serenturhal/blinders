from pymongo.collection import Collection
from redis.client import Redis
from blinders.explore_core.embedder import Embedder
from blinders.explore_core.types import MatchInfo
from blinders.explore_core.utils import CreateRedisMatchKey
from redis.commands.search.field import (
    TextField,
    VectorField,
)
from redis.commands.search.indexDefinition import IndexDefinition, IndexType


class Explore(object):
    redisClient: Redis
    embedder: Embedder
    matchCol: Collection
    vector_dimension = 384

    def __init__(
        self,
        RedisClient: Redis,
        Embbeder: Embedder,
        matchCol: Collection,
    ) -> None:
        if not RedisClient.ping():
            raise Exception("cannot ping to redis")

        self.redisClient = RedisClient
        self.embedder = Embbeder
        self.matchCol = matchCol
        self.initIndex()

    def initIndex(self):
        try:
            schema = (
                TextField("$.id", no_stem=True, as_name="id"),
                VectorField(
                    "$.embed",
                    "HNSW",
                    {
                        "TYPE": "FLOAT32",
                        "DIM": self.vector_dimension,
                        "DISTANCE_METRIC": "L2",
                    },
                    as_name="embed",
                ),
            )
            definition = IndexDefinition(prefix=["match:"], index_type=IndexType.JSON)
            res = self.redisClient.ft("idx:match_vss").create_index(
                fields=schema, definition=definition
            )
            print(res)

        except Exception as e:
            # maybe index defined
            print("error=>", e)

    def AddUserMatch(self, info: MatchInfo) -> None:
        embed = self.embedder.embed(info)
        self.redisClient.json().set(
            CreateRedisMatchKey(info.firebaseUID),
            "$",
            {
                "embed": embed,
                "id": info.firebaseUID,
            },
        )

        print(info)
        res = self.matchCol.insert_one(info.__dict__)
        if res.acknowledged:
            print("inserted match")
