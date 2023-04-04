use std::time::{SystemTime, UNIX_EPOCH};
use rabbitmq_stream_client::{types::Message, Environment, TlsConfiguration, NoDedup, Producer};
use tracing::info;
use tracing_subscriber::FmtSubscriber;

const BATCH_SIZE: usize = 100;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    
    let stream_name = String::from("mixing");
    let subscriber = FmtSubscriber::builder().finish();

    tracing::subscriber::set_global_default(subscriber).expect("setting default subscriber failed");

    let tls_configuration = TlsConfiguration::builder().build();
    //.trust_certificate(false)
    //.trust_hostname(false)
    //.build();
  


    let environment = Environment::builder()
        .host("localhost")
        .port(5551)
        .username("guest")
        .password("guest")
        .tls(tls_configuration)
        .build()
        .await?;

    start_publisher(
        environment.clone(),
        &stream_name,
    ).await?;

    info!("final ending");

    Ok(())

}

async fn start_publisher(
    env: Environment,
  //  opts: &Opts,
    stream: &String,
) -> Result<(), Box<dyn std::error::Error>> {

    info!("im inside start_publisher");
    //let _ = env.stream_creator().create(&stream).await;

    info!("after im inside start_publisher");
    

    let producer = env
        .producer()
        .batch_size(BATCH_SIZE)
        .build(&stream)
        .await?;

    let is_batch_send = true;
    info!("before sending");
    tokio::task::spawn(async move {
        info!(
            "Starting producer with batch size {} and batch send {}",
            BATCH_SIZE, is_batch_send
        );
        info!("Sending {} simple messages", BATCH_SIZE);
        batch_send_simple(&producer).await;
        //info!("Sending {} simple messages with properties", BATCH_SIZE);
        //batch_send_with_properties(&producer).await;
        //info!("Sending {} simple messages with properties and application properties", BATCH_SIZE);
        //batch_send_with_application_properties(&producer).await;
        
    }).await?;
    info!("end im inside start_publisher");
    Ok(())
}

async fn batch_send_simple(producer: &Producer<NoDedup>) {
    let mut msg = Vec::with_capacity(BATCH_SIZE);
    for i in 0..BATCH_SIZE {
        msg.push(Message::builder().body(format!("rust message{}", i)).build());
    }

    producer
        .batch_send(msg, move |_| async move {})
        .await
        .unwrap();

}

/* 
async fn batch_send_with_properties(producer: &Producer<NoDedup>) {
    let mut msg = Vec::with_capacity(BATCH_SIZE);
    for i in 0..BATCH_SIZE {
        let message_id = String::from(format!("{}{}", "MyMessageId", i));
        let correlation_id = String::from(format!("{}{}", "MyCorrelationId", i));
        let message_builder = Message::builder()
        .body(format!("rust message{}", i))
        .properties()
        .message_id(message_id)
        .correlation_id(correlation_id)
        .content_type("text/plain")
        .content_encoding("utf-8")
        .group_sequence(9999)
        .reply_to_group_id("MyReplyToGroupId")
        .user_id("guest".as_bytes())
        .message_builder()
                
        .message_annotations()
      
        .message_builder()
        .build();
    
        msg.push(message_builder);

    }

    producer
        .batch_send(msg, move |_| async move {})
        .await
        .unwrap();


}

async fn batch_send_with_application_properties(producer: &Producer<NoDedup>) {
    let mut msg = Vec::with_capacity(BATCH_SIZE);
    for i in 0..BATCH_SIZE {
        let message_id = String::from(format!("{}{}", "MyMessageId", i));
        let correlation_id = String::from(format!("{}{}", "MyCorrelationId", i));
        let message_builder = Message::builder()
        .body(format!("rust message{}", i))
        .properties()
        .message_id(message_id)
        .correlation_id(correlation_id)
        .content_type("text/plain")
        .content_encoding("utf-8")
        .group_sequence(9999)
        .reply_to_group_id("MyReplyToGroupId")
        .user_id("guest".as_bytes())
        .message_builder()
        
        .application_properties()
        .insert("key_string","value")
        .insert("key2_int","1111")
        .insert("key2_decimal","10_000_000_000")
        
        .message_builder()
        
        .message_annotations()
      
        .message_builder()
        .build();
    
        msg.push(message_builder);

    }

    producer
        .batch_send(msg, move |_| async move {})
        .await
        .unwrap();


}*/