import {Observer} from "mobx-react";
import React, {useEffect, useState} from "react";
import {Link, useHistory} from "react-router-dom";
import {useSession} from "../app";
import {Layout} from "../component/layout";
import {Dimmer, Form, Header, Loader, Table} from "semantic-ui-react";
import {OrdersModel, SubOrderModel} from "../../api/orders";
import {Session} from "../../store/session";
import {runInAction} from "mobx";

function renderRow(history: ReturnType<typeof useHistory>, session: Session, order: OrdersModel) {
    const rows = [<Table.Row warning
                             onClick={() => session.openedOrders[order.id] = !session.openedOrders[order.id]}
                             key={order.id}>
        <Table.Cell>{order.num}</Table.Cell>
        <Table.Cell singleLine>{order.order_caption}</Table.Cell>
        <Table.Cell>{order.customer}</Table.Cell>
        <Table.Cell>{order.address}</Table.Cell>
        <Table.Cell>{order.run}</Table.Cell>
        <Table.Cell>{order.amount_pallets}</Table.Cell>
        <Table.Cell>{order.amount_boxes}</Table.Cell>
    </Table.Row>];

    const next = (sub: SubOrderModel) => {
        if (sub.is_small) {
            if (sub.amount_boxes === 0) {
                history.push(`/orders/small/${order.id}`);
            }
        } else {
            history.push(`/orders/big/${order.id}`);
        }
    }

    if (session.openedOrders[order.id]) {
        let n = 0;
        for (const sub of order.sub_orders) {
            rows.push(
                <Table.Row disabled={sub.is_small && sub.amount_boxes > 0}
                           key={`${order.id}-${n}`} onClick={() => next(sub)}>
                    <Table.Cell/>
                    <Table.Cell>{sub.order_caption}</Table.Cell>
                    <Table.Cell/>
                    <Table.Cell/>
                    <Table.Cell/>
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
    const [filter, setFilter] = useState("");
    const normFilter = filter.trim().toLocaleLowerCase();

    useEffect(() => {
        runInAction(() => {
            session.curPage = "orders";
            session.breadcrumbs = [
                {key: 'orders', content: 'Комплектование', active: true},
            ];
            session.fetchOrdersToBuild().catch(console.error);
        });

        return () => {
            session.curPage = "none";
        }
    }, [session]);

    return <Observer>{() =>
        <Layout>
            <Dimmer inverted active={(session.ordersToBuild?.length ?? 0) === 0}>
                <Loader/>
            </Dimmer>

            <Form>
                <Form.Group>
                    <Form.Field>
                        <label>Фильтр:</label>
                        <input type="text" value={filter} onChange={e => setFilter(e.target.value)}/>
                    </Form.Field>
                </Form.Group>
            </Form>

            <Table celled selectable>
                <Table.Header>
                    <Table.Row>
                        <Table.HeaderCell width="1">№</Table.HeaderCell>
                        <Table.HeaderCell width="3">Заказ</Table.HeaderCell>
                        <Table.HeaderCell width="2">Заказчик</Table.HeaderCell>
                        <Table.HeaderCell width="3">Адрес</Table.HeaderCell>
                        <Table.HeaderCell width="1">Тираж</Table.HeaderCell>
                        <Table.HeaderCell width="1">Паллет</Table.HeaderCell>
                        <Table.HeaderCell width="1">Коробок</Table.HeaderCell>
                    </Table.Row>
                </Table.Header>

                <Table.Body>
                    {session.ordersToBuild?.filter(o => normFilter.length === 0 || o.order_caption.toLowerCase().includes(filter))
                        .map(renderRow.bind(null, history, session))}
                </Table.Body>
            </Table>

        </Layout>
    }</Observer>;
}
