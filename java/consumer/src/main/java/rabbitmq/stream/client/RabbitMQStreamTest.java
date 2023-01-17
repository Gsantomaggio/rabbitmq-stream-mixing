package rabbitmq.stream.client;


import com.rabbitmq.stream.Environment;
import com.rabbitmq.stream.Consumer;
import com.rabbitmq.stream.OffsetSpecification;
import java.util.UUID;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicLong;
import java.util.stream.IntStream;

class RabbitMQStreamTest {
    public static void main(String[] args) throws Exception {
        System.out.println("Hello World!"); // Display the string.

        System.out.println("Connecting...");


        AtomicLong sum = new AtomicLong(0);
        int messageCount = 1000000;
        CountDownLatch consumeLatch = new CountDownLatch(messageCount);

        Environment environment = Environment.builder().build();  
        String stream = "mixing";
        //environment.streamCreator().stream(stream).create();  
       
        Consumer consumer = environment.consumerBuilder()  // <1>
            .stream(stream)
            .offset(OffsetSpecification.first()) // <2>
            .messageHandler((offset, message) -> {  // <3>
            System.out.println(new String(message.getBodyAsBinary()));
            if (message.getProperties() != null) {
    
                if (message.getProperties().getMessageIdAsString() != null) {
                    System.out.println("MessageID: " + new String(message.getProperties().getMessageIdAsString()));
                }
                if (message.getProperties().getCorrelationIdAsString() != null) {
                    System.out.println("CorrelationID: " + new String(message.getProperties().getCorrelationIdAsString()));
                }
                if (message.getProperties().getContentType() != null) {
                    System.out.println("ContentType: " + new String(message.getProperties().getContentType()));
                }
                if (message.getProperties().getContentEncoding() != null) {

                    System.out.println("Encoding: " + new String(message.getProperties().getContentEncoding()));
                }
                if (message.getProperties().getUserId() != null) {
                    System.out.println("UserID: " + new String(message.getProperties().getUserId()));
                }
                System.out.println("GroupSequence: " + message.getProperties().getGroupSequence());
            }

            if (message.getApplicationProperties() != null) {

                System.out.println("Application Property key_string: " + message.getApplicationProperties().get("key_string").toString());
                System.out.println("Application Property key2_int: " + message.getApplicationProperties().get("key2_int").toString());
                System.out.println("Application Property key2_decimal: " + message.getApplicationProperties().get("key2_decimal").toString());


            }
      
            consumeLatch.countDown();  // <5>
            })
        .build();

        consumeLatch.await(2, TimeUnit.SECONDS);  // <6>

        //System.out.println("Sum: " + sum.get());  // <7>

        consumer.close();  // <8>
        // end::sample-consumer[]

        // tag::sample-environment-close[]
        //environment.deleteStream(stream);  // <1>
        environment.close();  // <2>
    }
}
