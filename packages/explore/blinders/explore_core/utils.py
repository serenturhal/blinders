def CreateRedisMatchKey(userID: str) -> str:
    return f"match:{userID}"
