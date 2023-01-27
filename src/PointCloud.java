import java.io.*;
import java.util.ArrayList;
import java.util.Iterator;
import java.util.List;

public class PointCloud {
    private List<Point3D> points;
    public PointCloud() {
        this.points = new ArrayList<>();
    }

    public PointCloud(String filename) {
        this.points = read(filename);
    }

    public void addPoint(Point3D pt) {
        points.add(pt);
    }

    public Point3D getPoint () {
        return points.get((int)(Math.random() * ((points.size()-1) +1)));
    }

    private List<Point3D> getPoints() {
        return this.points;
    }

    public void save (String filename) {
        //create output file
        File outputFile = new File(filename);
        try {
            PrintWriter writer = new PrintWriter(outputFile);

            //write column headers to file
            writer.println("x\ty\tz");

            //write the coordinates, cluster label and rgb values of each point to file
            for (Point3D point : getPoints()) {
                writer.printf("%f\t%f\t%f\n", point.getX(), point.getY(), point.getZ());
            }
            writer.close();

        } catch (FileNotFoundException e) {
            e.printStackTrace();
        }
    }

    /**
     * Reads the given file and creates list of 3D points
     *
     * @param filename The name of the input file to be read
     * @return a list of Point3D objects
     *
     */
    private ArrayList<Point3D> read(String filename) {
        //string array to hold the contents of each line to be used in object creation
        String[] temp;

        //string used to temporarily hold the contents of each line
        String tempLine = " ";

        //will contain the processed 3D points
        ArrayList<Point3D> outputList = new ArrayList<Point3D>();

        try {
            File file = new File(filename);
            BufferedReader buffReader = new BufferedReader(new FileReader(file));

            // skip 'x,y,z' header
            tempLine = buffReader.readLine();

            //read and split each line according to delimiter
            while ((tempLine = buffReader.readLine()) != null) {
                temp = tempLine.split("\t");

                //create new Point3D object from parsed information and add to list
                outputList.add(new Point3D(Double.parseDouble(temp[0]),
                        Double.parseDouble(temp[1]), Double.parseDouble(temp[2])));
            }
            buffReader.close();
        } catch (IOException ioe) {
            ioe.printStackTrace();
        }

        //return list of 3D points
        return outputList;
    }

    Iterator<Point3D> iterator() {
        return getPoints().iterator();
    }

    public static void main(String[] args) {
        PointCloud pc = new PointCloud("PointCloud1.xyz");
        pc.save("test.xyz");
        Iterator itr = pc.iterator();
        while (itr.hasNext()) {
            System.out.println(itr.next());
        }
    }

}
