import React, {useEffect, useState} from "react";
import {Observer} from "mobx-react";
import {Layout} from "../component/layout";
import {Dimmer, Form, Header, Loader, Table} from "semantic-ui-react";
import {useSession} from "../app";
import {useHistory} from "react-router-dom";
import {Session} from "../../store/session";
import {OrdersModel, SubOrderModel} from "../../api/orders";
import {ShipmentModel} from "../../api/shipment";
import {runInAction} from "mobx";


function renderRow(history: ReturnType<typeof useHistory>, session: Session, order: ShipmentModel) {
    const next = () => {
        if (order.shipped) {
            window.open(`/shipment/print/${order.id}`);
        } else {
            history.push(`/shipment/pallet/${order.id}`);
        }
    }

    return <Table.Row warning={order.shipped} onClick={next} key={order.id} disabled={!order.collected && !order.shipped}>
        <Table.Cell>{order.num}</Table.Cell>
        <Table.Cell>{order.order_caption}</Table.Cell>
        <Table.Cell>{order.customer}</Table.Cell>
        <Table.Cell>{order.address}</Table.Cell>
        <Table.Cell>{order.run}</Table.Cell>
        <Table.Cell>{order.amount_pallets}</Table.Cell>
        <Table.Cell>{order.amount_boxes}</Table.Cell>
    </Table.Row>;
}

export function ShipmentPage() {
    const session = useSession();
    const history = useHistory();
    const [filter, setFilter] = useState("");
    const normFilter = filter.trim().toLocaleLowerCase();

    useEffect(() => {
        runInAction(() => {
            session.curPage = "shipment";
            session.breadcrumbs = [
                {key: 'shipment', content: 'Отгрузка', active: true},
            ];
            session.fetchShipmentReady().catch(console.error);
        });

        return () => {
            session.curPage = "none";
        }
    }, [session]);

    return <Observer>{() =>
        <Layout>
            <Dimmer inverted active={(session.currentShipments?.length ?? 0) === 0}>
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


            <Table celled selectable singleLine>
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
                    {session.currentShipments?.filter(o => normFilter.length === 0 || o.order_caption.toLowerCase().includes(filter))
                        .map(renderRow.bind(null, history, session))}
                </Table.Body>
            </Table>

        </Layout>
    }</Observer>;
}
