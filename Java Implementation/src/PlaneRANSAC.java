import java.io.FileNotFoundException;
import java.util.ArrayList;
import java.util.Iterator;
import java.util.List;
import java.util.Scanner;

public class PlaneRANSAC {

    private PointCloud pc;
    private double eps;
    public PlaneRANSAC(PointCloud pc) {
        this.pc = pc;
        this.eps = 1;
    }

    public void setEps(double eps) {
        this.eps = eps;
    }

    public double getEps() {
        return this.eps;
    }

    public int getNumberOfIterations(double confidence, double percentageOfPointsOnPlane) {
        if (confidence > 1 || percentageOfPointsOnPlane > 1) {
            throw new IllegalArgumentException("Input must be less than 1");
        }
        return (int) Math.ceil(Math.log(1-confidence)/Math.log(1-Math.pow(percentageOfPointsOnPlane,3)));
    }

    public void run(int numberOfIterations, String filename) {

        int bestSupport = 0;
        Plane3D dominantPlane = null;

        //might need to move this into the loop
        Iterator<Point3D> itr;

        for (int i = 0; i < numberOfIterations; i++) {

            itr  = pc.iterator();

            int currentBestSupport = 0;

            Plane3D currentPlane = new Plane3D(pc.getPoint(),pc.getPoint(),pc.getPoint());

            while (itr.hasNext()) {
                Point3D tempPoint = itr.next();
                if (currentPlane.getDistance(tempPoint) < getEps()) {
                    currentBestSupport++;
                }
            }

            if (currentBestSupport > bestSupport) {
                bestSupport = currentBestSupport;
                dominantPlane = currentPlane;
            }
        }

        PointCloud dominantPlanePC = new PointCloud();
         itr = pc.iterator();


        while (itr.hasNext()) {
            Point3D tempPoint = itr.next();
            if (dominantPlane.getDistance(tempPoint) < getEps()) {
                dominantPlanePC.addPoint(tempPoint);
                itr.remove();
            }
        }

        dominantPlanePC.save(filename);
    }

    public static void main(String[] args) {
        Scanner input = new Scanner(System.in);
        try {
            System.out.println("Please enter a Point Cloud .xyz file to be read from (exclude extension):");
            String inputFilename = input.nextLine();
            if (inputFilename.isEmpty() || inputFilename.endsWith(".xyz")) {
                throw new IllegalArgumentException("Invalid input filename.");
            }

            PointCloud pc = new PointCloud(inputFilename);
            PlaneRANSAC pRANSAC = new PlaneRANSAC(pc);

            double confidence, percentage = 0;
            System.out.println("Please enter a confidence value (between 0 and 1):");
            confidence = input.nextDouble();

            System.out.println("Please enter a the percentage of points that support the dominant plane (between 0 and 1):");
            percentage = input.nextDouble();

            int iterations = pRANSAC.getNumberOfIterations(confidence,percentage);

            for (int i = 1; i <= 3; i++) {
                pRANSAC.run(iterations, inputFilename + "_p" + i);
            }

            pc.save(inputFilename + "_p0");

        } catch (IllegalArgumentException ie) {
            System.out.println("\n" + ie.getMessage());
        } finally {
            input.close();
        }
    }
}
