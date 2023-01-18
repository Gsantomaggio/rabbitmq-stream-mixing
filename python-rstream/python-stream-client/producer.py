import asyncio
import uamqp
from rstream import Producer, AMQPMessage
from uamqp import c_uamqp, constants, errors, utils

MESSAGES = 100

async def publish():
    async with Producer('localhost', username='guest', password='guest') as producer:
        # await producer.create_stream('mixing')

        # sending basic messages
        print ("Sending " + str(MESSAGES) + " simple messages")
        for i in range(MESSAGES):
            amqp_message = AMQPMessage(
                body='hello: {}'.format(i),
            )
            #print (amqp_message)
            await producer.publish('mixing', amqp_message)

        print ("Sending " + str(MESSAGES) + " simple messages with properties")   
        # sending messages with properties 
        for i in range(MESSAGES):
            message_properties = uamqp.message.MessageProperties("MessageId"+str(i), None, bytes("guest",'utf-8'), None, "CorrelationId"+str(i), "text/plain", "utf-8", None, None, None, None, 9999, "MyReplyToGroupId", None)
           
            amqp_message = AMQPMessage(
                body='hello: {}'.format(i),
                properties = message_properties,
            )
            #print (amqp_message)
            await producer.publish('mixing', amqp_message)   
    
        print ("Sending " +  str(MESSAGES) + " simple messages with properties and application properties")   
        # sending messages with properties and application properties
        for i in range(MESSAGES):
            message_properties = uamqp.message.MessageProperties("MessageId"+str(i), None, bytes("guest",'utf-8'), None, "CorrelationId"+str(i), "text/plain", "utf-8", None, None, None, None, 9999, "MyReplyToGroupId", None)
            application_properties = {
                "key_string": "values",
                "key2_int": "1111",
                "key2_decimal": 10_000_000_000
            }
            amqp_message = AMQPMessage(
                body='hello: {}'.format(i),
                properties = message_properties,
                application_properties = application_properties,
            
            )
            #print (amqp_message)
            await producer.publish('mixing', amqp_message)

asyncio.run(publish())
