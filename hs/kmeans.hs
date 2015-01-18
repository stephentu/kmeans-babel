import Data.List
import System.Environment
import System.IO
import Control.Monad (forM_)

data Vec a = Vec [a] deriving (Show)
data Dataset a = Dataset [Vec a]
data Clustering a = Clustering [Vec a]

vecDot :: Num a => Vec a -> Vec a -> a
vecDot (Vec x) (Vec y) = foldl (\acc t -> acc + (fst t) * (snd t)) 0 (zip x y)

vecAdd :: Num a => Vec a -> Vec a -> Vec a
vecAdd (Vec x) (Vec y) = Vec $ map (\t -> (fst t) + (snd t)) $ zip x y

vecSub :: Num a => Vec a -> Vec a -> Vec a
vecSub (Vec x) (Vec y) = Vec $ map (\t -> (fst t) - (snd t)) $ zip x y

-- x * a, where x is vec and a is scalar
vecScale :: Num a => Vec a -> a -> Vec a
vecScale (Vec x) s = Vec $ map (\t -> s * t) x

vecNorm :: RealFloat a => Vec a -> a
vecNorm x = sqrt $ vecDot x x

-- given a list of vectors, compute the average point
centroid :: RealFloat a => [Vec a] -> Vec a
centroid l@(x:xs) = vecScale (foldl vecAdd x xs) (1 / (fromIntegral $ length l))

dists :: RealFloat a => Vec a -> [Vec a] -> [a]
dists x pts = map (\t -> vecNorm $ vecSub x t) pts

argmin :: Ord a => [a] -> Int
argmin xs = snd (foldl update ((head xs, 0)) (zip (tail xs) [1..]))
  where update s@(vbest, ibest) t@(vcur, icur) = if vcur < vbest then t else s

assign0 :: RealFloat a => Vec a -> [Vec a] -> [[Vec a]]
assign0 pt clusters = r
  where
    d = dists pt clusters
    i = argmin d
    r = map (\t -> if t == i then [pt] else []) [0 .. (length clusters)]

assign :: RealFloat a => [Vec a] -> [Vec a] -> [[Vec a]]
assign pts clusters = map concat $ transpose a
  where
    a = map (\pt -> assign0 pt clusters) pts

-- given a list of pts and a current cluster assignment, return an updated
-- cluster assignment
kmeansIter :: RealFloat a => [Vec a] -> [Vec a] -> [Vec a]
kmeansIter pts clusters = newClusters
  where
    newClusters = map (\x -> if null $ fst x then snd x else centroid $ fst x) $ zip (assign pts clusters) clusters

kmeans0 :: RealFloat a => [Vec a] -> [Vec a] -> Int -> [Vec a]
kmeans0 _ clustering 0 = clustering
kmeans0 pts clusters n = kmeans0 pts (kmeansIter pts clusters) (n - 1)

kmeans :: RealFloat a => Dataset a -> Clustering a -> Int -> Clustering a
kmeans (Dataset pts) (Clustering clusters) n = Clustering $ kmeans0 pts clusters n

-- http://stackoverflow.com/questions/4503958/what-is-the-best-way-to-split-a-string-by-a-delimiter-functionally
splitBy delimiter = foldr f [[]]
  where f c l@(x:xs) | c == delimiter = []:l
                     | otherwise = (c:x):xs

parse :: [String] -> [Vec Float]
parse lines = map (\t -> Vec $ map read (splitBy ' ' t)) lines

prog :: [String] -> [String] -> [Vec Float]
prog rawPoints rawSeeds = finalClusters
  where
    dataset = Dataset $ parse rawPoints
    origClusters = Clustering $ parse rawSeeds
    Clustering finalClusters = kmeans dataset origClusters 20

main = do
  args <- getArgs
  case args of
    [pointsFile, seedsFile] -> do
        rawPoints <- readFile pointsFile
        rawSeeds <- readFile seedsFile
        forM_ (prog (lines rawPoints) (lines rawSeeds)) $ \p -> do
           print p
    _ -> error "USAGE: ./kmeans pointsFile seedsFile"
