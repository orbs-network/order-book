"""
A set of e2e tests for the Maker endpoints of the Orderbook API.
The tests use the Python SDK to interact with the API.
A local Orderbook instance is required to run the tests.
"""

import pytest
from orbs_orderbook import CreateOrderInput
from orbs_orderbook.exceptions import ErrApiRequest
from conftest import CLIENT_OID, SIZE, SYMBOL, API_KEY

# TODO 1: Add test for get market depth
# TODO 2: Add different test scenaries (eg. different failed states)


def test_create_order_success(ob_client, ob_signer):
    order_input = CreateOrderInput(
        price="0.86500000",
        size="40",
        symbol="MATIC-USDC",
        side="sell",
        client_order_id=CLIENT_OID,
    )

    signature, message = ob_signer.prepare_and_sign_order(order_input)

    res = ob_client.create_order(
        order_input=order_input,
        signature=signature,
        message=message,
    )

    assert res.order_id, "Order was not created"


def test_create_order_fails_with_same_clientoid(
    ob_client, ob_signer, create_new_orders
):
    order_input = CreateOrderInput(
        price="0.86500000",
        size="40",
        symbol="MATIC-USDC",
        side="sell",
        client_order_id=CLIENT_OID,
    )

    signature, message = ob_signer.prepare_and_sign_order(order_input)

    with pytest.raises(ErrApiRequest) as e:
        ob_client.create_order(
            order_input=order_input,
            signature=signature,
            message=message,
        )

    assert (
        e.value.status_code == 409
    ), "Order was created when it should have been rejected due to same clientOrderId"


def test_cancel_order_by_oid(ob_client, ob_signer, create_new_orders):
    new_oid = create_new_orders[0].order_id
    res = ob_client.cancel_order_by_id(new_oid)

    assert res.order_id == new_oid, "Order was not cancelled by orderId"


def test_cancel_order_by_oid_fails_when_cancelling_same_order(
    ob_client, ob_signer, create_new_orders
):
    new_oid = create_new_orders[0].order_id

    success = ob_client.cancel_order_by_id(new_oid)

    assert success.order_id == new_oid, "Order was not cancelled by orderId"

    with pytest.raises(ErrApiRequest) as e:
        ob_client.cancel_order_by_id(new_oid)

    assert e.value.status_code == 404, "Order was cancelled when it should have failed"


def test_cancel_order_by_clientoid(ob_client, ob_signer, create_new_orders):
    new_oid = create_new_orders[0].order_id

    res = ob_client.cancel_order_by_client_id(CLIENT_OID)

    assert res.order_id == new_oid, "Order was not cancelled by clientOrderId"


def test_cancel_order_by_clientoid_fails_when_cancelling_same_order(
    ob_client, ob_signer, create_new_orders
):
    new_oid = create_new_orders[0].order_id

    success = ob_client.cancel_order_by_client_id(CLIENT_OID)

    assert success.order_id == new_oid, "Order was not cancelled by clientOrderId"

    with pytest.raises(ErrApiRequest) as e:
        ob_client.cancel_order_by_client_id(CLIENT_OID)

    assert e.value.status_code == 404, "Order was cancelled when it should have failed"


def test_get_symbols(ob_client):
    res = ob_client.get_symbols()

    assert len(res) > 0, "No symbols returned"
    assert "symbol" in res[0], "Symbol not returned"
    assert "name" in res[0], "Name not returned"


def test_get_order_by_id(ob_client, ob_signer, create_new_orders):
    new_oid = create_new_orders[0].order_id

    res = ob_client.get_order_by_id(new_oid)

    assert res.order_id == new_oid, "Order was not returned by orderId"
    assert res.client_order_id == CLIENT_OID, "Order was not returned by clientOrderId"
    assert res.price == "0.865", "Order was not returned by price"
    assert res.size == SIZE, "Order was not returned by size"
    assert res.side == "sell", "Order was not returned by side"
    assert res.symbol == SYMBOL, "Order was not returned by symbol"


def test_get_order_by_clientoid(ob_client, ob_signer, create_new_orders):
    new_oid = create_new_orders[0].order_id

    res = ob_client.get_order_by_client_id(CLIENT_OID)

    assert res.order_id == new_oid, "Order was not returned by orderId"
    assert res.client_order_id == CLIENT_OID, "Order was not returned by clientOrderId"
    assert res.price == "0.865", "Order was not returned by price"
    assert res.size == SIZE, "Order was not returned by size"
    assert res.side == "sell", "Order was not returned by side"
    assert res.symbol == SYMBOL, "Order was not returned by symbol"


def test_get_orders_for_user(ob_client, ob_signer, create_new_orders):
    res = ob_client.get_orders_for_user(page=1, page_size=25)

    assert len(res.data) > 0, "No orders returned"
    assert res.data[0]["orderId"], "orderId not returned"
    assert res.data[0]["clientOrderId"], "clientOrderId not returned"
    assert res.data[0]["price"], "price not returned"
    assert res.data[0]["size"], "size not returned"
    assert res.data[0]["side"], "side not returned"
    assert res.data[0]["symbol"], "symbol not returned"
