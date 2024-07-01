import asyncio
import websockets
import os


async def connect_and_listen():
    # url = "ws://localhost:8080/ws/orders"
    prod = "orderbook.orbs.network"
    stg = "stg.orderbook.orbs.network"
    dev = "ob-server-development-c61f13bc12ed.herokuapp.com"

    url = f"wss://{dev}/api/v1/ws/orders"
    print("url: ", url)
    api_key = os.getenv("API_KEY")

    if api_key is None:
        print("Error: API_KEY environment variable not set.")
        return

    headers = {"X-API-KEY": f"bearer {api_key}"}

    async with websockets.connect(url, extra_headers=headers) as websocket:
        while True:
            message = await websocket.recv()
            print(message)


asyncio.get_event_loop().run_until_complete(connect_and_listen())
