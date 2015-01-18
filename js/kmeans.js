var fs = require('fs');

var pointsFile = process.argv[2];
var seedsFile = process.argv[3];

function zip(lhs, rhs) {
  return lhs.map(function(l,i) { return [l, rhs[i]]; });
}

function Vec(pts) {
  this.pts = pts;
}

Vec.prototype.dot = function(that) {
  return zip(this.pts, that.pts)
    .map(function (x) { return x[0] * x[1]; })
    .reduce(function (a, b) { return a + b; });
}

Vec.prototype.add = function(that) {
  return new Vec(
    zip(this.pts, that.pts)
    .map(function (x) { return x[0] + x[1]; }));
}

Vec.prototype.sub = function(that) {
  return new Vec(
    zip(this.pts, that.pts)
    .map(function (x) { return x[0] - x[1]; }));
}

Vec.prototype.mult = function(a) {
  return new Vec(this.pts.map(function (e) { return e * a; }));
}

Vec.prototype.norm = function() {
  return this.dot(this);
}

function parsePoint(s) {
  return new Vec(s.trim().split(" ").map(parseFloat));
}

function bootstrap(pointsData, seedsData) {
  var points = pointsData.trim().split("\n").map(parsePoint);
  var seeds = seedsData.trim().split("\n").map(parsePoint);
  kmeans(points, seeds, 20).forEach(function (e) { console.log(e); });
}

function argmin(elems) {
  var minSoFar = elems[0];
  var minIdxSoFar = 0;
  for (var i = 1; i < elems.length; i++) {
    if (elems[i] < minSoFar) {
      minSoFar = elems[i];
      minIdxSoFar = i;
    }
  }
  return minIdxSoFar;
}

function which(pt, clusters) {
  return argmin(clusters.map(function (c) { return pt.sub(c).norm(); }));
}

function centroid(pts) {
  return pts
    .reduce(function (a, b) { return a.add(b); })
    .mult(1. / pts.length);
}

function kmeans(points, seeds, iters) {
  var clusters = seeds.slice(0); // clone
  for (var i = 0; i < iters; i++) {
    var assignments = [];
    for (var j = 0; j < clusters.length; j++) {
      assignments.push([]);
    }
    points.forEach(function (pt) {
      assignments[which(pt, clusters)].push(pt);
    });
    clusters = zip(assignments, clusters)
      .map(function (e) { return e[0].length === 0 ? e[1] : centroid(e[0]); });
  }
  return clusters
}

fs.readFile(pointsFile, 'ascii', function (err, pointsData) {
  fs.readFile(seedsFile, 'ascii', function (err, seedsData) {
    bootstrap(pointsData, seedsData);
  });
});
