public class Point3D {
    //Instance Variables

    /**
     * The x coordinate of the point
     */
    private double x;

    /**
     * The y coordinate of the point
     */
    private double y;

    /**
     * The z coordinate of the point
     */
    private double z;

    /**
     * Constructs an instance of Point3D
     *
     * @param x the x coordinate of the point
     * @param y the y coordinate of the point
     * @param z the z coordinate of the point
     *
     */
    public Point3D (double x, double y, double z) {
        this.x = x;
        this.y = y;
        this.z = z;
    }

    //Getters
    /**
     * Getter for x coordinate
     *
     * @return the x coordinate of the point
     *
     */
    public double getX () {
        return x;
    }

    /**
     * Getter for the y coordinate
     *
     * @return the y coordinate of the point
     *
     */
    public double getY () {
        return y;
    }

    /**
     * Getter for the z coordinate
     *
     * @return the z coordinate of the point
     *
     */
    public double getZ() {
        return z;
    }
}
