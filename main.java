import com.google.cloud.functions.HttpFunction;
import com.google.cloud.functions.HttpRequest;
import com.google.cloud.functions.HttpResponse;
import com.google.cloud.storage.Blob;
import com.google.cloud.storage.Bucket;
import com.google.cloud.storage.Storage;
import com.google.cloud.storage.StorageOptions;

import java.awt.image.BufferedImage;
import java.io.ByteArrayInputStream;
import java.io.ByteArrayOutputStream;
import java.io.IOException;

import javax.imageio.ImageIO;

public class ImageProcessorHttpFunction implements HttpFunction {

    public HttpResponse process(HttpRequest request) {
        try {
          //Receiving the Raw data
            String imageDataString = request.getPayload();
          //conversion to bytes
            byte[] imageDataBytes = Base64.getDecoder().decode(imageDataString);
          //bytecode to image
            BufferedImage image = ImageIO.read(new ByteArrayInputStream(imageDataBytes));
          //operations on image 
            int newWidth = 200;
            int newHeight = 200;
            BufferedImage resizedImage = new BufferedImage(newWidth, newHeight, BufferedImage.TYPE_INT_RGB);
            resizedImage.getGraphics().drawImage(image.getScaledInstance(newWidth, newHeight, BufferedImage.SCALE_SMOOTH), 0, 0, null);
          //exporting to bytecode
            ByteArrayOutputStream outputStream = new ByteArrayOutputStream();
            ImageIO.write(resizedImage, "jpg", outputStream);
            byte[] processedImageData = outputStream.toByteArray();
            String processedImageDataString = Base64.getEncoder().encodeToString(processedImageData);
          //Successfull HTTP Response
            return HttpResponse.newBuilder()
                    .setStatusCode(200)
                    .setHeader("Content-Type", "text/plain")
                    .setBody(processedImageDataString)
                    .build();
        } catch (IOException e) {
          //Failed HTTP Response
            return HttpResponse.newBuilder()
                    .setStatusCode(500)
                    .setHeader("Content-Type", "text/plain")
                    .setBody("Error processing image: " + e.getMessage())
                    .build();
        }
    }
}
