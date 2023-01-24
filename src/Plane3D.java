public class Plane3D {
    private double a;
    private double b;
    private double c;
    private double d;
    public Plane3D(Point3D p1, Point3D p2, Point3D p3) {
        Point3D vector1 = new Point3D(p2.getX() - p1.getX(),p2.getY() - p1.getY(), p2.getZ() - p1.getZ());
        Point3D vector2 = new Point3D(p3.getX() - p1.getX(),p3.getY() - p1.getY(), p3.getZ() - p1.getZ());
        this.a = ((vector1.getY() * vector2.getZ()) - (vector2.getY() * vector1.getZ()));
        this.b = (vector1.getZ() * vector2.getX()) - (vector1.getX() * vector2.getZ());
        this.c = (vector1.getX() * vector2.getY()) - (vector1.getY() * vector2.getX());
        this.d = (a * p1.getX()) + (b * p1.getY()) + (c * p1.getZ());
    }

    public Plane3D(double a, double b, double c, double d) {
        this.a = a;
        this.b = b;
        this.c = c;
        this.d = d;
    }

    public double getDistance(Point3D pt) {
        return (Math.abs((a*pt.getX()) + (b*pt.getY()) + (c*pt.getZ()) + d))/
                (Math.sqrt(Math.pow(a,2) + Math.pow(b,2) + Math.pow(c,2)));
    }

    @Override
    public String toString() {
        return String.format("a=%f, b=%f, c=%f, d=%f", a,b,c,d);
    }
}
