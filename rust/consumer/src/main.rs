#![allow(non_snake_case)]

use futures::StreamExt;
use rabbitmq_stream_client::{Environment, types::OffsetSpecification};
use tracing::info;
use tracing_subscriber::FmtSubscriber;
use std::thread;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    //let opts = Opts::parse();
    let stream_name = String::from("mixing");
    let subscriber = FmtSubscriber::builder().finish();

    tracing::subscriber::set_global_default(subscriber).expect("setting default subscriber failed");

    let environment = Environment::builder()
        .host("localhost")
        .port(5552)
        .build()
        .await?;

    info!("im before start_consumer");

    start_consumer(environment, &stream_name).await?;

    Ok(())
}

async fn start_consumer(
    env: Environment,
  //  _opts: &Opts,
    stream: &String,
) -> Result<(), Box<dyn std::error::Error>> {

    println!("inside start_consumer");
   
    let mut consumer = env
    .consumer()
    .offset(OffsetSpecification::First)
    .build(stream)
    .await?;
   
    tokio::task::spawn(async move {
       
        loop  {
            
            println!("looping");
             let delivery = consumer.next().await;
             if delivery.is_some() {
                // TODO
                let delivery = delivery.unwrap();
                let delivery_tmp = delivery.unwrap();
                let message = delivery_tmp.message();
               
                // let val = delivery.unwrap();
                let data = String::from_utf8(message.data().unwrap().to_vec()).unwrap();
                println!("Got message {:?}", data);
                match &message.properties() {
                    Some(_T) => {
    
                        match &message.properties().unwrap().message_id {
                            Some(T) => println!("Message id: {:?}", T),
                            None => println!(""),
                        }
                        match &message.properties().unwrap().correlation_id {
                            Some(T) => println!("Correlation id: {:?}", T),
                            None => println!(""),
                        }
                        match &message.properties().unwrap().content_type {
                            Some(T) => println!("Content Type: {:?}", T),
                            None => println!(""),
                        }
                        match &message.properties().unwrap().content_encoding {
                            Some(T) => println!("Content Encoding: {:?}", T),
                            None => println!(""),
                        }
                        match &message.properties().unwrap().group_sequence {
                            Some(T) => println!("Group Sequence: {:?}", T),
                            None => println!(""),
                        }
                        match &message.properties().unwrap().user_id {
                            Some(T) => println!("User id: {:?}", T),
                            None => println!(""),
                        }
                  
                    }

                    None => println!("No properties for this message"),
                }

                match &message.application_properties() {
                    Some(T) => {
    
                        println!("application_properties {:?}", T)
                  
                    }

                    None => println!("No application properties for this message"),
                }
              

            
            }
            else {
                println!("breaking");
                break;
            }

            println!("ending looping");

        }

      
    });

    thread::sleep_ms(2000);
    println!("ending function");
    Ok(())

}