import React, {useState} from "react";
import {Observer} from "mobx-react";
import {Layout} from "../component/layout";
import {Form, Header, Table} from "semantic-ui-react";
import {useSession} from "../app";
import {useHistory} from "react-router-dom";
import {Session} from "../../store/session";
import {OrdersModel, SubOrderModel} from "../../api/orders";
import {ShipmentModel} from "../../api/shipment";


function renderRow(history: ReturnType<typeof useHistory>, session: Session, order: ShipmentModel) {
    const next = () => {
        history.push(`/shipment/pallet/${order.id}`);
    }

    return <Table.Row warning onClick={next} key={order.id}>
        <Table.Cell width="1">{order.num}</Table.Cell>
        <Table.Cell width="3">{order.order_caption}</Table.Cell>
        <Table.Cell width="2">{order.customer}</Table.Cell>
        <Table.Cell width="3">{order.address}</Table.Cell>
        <Table.Cell width="1">{order.run}</Table.Cell>
        <Table.Cell width="1">{order.amount_pallets}</Table.Cell>
        <Table.Cell width="1">{order.amount_boxes}</Table.Cell>
    </Table.Row>;
}

export function ShipmentPage() {
    const session = useSession();
    const history = useHistory();
    const [filter, setFilter] = useState("");
    const normFilter = filter.trim().toLocaleLowerCase();

    return <Observer>{() =>
        <Layout>
            <Header>Отгрузка</Header>

            <Form>
                <Form.Group>
                    <Form.Field>
                        <label>Фильтр:</label>
                        <input type="text" value={filter} onChange={e => setFilter(e.target.value)}/>
                    </Form.Field>
                </Form.Group>
            </Form>


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
                    {session.currentShipments?.filter(o => normFilter.length === 0 || o.order_caption.toLowerCase().includes(filter))
                        .map(renderRow.bind(null, history, session))}
                </Table.Body>
            </Table>

        </Layout>
    }</Observer>;
}
