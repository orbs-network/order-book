"""
A set of e2e tests for the Maker endpoints of the Orderbook API.
The tests use the Python SDK to interact with the API.
A local Orderbook instance is required to run the tests.
"""


import os

import pytest
import uuid
import json
import requests
import math
from decimal import *

from orbs_orderbook import CreateOrderInput, OrderBookSDK, OrderSigner
from orbs_orderbook.exceptions import ErrApiRequest

TOKEN_DEC = {"MATIC": 18, "USDC": 6}

PRIVATE_KEY = os.environ.get(
    "PRIVATE_KEY", "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
)

BASE_URL = os.environ.get("BASE_URL", "http://localhost")

API_KEY = os.environ.get("API_KEY", "og4lpqQUILyciacspkFESHE1qrXIxpX1")

SYMBOL = "MATIC-USDC"
SPREAD_PRICE = 0.8
PRICER_OFFSET = 0.02
SPREAD_SIZE = 10
ORDER_SIDE_COUNT = 3


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

    price = SPREAD_PRICE
    size = SPREAD_SIZE
    orders = []

    for i in range(ORDER_SIDE_COUNT):
        sell_price = price + (i + 1) * PRICER_OFFSET
        sell_price = math.floor(sell_price * 100) / 100

        print("sell_price: " + str(sell_price))
        order_input = CreateOrderInput(
            price=str(sell_price),
            size=str(size + i * SPREAD_SIZE),
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

        buy_price = price - (i + 1) * PRICER_OFFSET
        buy_price = math.floor(buy_price * 100) / 100
        print("buy_price: " + str(buy_price))
        order_input = CreateOrderInput(
            price=str(buy_price),
            size=str(size + i * SPREAD_SIZE),
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


def toTokenDec(num, dec):
    return str(int(num * 10**dec))


def call_quote(inAmount, inToken, outToken):
    json_data = {"InAmount": str(inAmount), "InToken": inToken, "OutToken": outToken}
    # Convert the JSON data to a string
    json_string = json.dumps(json_data)
    # Set the headers for the request
    headers = {"Content-Type": "application/json", "Accept": "application/json"}

    url = f"{BASE_URL}/taker/v1/quote"
    res = requests.get(url, data=json_string, headers=headers)
    assert res.status_code == 200, "res is not 200"
    return res


def call_quote_size(inToken, inSize, outToken, outSize):
    # test simple quote on depth spread
    res = call_quote(toTokenDec(inSize, TOKEN_DEC[inToken]), inToken, outToken)

    obj = res.json()
    assert obj is not None, "json is none"

    expectedOutAmount = toTokenDec(outSize, TOKEN_DEC[outToken])
    outTokenDec = TOKEN_DEC[outToken] - 2
    # breakpoint()
    assert (
        obj["outAmount"][:outTokenDec] == expectedOutAmount[:outTokenDec]
    ), "outAmount is wrong"
    assert obj["outToken"] == outToken, "outToken is wrong"
    assert obj["inToken"] == inToken, "inToken is wrong"
    assert obj["swapId"] == "", "swapId should be empty"


def test_quote(cancel_all_orders, create_spread):
    # sell to bids single
    call_quote_size("MATIC", 1, "USDC", SPREAD_PRICE - PRICER_OFFSET)

    # entire sell side
    tot_matic = 0
    tot_usdc = 0
    for i in range(ORDER_SIDE_COUNT):
        matic_size = SPREAD_SIZE + i * SPREAD_SIZE
        price = SPREAD_PRICE - (i + 1) * PRICER_OFFSET
        tot_usdc += matic_size * price
        tot_matic += matic_size

    call_quote_size("MATIC", tot_matic, "USDC", tot_usdc)

    # buy from ask
    call_quote_size(
        "USDC", 1, "MATIC", Decimal(1) / Decimal(SPREAD_PRICE + PRICER_OFFSET)
    )

    tot_matic = 0
    tot_usdc = 0
    for i in range(ORDER_SIDE_COUNT):
        matic_size = SPREAD_SIZE + i * SPREAD_SIZE
        price = SPREAD_PRICE + (i + 1) * PRICER_OFFSET
        tot_usdc += matic_size * price
        tot_matic += matic_size

    call_quote_size("USDC", tot_usdc, "MATIC", tot_matic)
