import os
import time
from dataclasses import dataclass
from decimal import Decimal
from uuid import uuid4

import requests
from orbs_orderbook import CreateOrderInput, OrderBookSDK, OrderSigner

IS_DISABLED = os.environ.get("IS_DISABLED", "false").lower() == "true"
BASE_URL = os.environ.get("BASE_URL", "http://localhost")
API_KEY = os.environ.get("API_KEY", "abc123")
PRIVATE_KEY = os.environ.get(
    "PRIVATE_KEY", "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
)
TICKER_URL = os.environ.get(
    "TICKER_URL", "https://www.binance.com/api/v3/ticker/price?symbol=MATICUSDC"
)
TICKER_SYMBOL = os.environ.get("TICKER_SYMBOL", "MATIC-USDC")
DEPTH_SIZE = os.environ.get("DEPTH_SIZE", "5")
SLEEP_TIME = os.environ.get("SLEEP_TIME", "10")


class Ticker:
    def __init__(self, price, symbol):
        self.price = price
        self.symbol = symbol


@dataclass
class Client:
    ob_sdk: OrderBookSDK
    signer: OrderSigner

    def update_orders(self, price):
        try:
            self.ob_sdk.cancel_all_orders()
        except requests.exceptions.HTTPError as e:
            if e.response.status_code == 404:
                print("No orders to cancel")
            else:
                print("Error cancelling orders:", e)
                return
        print("Cancelled all orders")

        factor = Decimal("1.001")
        cur_price = price
        print("------ Market Price: ", price)

        # ASK
        for i in range(int(DEPTH_SIZE)):
            cur_price *= factor
            print("Ask Price: ", cur_price)
            cur_size = Decimal((i + 1) * 10)

            order_input = CreateOrderInput(
                price=str(self._round(cur_price)),
                size=str(cur_size),
                side="sell",
                symbol=TICKER_SYMBOL,
                clientOrderId=str(uuid4()),
            )

            signature, message_data = self.signer.prepare_and_sign_order(order_input)

            try:
                self.ob_sdk.create_order(
                    order_input=order_input,
                    signature=signature,
                    message_data=message_data,
                )
            except Exception as e:
                print("Error creating sell order:", e)
                continue

            print(f"Created sell order of size {cur_size} at price {cur_price}")

        print("\n")
        print("----------------------------------")
        print("\n")

        # BIDS
        factor = Decimal("0.999")
        cur_price = price
        for i in range(int(DEPTH_SIZE)):
            cur_price *= factor
            print("Bid Price: ", cur_price)
            cur_size = Decimal((i + 1) * 10)

            order_input = CreateOrderInput(
                price=str(self._round(cur_price)),
                size=str(cur_size),
                side="buy",
                symbol=TICKER_SYMBOL,
                clientOrderId=str(uuid4()),
            )

            signature, message_data = self.signer.prepare_and_sign_order(order_input)

            try:
                self.ob_sdk.create_order(
                    order_input=order_input,
                    signature=signature,
                    message_data=message_data,
                )
            except Exception as e:
                print("Error creating buy order:", e)
                continue

            print(f"Created buy order of size {cur_size} at price {cur_price}")

        print("\n")
        print("----------------------------------")
        print("\n")

    def on_tick(self, url):
        response = requests.get(url)
        if response.status_code != 200:
            print(f"HTTP request failed with status code {response.status_code}")
            return None

        ticker_data = response.json()
        ticker = Ticker(ticker_data["price"], ticker_data["symbol"])
        print(f"ETH-USD Price: {ticker.price}")
        return ticker

    def _round(self, value):
        decimal_value = Decimal(value)
        # 8 decimal places
        EIGHT_PLACES = Decimal("0.00000001")
        rounded_value = decimal_value.quantize(EIGHT_PLACES)
        return rounded_value


def main():
    sdk = OrderBookSDK(base_url=BASE_URL, api_key=API_KEY)
    ob_client = Client(
        ob_sdk=sdk,
        signer=OrderSigner(private_key=PRIVATE_KEY, sdk=sdk),
    )

    print("----------------------------------")
    print(f"BASE_URL: {BASE_URL}")
    print(f"TICKER_URL: {TICKER_URL}")
    print(f"TICKER_SYMBOL: {TICKER_SYMBOL}")
    print("----------------------------------")

    try:
        while True:
            ticker = ob_client.on_tick(TICKER_URL)
            if not ticker:
                print(f"No price data. Sleeping for {SLEEP_TIME} seconds...")
                time.sleep(int(SLEEP_TIME))
                continue

            price = Decimal(ticker.price)
            ob_client.update_orders(price)

            print(f"Sleeping for {SLEEP_TIME} seconds...")
            time.sleep(int(SLEEP_TIME))
    except KeyboardInterrupt:
        print("\nExiting...")


if __name__ == "__main__":
    if IS_DISABLED:
        print("Maker is disabled")
        exit(0)

    main()
