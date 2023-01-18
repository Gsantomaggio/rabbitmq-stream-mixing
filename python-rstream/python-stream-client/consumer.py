import asyncio
import signal
from rstream import Consumer, amqp_decoder, AMQPMessage

async def consume():
    consumer = Consumer(
        host='localhost',
        port=5552,
        vhost='/',
        username='guest',
        password='guest',
    )

    loop = asyncio.get_event_loop()
    loop.add_signal_handler(signal.SIGINT, lambda: asyncio.create_task(consumer.close()))

    def on_message(msg: AMQPMessage):
        print('Got message: {}'.format(msg))
        print('Got properties: {}'.format(msg.properties))
        print('Got application_properties: {}'.format(msg.application_properties))

    await consumer.start()
    await consumer.subscribe('mixing', on_message, decoder=amqp_decoder)
    await consumer.run()

# main coroutine
async def main():
    # schedule the task
    task = asyncio.create_task(consume())
    # suspend a moment
      # wait a moment
    await asyncio.sleep(5)
    # cancel the task
    was_cancelled = task.cancel()
   
    # report a message
    print('Main done')
 
# run the asyncio program
asyncio.run(main())

