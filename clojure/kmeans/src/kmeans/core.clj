(ns kmeans.core 
    (:require [clojure.math.numeric-tower :as math]
              [clojure.string :as str]))

(defn vecadd [lhs rhs] (mapv + lhs rhs))

(defn vecsub [lhs rhs] (mapv - lhs rhs))

(defn vecscale [v s] (mapv #(* % s) v))

(defn dot [lhs rhs] (reduce + (mapv * lhs rhs)))
 
(defn norm [v] (math/sqrt (dot v v)))

(defn centroid [points] (vecscale (reduce vecadd points) (/ 1. (count points))))

(defn dists [pt centers] (map #(norm (vecsub pt %)) centers))

; wow this is ugly, what's the more idiomatic way to do this?
(defn argmin [elems] 
  (first 
    (reduce #(if (< %2 (%1 1)) 
                [(%1 2) %2 (+ 1 (%1 2))] 
                [(%1 0) (%1 1) (+ 1 (%1 2))]) 
            [0 (first elems) 1] 
            (rest elems))))

(defn transpose [x] (apply map list x))

(defn assignment [pt centers] 
  (let [idx (argmin (dists pt centers))] 
    (map #(if (== idx %) (list pt) '()) (range (count centers)))))

(defn kmeans-iter [pts centers] 
  (map #(if (empty? %1) %2 (centroid %1))
    (map #(mapcat identity %) (transpose (map #(assignment % centers) pts)))
    centers))

(defn kmeans [pts centers iters]
  (if (== iters 0) centers (recur pts (kmeans-iter pts centers) (dec iters))))

(defn parse-point [s] (mapv read-string (str/split (str/trim s) #" ")))

(defn parse [s] (map parse-point (str/split-lines (str/trim s))))

(defn parse-file [fname] (parse (slurp fname)))

(defn -main [& args] 
  (let [pointsFile (nth args 0) 
        seedsFile (nth args 1)] 
    (doseq [elem (kmeans (parse-file pointsFile) (parse-file seedsFile) 20)]
      (println elem))))
