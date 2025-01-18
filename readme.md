# RANSAC
During my winter 2023 school term I worked on developing several modified implementations of the RANSAC (Random Sample Consensus) algorithm. These modified implementations are written in Java and Go respectively and aim to identify the largest plane that can be formed within a point cloud dataset. That is, the plane within the point cloud that contains the most points.

> Note that the Go implementation uses goroutines to form a concurrent pipeline that speeds up the execution of the algorithm.

# Running the Application
Both the Java and Go implementations expect a **confidence** and **percentage** value to be supplied.

> The **confidence** value controls the number of iterations the algorithm runs to guarantee a certain confidence that the dominant plane has been found.
>
> The **percentage** value is the percentage of points that support the dominant plane.

The program expects as input a *csv* file with the following format:
```csv
x,y,z
```
Every row after the *x,y,z* header should be a point in the point cloud dataset.

## Java
The Java implementation can be run from the **PlaneRANSAC.java** file.

## Go
The Go implementation can be run from the planeRANSAC.go file.