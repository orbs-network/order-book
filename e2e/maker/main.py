import json
import requests
import os
import time
from uuid import uuid4
from decimal import Decimal

mock_api_key = "abcdef12345"
depth_size = 5
HOST = "localhost"


class Ticker:
    def __init__(self, price, symbol):
        self.price = price
        self.symbol = symbol


class AddOrderReq:
    def __init__(self, price, size, side, symbol, client_order_id):
        self.price = price
        self.size = size
        self.side = side
        self.symbol = symbol
        self.client_order_id = client_order_id


def on_tick(url):
    response = requests.get(url)
    if response.status_code != 200:
        print(f"HTTP request failed with status code {response.status_code}")
        return None

    ticker_data = response.json()
    ticker = Ticker(ticker_data["price"], ticker_data["symbol"])
    print(f"ETH-USD Price: {ticker.price}")
    return ticker


def cancel_all_orders():
    url = f"{HOST}/api/v1/orders"
    headers = {"X-API-Key": f"Bearer {mock_api_key}"}

    response = requests.delete(url, headers=headers)
    response.raise_for_status()
    print("Canceled all orders")


def place_order(side, price, size):
    cOId = str(uuid4())
    body = AddOrderReq(
        price=str(price),
        size=str(size),
        side=side,
        symbol="ETH-USD",
        client_order_id=cOId,
    )
    url = f"{HOST}/api/v1/order"
    headers = {
        "X-API-Key": f"Bearer {mock_api_key}",
        "Content-Type": "application/json",
    }

    response = requests.post(url, data=json.dumps(body.__dict__), headers=headers)
    response.raise_for_status()
    print("Created order with clientOrderId:", cOId)
    print("Status code:", response.status_code)


def update_orders(price):
    cancel_all_orders()
    factor = Decimal("1.001")
    cur_price = price
    print("------ Market Price: ", price)

    # ASK
    for i in range(depth_size):
        cur_price *= factor
        print("Ask Price: ", cur_price)
        cur_size = Decimal((i + 1) * 10)
        place_order("sell", cur_price, cur_size)

    # BIDS
    factor = Decimal("0.999")
    cur_price = price
    for i in range(depth_size):
        cur_price *= factor
        print("Bid Price: ", cur_price)
        cur_size = Decimal((i + 1) * 10)
        place_order("buy", cur_price, cur_size)


def main():
    url = "https://www.binance.com/api/v3/ticker/price?symbol=ETHUSDT"
    print("Ticker URL: ", url)
    host = os.getenv("ORDERBOOK_HOST")
    if host:
        global HOST
        HOST = host

    while True:
        ticker = on_tick(url)
        if ticker:
            price = Decimal(ticker.price)
            update_orders(price)

            print("Sleeping for 10 seconds...")
            time.sleep(10)


if __name__ == "__main__":
    main()
