import {Observer} from "mobx-react";
import React, {useEffect} from "react";
import {useHistory} from "react-router-dom";
import {useSession} from "../app";
import {Layout} from "../component/layout";
import {Header, Table} from "semantic-ui-react";
import {OrdersModel, SubOrderModel} from "../../api/orders";
import * as H from "history";

function renderRow(history: ReturnType<typeof useHistory>, order: OrdersModel) {
    const rows = [<Table.Row positive onClick={() => order.opened = !order.opened} key={order.id}>
        <Table.Cell width="1">{order.num}</Table.Cell>
        <Table.Cell width="3">{order.order_caption}</Table.Cell>
        <Table.Cell width="2">{order.customer}</Table.Cell>
        <Table.Cell width="3">{order.address}</Table.Cell>
        <Table.Cell width="1">{order.run}</Table.Cell>
        <Table.Cell width="1">{order.amount_pallets}</Table.Cell>
        <Table.Cell width="1">{order.amount_boxes}</Table.Cell>
    </Table.Row>];

    const next = (sub: SubOrderModel) => {
        if (sub.is_small) {
            history.push(`/orders/small/${order.id}`);
        } else {
            history.push(`/orders/big/${order.id}`);
        }
    }

    if (order.opened) {
        let n = 0;
        for (const sub of order.sub_orders) {
            rows.push(
                <Table.Row key={`${order.id}-${n}`} onClick={() => next(sub)}>
                    <Table.Cell />
                    <Table.Cell>{sub.order_caption}</Table.Cell>
                    <Table.Cell />
                    <Table.Cell />
                    <Table.Cell />
                    <Table.Cell>{sub.amount_pallets}</Table.Cell>
                    <Table.Cell>{sub.amount_boxes}</Table.Cell>
                </Table.Row>
            );
            n++;
        }
    }
    return rows;
}

export function OrdersPage() {
    const session = useSession();
    const history = useHistory();

    return <Observer>{() =>
        <Layout>
            <Header>Комплектование</Header>

            <Table celled selectable singleLine>
                <Table.Header>
                    <Table.Row>
                        <Table.HeaderCell>№</Table.HeaderCell>
                        <Table.HeaderCell>Заказ</Table.HeaderCell>
                        <Table.HeaderCell>Заказчик</Table.HeaderCell>
                        <Table.HeaderCell>Адрес</Table.HeaderCell>
                        <Table.HeaderCell>Тираж</Table.HeaderCell>
                        <Table.HeaderCell>Паллет</Table.HeaderCell>
                        <Table.HeaderCell>Коробок</Table.HeaderCell>
                    </Table.Row>
                </Table.Header>

                <Table.Body>
                    {session.ordersToBuild?.map?.(renderRow.bind(null, history))}
                </Table.Body>
            </Table>

        </Layout>
    }</Observer>;
}
