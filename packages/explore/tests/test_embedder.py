from blinders.explore_core.embedder import Embedder
import unittest
from numpy import dot
from numpy.linalg import norm

from blinders.explore_core.types import MatchInfo


class TestEmbedder(unittest.TestCase):
    model: Embedder

    def __init__(self, methodName: str = "runTest") -> None:
        super().__init__(methodName)
        self.model = Embedder()

    def test_embed(self):
        matchInfo = MatchInfo(
            firebaseUID="firebaseUID",
            name="Hnimtadd",
            gender="male",
            major="software developer",
            native="vietnamese",
            country="viet nam",
            learnings=["english", "spanish"],
            interests=["coding", "music", "football"],
            userID="userID",
            age=22,
        )
        embed = self.model.embed(matchInfo)
        embed1 = self.model.embed(matchInfo)
        distance = dot(embed, embed1) / (norm(embed) * norm(embed1))
        print("=>distance", distance)
        self.assertEqual(len(embed), 384)
        self.assertEqual(len(embed1), 384)
        self.assertEqual(distance, 1.0)
        self.assertEqual(embed, embed1)


if __name__ == "__main__":
    unittest.main()
