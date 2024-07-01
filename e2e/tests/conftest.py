import os

import pytest
from orbs_orderbook import CreateOrderInput, OrderBookSDK, OrderSigner
from orbs_orderbook.exceptions import ErrApiRequest

BASE_URL = os.environ.get("BASE_URL", "http://localhost")
PRIVATE_KEY = os.environ.get(
    "PRIVATE_KEY", "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
)
API_KEY = os.environ["API_KEY"]

CLIENT_OID = "550e8400-e29b-41d4-a716-446655440000"
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
    try:
        ob_client.cancel_all_orders_by_symbol(SYMBOL)
    except Exception:
        print("No orders to cancel")
        pass


@pytest.fixture
def create_new_orders(ob_client, ob_signer, cancel_all_orders):
    order_input = CreateOrderInput(
        price=PRICE,
        size=SIZE,
        symbol=SYMBOL,
        side="sell",
        client_order_id=CLIENT_OID,
    )

    signature, message = ob_signer.prepare_and_sign_order(order_input)

    yield [
        ob_client.create_order(
            order_input=order_input,
            signature=signature,
            message=message,
        )
    ]
