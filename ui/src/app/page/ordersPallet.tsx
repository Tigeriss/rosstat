import React, {useEffect} from "react";
import {Observer} from "mobx-react";
import {Layout} from "../component/layout";
import {useHistory, useParams} from "react-router-dom";
import {useSession} from "../app";
import {BigOrdersModel, BigPalletModel, OrdersModel} from "../../api/orders";
import {Button, Divider, Form, Grid, Header, Input, Message, Table} from "semantic-ui-react";
import {Session} from "../../store/session";

function renderTypes(type: {type: BigOrdersModel, barcode: string | null}, i: number) {
    return <Table.Row key={i}>
        <Table.Cell>{type.type.form_name}</Table.Cell>
        <Table.Cell width={3}>{type.barcode}</Table.Cell>
    </Table.Row>
}

function renderOrder(order: OrdersModel | null, pallet: BigPalletModel, history: ReturnType<typeof useHistory>, session: Session) {
    if (order == null) {
        return <Message warning>
            Заказ не найден
        </Message>;
    }

    const addBox = async (ev: React.KeyboardEvent<HTMLInputElement>) => {
        const el = ev.currentTarget as HTMLInputElement;
        if (ev.key === "Enter" && el.value.trim() !== "") {
            const barcode = el.value.trim();
            session.lastError = "";
            if (session.bigPalletOrderMatches.some(v => v.barcode === barcode)) {
                el.value = "";
                session.lastError = "Такой штрих-код уже добавлен";
                return;
            }

            try {
                const type = await session.requestPalletType(barcode);
                if (type.success) {
                    session.matchPalletBarcode(type.type, barcode);
                } else {
                    session.lastError = type.error ?? "";
                }
            } catch (ex) {
                console.log(ex);
                session.lastError = "Произошла ошибка при добавлении штрих-кода";
            }
            el.value = "";
        }
    }

    const createPallet = async () => {
        try {
            session.lastError = "";
            const resp = await session.finishBigPallet({
                pallet_num: pallet.pallet_num,
                barcodes: session.bigPalletOrderMatches.filter(m => m.barcode).map(m => m.barcode ?? ""),
            });

            if (resp.success) {
                window.open(`/orders/pallet/${order.id}/print/${pallet.pallet_num}`);
                if (resp.last_pallet) {
                    history.push(`/orders`);
                } else {
                    history.push(`/orders/big/${order?.id}`);
                }
            } else {
                session.lastError = resp.error ?? "";
            }

        } catch (ex) {
            session.lastError = "Произошла ошибка при добавлении паллеты";
        }
    }

    return <>
        <Grid columns={3}>
            <Grid.Row>
                <Grid.Column>
                    <Header sub>Паллета №:</Header> {pallet.pallet_num}<br/><br/>
                    <Form error={session.lastError.length > 0}>
                        <Form.Field>
                            <label>Сканируйте штрих-код коробки:</label>
                            <Input placeholder='202700030' onKeyPress={addBox}/>
                            <Message error>
                                {session.lastError}
                            </Message>
                        </Form.Field>
                    </Form>
                </Grid.Column>
                <Grid.Column>
                    <Header sub>Для заказа №:</Header> {order.order_caption}
                </Grid.Column>
            </Grid.Row>
        </Grid>
        <Table celled singleLine collapsing>
            <Table.Body>
                {session.bigPalletOrderMatches.map(renderTypes)}
            </Table.Body>
            <Table.Footer>
                <Table.Row>
                    <Table.HeaderCell>
                        Итого:
                    </Table.HeaderCell>
                    <Table.HeaderCell>
                        {session.bigPalletOrderMatches.filter(f => (f.barcode?.length ?? 0) > 0).length}
                    </Table.HeaderCell>
                </Table.Row>
            </Table.Footer>
        </Table>
        <Button primary onClick={createPallet}>Паллета собрана</Button>
    </>;
}


export function OrdersPalletPage() {
    const {id} = useParams<{ id: string }>();
    const session = useSession();
    const history = useHistory();

    useEffect(() => {
        session.lastError = "";
        session.currentOrderId = parseInt(id);
        session.fetchBigPallet().then(() => {
            session.clearPalletBarcode();
        }).catch(console.error);
    }, [session, id])

    return <Observer>{() =>
        <Layout>
            {renderOrder(session.findOrder(parseInt(id)), session.currentBigPalletOrder, history, session)}
        </Layout>
    }</Observer>;
}
