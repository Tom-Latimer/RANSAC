// Tom Latimer, 300250278
package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Point3D struct {
	X float64
	Y float64
	Z float64
}
type Plane3D struct {
	A float64
	B float64
	C float64
	D float64
}
type Plane3DwSupport struct {
	Plane3D
	SupportSize int
}

// reads an XYZ file and returns a slice of Point3D
func ReadXYZ(filename string) []Point3D {
	f, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	outputSlice := make([]Point3D, 1)

	//skip x,y,z header
	scanner.Scan()

	//scan each line of the file
	for scanner.Scan() {
		//split the coordinate line and place it into a slice
		s := strings.Split(scanner.Text(), "\t")

		//convert coordinate to float64
		x, errX := strconv.ParseFloat(s[0], 64)

		if errX != nil {
			log.Fatal(errX)
		}

		//convert coordinate to float64
		y, errY := strconv.ParseFloat(s[1], 64)

		if errY != nil {
			log.Fatal(errY)
		}

		//convert coordinate to float64
		z, errZ := strconv.ParseFloat(s[2], 64)

		if errZ != nil {
			log.Fatal(errZ)
		}

		//add point to output slice
		outputSlice = append(outputSlice, Point3D{x, y, z})
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return outputSlice
}

// saves a slice of Point3D into an XYZ file
func SaveXYZ(filename string, points []Point3D) {

	//create file
	dat, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	defer dat.Close()

	//add header
	_, errI := dat.WriteString("x\ty\tz\n")
	if errI != nil {
		log.Fatal(errI)
	}

	//format a Point3D into string and add it to output file
	for _, point := range points {
		_, err := dat.WriteString(fmt.Sprintf("%f\t%f\t%f\n", point.X, point.Y, point.Z))
		if err != nil {
			log.Fatal(err)
		}
	}
}

// computes the distance between points p1 and p2
func (p1 *Point3D) GetDistance(p2 *Point3D) float64 {
	return math.Sqrt(math.Pow((p2.X-p1.X), 2) + math.Pow((p2.Y-p1.Y), 2) + math.Pow((p2.Z-p1.Z), 2))
}

// calcualtes the distance between a point and a plane
func (plane *Plane3D) CalcDistance(pt *Point3D) float64 {
	return math.Abs((plane.A*pt.X)+(plane.B*pt.Y)+(plane.C*pt.Z)+plane.D) / math.Sqrt(math.Pow(plane.A, 2)+math.Pow(plane.B, 2)+math.Pow(plane.C, 2))
}

// computes the plane defined by a set of 3 points
func GetPlane(points []Point3D) Plane3D {
	vector1 := Point3D{points[1].X - points[0].X, points[1].Y - points[0].Y, points[1].Z - points[0].Z}
	vector2 := Point3D{points[2].X - points[0].X, points[2].Y - points[0].Y, points[2].Z - points[0].Z}

	a := ((vector1.Y * vector2.Z) - (vector2.Y * vector1.Z))
	b := ((vector1.Z * vector2.X) - (vector1.X * vector2.Z))
	c := ((vector1.X * vector2.Y) - (vector1.Y * vector2.X))
	d := (a * points[0].X) + (b * points[0].Y) + (c * points[0].Z)

	return Plane3D{a, b, c, d}
}

// computes the number of required RANSAC iterations
func GetNumberOfIterations(confidence float64, percentageOfPointsOnPlane float64) int {
	return int(math.Log(1-confidence) / (math.Log(1 - math.Pow(percentageOfPointsOnPlane, 3))))
}

// computes the support of a plane in a set of points
func GetSupport(plane Plane3D, points []Point3D, eps float64) Plane3DwSupport {
	support := Plane3DwSupport{plane, 0}
	for _, point := range points {
		if dist := plane.CalcDistance(&point); dist < eps {
			support.SupportSize++
		}
	}
	return support
}

// extracts the points that supports the given plane
// and returns them as a slice of points
func GetSupportingPoints(plane Plane3D, points []Point3D, eps float64) []Point3D {
	outputSlice := make([]Point3D, 1)
	for _, point := range points {
		if dist := plane.CalcDistance(&point); dist < eps {
			outputSlice = append(outputSlice, point)
		}
	}
	return outputSlice
}

// creates a new slice of points in which all points
// belonging to the plane have been removed
func RemovePlane(plane Plane3D, points []Point3D, eps float64) []Point3D {
	outputSlice := make([]Point3D, 1)
	for _, point := range points {
		if dist := plane.CalcDistance(&point); dist >= eps {
			outputSlice = append(outputSlice, point)
		}
	}
	return outputSlice
}

// generates random points from the input slice and
// sends them to the returned channel
func RandPointGenerator(input []Point3D, stop <-chan bool) <-chan Point3D {
	pointStream := make(chan Point3D)
	rand.Seed(time.Now().Unix())
	go func() {
		defer close(pointStream)
		for {
			select {
			case <-stop:
				return
			case pointStream <- input[rand.Intn(len(input))]:
			}
		}
	}()
	return pointStream
}

// aggregates 3 random points form RandPointGenerator and send them
// to the returned channel
func TripletPointGenerator(points <-chan Point3D, stop <-chan bool) <-chan [3]Point3D {
	var output [3]Point3D
	i := 0
	outputStream := make(chan [3]Point3D)

	go func() {
		defer close(outputStream)
		for {
			select {
			case <-stop:
				return
			case newPoint := <-points:
				output[i] = newPoint
				i++

				//if received 3 points, output triplet
				if i == 3 {
					outputStream <- output
					i = 0

				}

			}
		}
	}()
	return outputStream
}

// counts the number of triplet arrays that have passed through
// sends a stop signal to other goroutines when n triplets have
// passed through
func TakeN(pArr <-chan [3]Point3D, stop chan<- bool, n int) <-chan [3]Point3D {
	i := 0
	outputStream := make(chan [3]Point3D)
	go func() {
		defer close(outputStream)
		for {
			select {
			case buffer := <-pArr:
				//pass the triplets through
				outputStream <- buffer
				i++
				if i == n {
					//send stop stignal
					stop <- true
					return
				}
			}
		}
	}()
	return outputStream
}

// takes triplets of points and sends planes on its
// output channel
func PlaneEstimator(pArr <-chan [3]Point3D) <-chan Plane3D {
	planeStream := make(chan Plane3D)
	go func() {
		defer close(planeStream)
		for points := range pArr {
			planeStream <- GetPlane(points[:])
		}
	}()
	return planeStream
}

// takes planes from its input channel, calculates the supporting number of points
// for that plane and outputs it to its output channel
func SuppPointFinder(planeStream <-chan Plane3D, pc []Point3D, eps float64) chan Plane3DwSupport {
	supportStream := make(chan Plane3DwSupport)
	go func() {
		defer close(supportStream)
		for plane := range planeStream {
			supportStream <- GetSupport(plane, pc, eps)
		}
	}()
	return supportStream
}

// multiplexes the goroutines in a given slice into a single channel
// that channel is then returned
func FanIn(inputStreams []chan Plane3DwSupport) <-chan Plane3DwSupport {
	var wg sync.WaitGroup

	//the channel that will be multipleced to
	outputStream := make(chan Plane3DwSupport)

	//anonymous function that takes planeSupports from a
	//SuppPointFinder goroutine and adds it to the output channel
	output := func(planeSupport <-chan Plane3DwSupport) {
		for i := range planeSupport {
			outputStream <- i
		}
		wg.Done()
	}
	wg.Add(len(inputStreams))
	for _, support := range inputStreams {
		//activates the goroutine to take the supporting points
		//and multiplex it to the output channel
		go output(support)
	}

	//waits until the multiplexing is complete to close the channel
	go func() {
		wg.Wait()
		close(outputStream)
	}()
	return outputStream
}

// receives planes from its input channel and stores the dominant plane
// returns the dominant plane to its output channel
func domPlaneFinder(supportStream <-chan Plane3DwSupport) <-chan Plane3DwSupport {
	var currentBestSupport Plane3DwSupport
	//set current best support to 0
	currentBestSupport = Plane3DwSupport{Plane3D{0, 0, 0, 0}, 0}

	output := make(chan Plane3DwSupport)
	go func() {

		defer close(output)
		for p := range supportStream {

			if p.SupportSize > currentBestSupport.SupportSize {
				//fmt.Println(p.SupportSize, currentBestSupport.SupportSize)
				currentBestSupport = p

			}

		}

		output <- (currentBestSupport)
	}()
	return output
}

func main() {
	//read arguments
	inputArgs := os.Args[1:5]

	inputFile := inputArgs[0]

	//convert confidence to float64
	confidence, errCon := strconv.ParseFloat(inputArgs[1], 64)
	if errCon != nil {
		log.Fatal(errCon)
	}

	//convert plane percentage to float64
	percentage, errPerc := strconv.ParseFloat(inputArgs[2], 64)
	if errPerc != nil {
		log.Fatal(errPerc)
	}

	//convert epsilon value to float64
	eps, errEps := strconv.ParseFloat(inputArgs[3], 64)
	if errEps != nil {
		log.Fatal(errEps)
	}

	pointCloud := ReadXYZ(inputFile)

	//format the file name to remove file extension
	//used later for output file name
	inputFile = strings.ReplaceAll(inputFile, ".xyz", "")

	//get number of iterations to satisfy confidence value
	numIterations := GetNumberOfIterations(confidence, percentage)

	for i := 0; i < 3; i++ {
		bestSupport := Plane3DwSupport{}
		stop := make(chan bool)

		//start concurrent pipeline
		pointGen := RandPointGenerator(pointCloud, stop)
		tripPointGen := TripletPointGenerator(pointGen, stop)
		takeN := TakeN(tripPointGen, stop, numIterations)
		estimator := PlaneEstimator(takeN)

		pointfinder := make([]chan Plane3DwSupport, 20)

		//fan out support finder goroutines
		for j := 0; j < 20; j++ {
			pointfinder[j] = SuppPointFinder(estimator, pointCloud, eps)
		}

		multiplexedOut := FanIn(pointfinder)

		//get dominant plane
		bestSupport = <-domPlaneFinder(multiplexedOut)

		//save plane points to file
		SaveXYZ(fmt.Sprintf(inputFile+"_p%d.xyz", i+1), GetSupportingPoints(bestSupport.Plane3D, pointCloud, eps))

		pointCloud = RemovePlane(bestSupport.Plane3D, pointCloud, eps)
	}
	SaveXYZ((inputFile + "_p0.xyz"), pointCloud)
}

//Code Used for testing

/*
func main() {
	//read arguments
	inputArgs := os.Args[1:5]

	output := make([]timeOutput, 1)

	for num := 10; num < 51; num++ {
		avg := 0.0
		for l := 0; l < 5; l++ {
			start := time.Now()

			inputFile := inputArgs[0]

			//convert confidence to float64
			confidence, errCon := strconv.ParseFloat(inputArgs[1], 64)
			if errCon != nil {
				log.Fatal(errCon)
			}

			//convert plane percentage to float64
			percentage, errPerc := strconv.ParseFloat(inputArgs[2], 64)
			if errPerc != nil {
				log.Fatal(errPerc)
			}

			//convert epsilon value to float64
			eps, errEps := strconv.ParseFloat(inputArgs[3], 64)
			if errEps != nil {
				log.Fatal(errEps)
			}

			pointCloud := ReadXYZ(inputFile)

			//format the file name to remove file extension
			//used later for output file name
			inputFile = strings.ReplaceAll(inputFile, ".xyz", "")

			//get number of iterations to satisfy confidence value
			numIterations := GetNumberOfIterations(confidence, percentage)

			for i := 0; i < 3; i++ {
				bestSupport := Plane3DwSupport{}
				stop := make(chan bool)

				//start concurrent pipeline
				pointGen := RandPointGenerator(pointCloud, stop)
				tripPointGen := TripletPointGenerator(pointGen, stop)
				takeN := TakeN(tripPointGen, stop, numIterations)
				estimator := PlaneEstimator(takeN)

				pointfinder := make([]chan Plane3DwSupport, num)

				//fan out support finder goroutines
				for j := 0; j < num; j++ {
					pointfinder[j] = SuppPointFinder(estimator, pointCloud, eps)
				}

				multiplexedOut := FanIn(pointfinder)

				//get dominant plane
				bestSupport = <-domPlaneFinder(multiplexedOut)

				//save plane points to file
				SaveXYZ(fmt.Sprintf(inputFile+"p%d.xyz", i+1), GetSupportingPoints(bestSupport.Plane3D, pointCloud, eps))

				pointCloud = RemovePlane(bestSupport.Plane3D, pointCloud, eps)
			}
			SaveXYZ((inputFile + "p0.xyz"), pointCloud)

			elapsed := float64(time.Since(start).Nanoseconds())
			elapsed = elapsed / 1000000

			avg += elapsed
		}
		avg = avg / 5
		output = append(output, timeOutput{num, avg})
	}
	SaveTime("Time3-6.csv", output)
}

func SaveTime(filename string, arr []timeOutput) {
	dat, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	defer dat.Close()
	_, errI := dat.WriteString("Number,Time\n")
	if errI != nil {
		log.Fatal(errI)
	}

	for _, out := range arr {
		_, err := dat.WriteString(fmt.Sprintf("%d,%f\n", out.num, out.time))
		if err != nil {
			log.Fatal(err)
		}
	}
}

type timeOutput struct {
	num  int
	time float64
}
*/
