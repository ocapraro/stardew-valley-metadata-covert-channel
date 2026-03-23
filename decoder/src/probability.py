from math import factorial
from collections import Counter

def unique_permutation_count(items):
    counts = Counter(items)
    n = len(items)

    denom = 1
    for c in counts.values():
        denom *= factorial(c)

    return factorial(n) // denom


def distribute_balls_unlabeled(balls: int, holes: int):
    """
    Return all ways to put `balls` identical balls into `holes` unlabeled holes,
    with at least 1 ball in each hole.
    
    Each result is a tuple of length `holes` in non-increasing order.
    Example: balls=6, holes=3 -> (4, 1, 1), (3, 2, 1), (2, 2, 2)
    """
    if balls < holes:
        # Can't give at least 1 ball per hole
        return []

    results = []

    def backtrack(remaining: int, holes_left: int, max_next: int, current: list[int]):
        # remaining: balls left to place
        # holes_left: holes still to fill
        # max_next: max balls allowed in the next hole (for non-increasing order)
        if holes_left == 0:
            if remaining == 0:
                results.append(tuple(current))
            return

        # We need at least 1 ball in each remaining hole
        max_can_put = min(max_next, remaining - (holes_left - 1))

        # Place x balls in the next hole
        for x in range(max_can_put, 0, -1):
            current.append(x)
            backtrack(remaining - x, holes_left - 1, x, current)
            current.pop()

    backtrack(balls, holes, balls, [])
    return results