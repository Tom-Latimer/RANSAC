import java.util.Iterator;

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
        Plane3D dominantPlane;
        Iterator<Point3D> itr = pc.iterator();

        for (int i = 0; i < numberOfIterations; i++) {

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
    }
}
