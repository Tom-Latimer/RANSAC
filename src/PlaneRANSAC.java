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

}
