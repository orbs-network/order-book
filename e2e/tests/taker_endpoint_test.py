"""
A set of e2e tests for the Maker endpoints of the Orderbook API.
The tests use the Python SDK to interact with the API.
A local Orderbook instance is required to run the tests.
"""


import os

import pytest
import uuid
import requests
import math

from orbs_orderbook import CreateOrderInput, OrderBookSDK, OrderSigner
from orbs_orderbook.exceptions import ErrApiRequest


PRIVATE_KEY = os.environ.get(
    "PRIVATE_KEY", "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
)

BASE_URL = os.environ.get("BASE_URL", "http://localhost")

API_KEY = os.environ.get("API_KEY", "og4lpqQUILyciacspkFESHE1qrXIxpX1")

SYMBOL = "MATIC-USDC"
PRICE = "0.86500000"
SIZE = "40"


@pytest.fixture
def ob_client():
    yield OrderBookSDK(base_url=BASE_URL, api_key=API_KEY)


@pytest.fixture
def ob_signer(ob_client):
    yield OrderSigner(
        private_key=PRIVATE_KEY,
        sdk=ob_client,
    )


@pytest.fixture(autouse=True, scope="function")
def cancel_all_orders(ob_client, ob_signer):
    print("canceling all orders ---------")
    try:
        ob_client.cancel_all_orders()
    except ErrApiRequest as e:
        if e.status_code != 404:
            raise e
        print("No orders to cancel")
        pass


@pytest.fixture
def create_spread(ob_client, ob_signer):
    print("create spread ---------")
    # clear ob
    # cancel_all_orders(ob_client, ob_signer)

    price = 0.8
    size = 10
    orders = []

    for i in range(3):
        sell_price = price + (i + 1) * 0.01
        sell_price = math.floor(sell_price * 100) / 100

        print("sell_price: " + str(sell_price))
        order_input = CreateOrderInput(
            price=str(sell_price),
            size=str(size + i * 10),
            symbol=SYMBOL,
            side="sell",
            client_order_id=str(uuid.uuid4()),
        )

        signature, message = ob_signer.prepare_and_sign_order(order_input)

        order = ob_client.create_order(
            order_input=order_input,
            signature=signature,
            message=message,
        )
        orders.append(order)

        buy_price = price - (i + 1) * 0.01
        buy_price = math.floor(buy_price * 100) / 100
        print("buy_price: " + str(buy_price))
        order_input = CreateOrderInput(
            price=str(buy_price),
            size=str(size + i * 10),
            symbol=SYMBOL,
            side="buy",
            client_order_id=str(uuid.uuid4()),
        )

        signature, message = ob_signer.prepare_and_sign_order(order_input)
        order = ob_client.create_order(
            order_input=order_input,
            signature=signature,
            message=message,
        )
        orders.append(order)

    # yield orders
    yield orders


# test simple quote on depth spread 1200-700
def test_quote(cancel_all_orders, create_spread):
    print("tes quote")
    # json_data = {"InAmount": "1000000", "InToken": "USDC", "OutToken": "MATIC"}
    # # Convert the JSON data to a string
    # json_string = json.dumps(json_data)
    # # Set the headers for the request
    # headers = {"Content-Type": "application/json", "Accept": "application/json"}

    # url = f"{BASE_URL}/taker/v1/quote"
    # print(url)

    # response = requests.get(url, data=json_string, headers=headers)
    # print(response)
