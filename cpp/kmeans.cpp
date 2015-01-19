#include <vector>
#include <iostream>
#include <sstream>
#include <fstream>
#include <string>
#include <cmath>
#include <cstdlib>

template <typename Scalar>
class vec {
public:

  vec() = default;
  vec(size_t n) : p_(n) {}

  vec(const std::vector<Scalar> &p) : p_(p) {}
  vec(std::vector<Scalar> &&p) : p_(std::move(p)) {}

  inline Scalar dot(const vec<Scalar>& that) const
  {
    Scalar acc(0);
    for (size_t i = 0; i < p_.size(); i++)
      acc += p_[i] * that.p_[i];
    return acc;
  }

  inline Scalar norm() const
  {
    return sqrt(dot(*this));
  }

  inline vec<Scalar> operator+(const vec<Scalar> &that) const
  {
    vec<Scalar> ret(p_);
    ret += that;
    return ret;
  }

  inline vec<Scalar> operator-(const vec<Scalar> &that) const
  {
    vec<Scalar> ret(p_);
    ret -= that;
    return ret;
  }

  inline vec<Scalar> operator*(Scalar s) const
  {
    vec<Scalar> ret(p_);
    ret *= s;
    return ret;
  }

  inline const std::vector<Scalar> & repr() const { return p_; }

  inline vec<Scalar> &operator+=(const vec<Scalar> &that) {
    for (size_t i = 0; i < that.p_.size(); i++)
      p_[i] += that.p_[i];
    return *this;
  }

  inline vec<Scalar> &operator-=(const vec<Scalar> &that) {
    for (size_t i = 0; i < that.p_.size(); i++)
      p_[i] -= that.p_[i];
    return *this;
  }

  inline vec<Scalar> &operator*=(Scalar s) {
    for (size_t i = 0; i < p_.size(); i++)
      p_[i] *= s;
    return *this;
  }

  inline size_t size() const { return p_.size(); }

private:
  std::vector<Scalar> p_;
};

template <typename Scalar>
static inline vec<Scalar>
operator*(const Scalar s, const vec<Scalar> &v)
{
  return v * s;
}

template <typename Scalar>
static inline std::ostream &
operator<<(std::ostream &o, const vec<Scalar> &v)
{
  o << "[";
  bool first = true;
  for (auto e : v.repr()) {
    if (first)
      first = false;
    else
      o << ", ";
    o << e;
  }
  o << "]";
  return o;
}

typedef std::vector<vec<double>> dataset;
typedef std::vector<vec<double>> clustering;

static std::vector<std::string>
Split(const std::string &s, char delim)
{
  std::stringstream ss(s);
  std::string elem;
  std::vector<std::string> elems;
  while (std::getline(ss, elem, delim))
    elems.emplace_back(elem);
  return elems;
}

static std::vector<double>
Parse(const std::vector<std::string> &tokens)
{
  std::vector<double> ret;
  for (const auto &p : tokens)
    ret.push_back(strtod(p.c_str(), nullptr));
  return ret;
}

static dataset
ParseFile(const std::string &fname)
{
  std::ifstream ifs(fname);
  std::string line;
  dataset d;
  while (std::getline(ifs, line)) {
    const auto tokens = Split(line, ' ');
    const auto v = Parse(tokens);
    d.emplace_back(std::move(v));
  }
  return d;
}

static size_t
Closest(const vec<double> &pt, const clustering &c)
{
  size_t idx = 0;
  double val = (pt - c.front()).norm();
  for (size_t i = 1; i < c.size(); i++) {
    const double d = (pt - c[i]).norm();
    if (d < val) {
      idx = i;
      val = d;
    }
  }
  return idx;
}

static vec<double>
Centroid(const std::vector<const vec<double> *> &pts)
{
  vec<double> r(pts.front()->size());
  for (auto px : pts)
    r += (*px);
  r *= 1. / float(pts.size());
  return r;
}

static clustering
KMeans(const dataset &d, const clustering &c, size_t niters)
{
  clustering res(c);
  for (size_t t = 0; t < niters; t++) {
    std::vector<std::vector<const vec<double> *>> assignments(res.size());
    for (const auto &pt : d)
      assignments[Closest(pt, res)].push_back(&pt);
    for (size_t i = 0; i < assignments.size(); i++) {
      if (assignments[i].empty())
        continue;
      res[i] = Centroid(assignments[i]);
    }
  }
  return res;
}

int
main(int argc, char **argv)
{
  const std::string pointsFile = argv[1];
  const std::string seedsFile = argv[2];

  const dataset d = ParseFile(pointsFile);
  const clustering c = ParseFile(seedsFile);

  for (const auto &p : KMeans(d, c, 20))
    std::cout << p << std::endl;

  return 0;
}
