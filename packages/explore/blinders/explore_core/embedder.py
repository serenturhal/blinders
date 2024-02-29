from blinders.explore_core.types import MatchInfo
from sentence_transformers import SentenceTransformer


class Embedder(object):
    model: SentenceTransformer

    def __init__(self, modelName: str = "all-MiniLM-L6-v2"):
        self.model = SentenceTransformer(modelName)

    def embed(self, info: MatchInfo) -> list[float]:
        embedString = info.firebaseUID
        embeddings = self.model.encode([embedString])
        return [float(v) for v in embeddings[0]]
