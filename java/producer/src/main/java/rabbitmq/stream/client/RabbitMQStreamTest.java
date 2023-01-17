package rabbitmq.stream.client;


import com.rabbitmq.stream.Environment;
import com.rabbitmq.stream.Producer;
import com.rabbitmq.stream.OffsetSpecification;
import java.util.UUID;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicLong;
import java.util.stream.IntStream;

class RabbitMQStreamTest {

    private static int  MESSAGES = 100;

    public static void main(String[] args) throws Exception {
        Environment environment = Environment.builder().build();  // <1>
        String stream = "mixing";
       
        // end::sample-environment[]
        // tag::sample-publisher[]
        System.out.println("Starting publishing...");

        SendSimpleMessages(environment, stream);
        SendWithProperties(environment, stream);
        SendWithApplicationProperties(environment, stream);

        System.exit(0);
       
}

public static void SendSimpleMessages(Environment environment, String stream) throws Exception {

    CountDownLatch publishConfirmLatch = new CountDownLatch(MESSAGES);
    Producer producer = environment.producerBuilder()  // <1>
             .stream(stream)
             .build();
        IntStream.range(0, MESSAGES)
            .forEach(i -> producer.send(  // <2>
                producer.messageBuilder()                    // <3>
                .addData(("message"+i).getBytes())   // <3>
                .build(),                                // <3>
            confirmationStatus -> publishConfirmLatch.countDown()  // <4>
        ));
        publishConfirmLatch.await(10, TimeUnit.SECONDS);  // <5>
        producer.close();  // <6>
        System.out.printf("Published %,d messages%n", MESSAGES);

}

public static void SendWithProperties(Environment environment, String stream) throws Exception {

    CountDownLatch publishConfirmLatch = new CountDownLatch(MESSAGES);
    Producer producer = environment.producerBuilder()  // <1>
            
             .stream(stream)
             .build();
        IntStream.range(0, MESSAGES)
            .forEach(i -> producer.send(  // <2>
                producer.messageBuilder() 
                .properties()
                    .messageId("MyMessageId"+i)
                    .correlationId("MyCorrelationId"+i)
                    .contentType("text/plain")   
                    .contentEncoding("utf-8")
                    .groupSequence(9999)
                    .replyToGroupId("MyReplyToGroupId")
                    .userId("guest".getBytes())                // <3>
                .messageBuilder()  
                .addData(("message"+i).getBytes())   // <3>
                .build(),                                // <3>
            confirmationStatus -> publishConfirmLatch.countDown()  // <4>
        ));
        publishConfirmLatch.await(10, TimeUnit.SECONDS);  // <5>
        producer.close();  // <6>
        System.out.printf("Published %,d messages%n", MESSAGES);

}

public static void SendWithApplicationProperties(Environment environment, String stream) throws Exception {

    CountDownLatch publishConfirmLatch = new CountDownLatch(MESSAGES);
    Producer producer = environment.producerBuilder()  // <1>
            
             .stream(stream)
             .build();
        IntStream.range(0, MESSAGES)
            .forEach(i -> producer.send(  // <2>
                producer.messageBuilder() 
                .properties()
                    .messageId("MyMessageId"+i)
                    .correlationId("MyCorrelationId"+i)
                    .contentType("text/plain")   
                    .contentEncoding("utf-8")
                    .groupSequence(9999)
                    .replyToGroupId("MyReplyToGroupId")
                    .userId("guest".getBytes())  
                .messageBuilder()       // <3>
                .applicationProperties()
                .entry("key_string","value")
                .entry("key2_int","1111")
                .entry("key2_decimal","10_000_000_000")
                .messageBuilder()  
                .addData(("message"+i).getBytes())   // <3>
                .build(),                                // <3>
            confirmationStatus -> publishConfirmLatch.countDown()  // <4>
        ));
        publishConfirmLatch.await(10, TimeUnit.SECONDS);  // <5>
        producer.close();  // <6>
        System.out.printf("Published %,d messages%n", MESSAGES);

}
}