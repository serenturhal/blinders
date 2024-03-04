from blinders.explore_core.types import MatchInfo
from sentence_transformers import SentenceTransformer


class Embedder(object):
    model: SentenceTransformer

    def __init__(self, model_name: str = "all-MiniLM-L6-v2"):
        self.model = SentenceTransformer(model_name)

    def embed(self, info: MatchInfo) -> list[float]:
        embed_string = str(info)
        embeddings = self.model.encode([embed_string])
        return [float(v) for v in embeddings[0]]
