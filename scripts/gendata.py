"""gendata.py - Generate random data for k-means

"""

import numpy as np


def main():
    # Toy example for now:

    # Generate k clusters of gaussians in d dimension
    k, d = 4, 10
    centers = np.random.multivariate_normal(mean=np.zeros(d), cov=2.*np.eye(d), size=k)

    # Now generate 100 points from each cluster
    n = 100
    pts = np.vstack([np.random.multivariate_normal(mean=c, cov=0.5*np.eye(d), size=n) for c in centers])

    # Permute the points
    pts = pts[np.random.permutation(pts.shape[0])]

    # Generate k starting points
    seeds = np.random.multivariate_normal(mean=np.zeros(d), cov=2.*np.eye(d), size=k)

    # Write two files
    np.savetxt('data/points.txt', pts)
    np.savetxt('data/seeds.txt', seeds)


if __name__ == '__main__':
    main()
