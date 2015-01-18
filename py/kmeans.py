import math
import sys
import itertools as it


class Vec(object):
    def __init__(self, pts):
        self.pts = pts

    def dot(self, that):
        return sum(x * y for x, y in zip(self.pts, that.pts))

    def norm(self):
        return math.sqrt(self.dot(self))

    def __add__(self, that):
        return Vec([x + y for x, y in zip(self.pts, that.pts)])

    def __sub__(self, that):
        return Vec([x - y for x, y in zip(self.pts, that.pts)])

    def __mul__(self, that):
        """Scalar multiplication

        """
        return Vec([that * x for x in self.pts])

    def __str__(self):
        return str(self.pts)

    def __repr__(self):
        return repr(self.pts)


def argmin(elems):
    minSoFar, minIdxSoFar = elems[0], 0
    for i, e in it.islice(enumerate(elems), 1, None):
        if e < minSoFar:
            minSoFar, minIdxSoFar = e, i
    return minIdxSoFar


def kmeans(points, seeds, iters):

    def centroid(pts):
        return reduce(lambda x, y: x + y, pts) * (1. / len(pts))

    clusters = list(seeds)
    for t in xrange(iters):
        assignments = [[] for _ in clusters]
        for pt in points:
            dists = map(lambda c: (c - pt).norm(), clusters)
            assignments[argmin(dists)].append(pt)
        clusters = [c if not elems else centroid(elems)
                    for c, elems in zip(clusters, assignments)]
    return clusters


def parse(fname):
    with open(fname, 'r') as fp:
        for l in fp:
            yield Vec(map(float, l.strip().split()))


def main():
    pointsFile, seedsFile = sys.argv[1:]
    points = list(parse(pointsFile))
    seeds = list(parse(seedsFile))
    clusters = kmeans(points, seeds, 20)
    for c in clusters:
        print c


if __name__ == '__main__':
    main()
