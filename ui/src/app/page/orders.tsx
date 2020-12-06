import {Observer} from "mobx-react";
import React, {useState} from "react";
import {useHistory} from "react-router-dom";
import {useSession} from "../app";
import {Layout} from "../component/layout";
import {Form, Header, Table} from "semantic-ui-react";
import {OrdersModel, SubOrderModel} from "../../api/orders";
import {Session} from "../../store/session";

function renderRow(history: ReturnType<typeof useHistory>, session: Session, order: OrdersModel) {
    const rows = [<Table.Row warning
                             onClick={() => session.openedOrders[order.id] = !session.openedOrders[order.id]}
                             key={order.id}>
        <Table.Cell width="1">{order.num}</Table.Cell>
        <Table.Cell width="3" singleLine>{order.order_caption}</Table.Cell>
        <Table.Cell width="2">{order.customer}</Table.Cell>
        <Table.Cell width="3">{order.address}</Table.Cell>
        <Table.Cell width="1">{order.run}</Table.Cell>
        <Table.Cell width="1">{order.amount_pallets}</Table.Cell>
        <Table.Cell width="1">{order.amount_boxes}</Table.Cell>
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
                <Table.Row disabled={sub.amount_boxes > 0}
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

    return <Observer>{() =>
        <Layout>
            <Header>Комплектование</Header>

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
                    {session.ordersToBuild?.filter(o => normFilter.length === 0 || o.order_caption.toLowerCase().includes(filter))
                        .map(renderRow.bind(null, history, session))}
                </Table.Body>
            </Table>

        </Layout>
    }</Observer>;
}
