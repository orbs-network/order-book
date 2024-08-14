import asyncio
import websockets
import os


async def connect_and_listen(client_id):
    url = "ws://127.0.0.1/api/v1/ws/orders"
    api_key = os.getenv("API_KEY")

    if api_key is None:
        print("Error: API_KEY environment variable not set.")
        return

    headers = {"X-API-KEY": f"Bearer {api_key}"}

    while True:
        try:
            async with websockets.connect(url, extra_headers=headers, ping_interval=None) as websocket:
                print(f"Client {client_id} connected to server")
                while True:
                    try:
                        message = await websocket.recv()
                        # print(f"Client {client_id} received message: {message}")
                        print(f"Client {client_id} received message")
                    except websockets.exceptions.ConnectionClosedError as e:
                        print(f"Client {client_id} connection closed with error: {e}")
                        break
        except (websockets.exceptions.ConnectionClosedError, websockets.exceptions.ConnectionClosed) as e:
            print(f"Client {client_id} failed to connect or connection lost: {e}")
        except Exception as e:
            print(f"Client {client_id} unexpected error: {e}")

        print(f"Client {client_id} reconnecting in 5 seconds...")
        await asyncio.sleep(5)  # Wait before attempting to reconnect


async def main(num_clients):
    # Create multiple WebSocket connections
    tasks = [asyncio.create_task(connect_and_listen(client_id)) for client_id in range(num_clients)]
    await asyncio.gather(*tasks)


if __name__ == "__main__":
    num_clients = 5  # Set the number of WebSocket connections to create
    asyncio.get_event_loop().run_until_complete(main(num_clients))