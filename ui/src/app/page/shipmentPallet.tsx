import React, {useEffect} from "react";
import {Observer} from "mobx-react";
import {Layout} from "../component/layout";
import {useSession} from "../app";
import {Link, useHistory, useParams} from "react-router-dom";
import {Dimmer, Form, Grid, Header, Icon, Input, Loader, Message, Table} from "semantic-ui-react";
import {ShipmentModel, ShipmentPalletModel} from "../../api/shipment";
import {Session} from "../../store/session";
import {runInAction} from "mobx";

function renderLabels(session: Session, pallet: ShipmentPalletModel) {
    return <Table.Row key={pallet.barcode} positive={session.sentPallets[pallet.barcode]}>
        <Table.Cell>
            {pallet.num}
        </Table.Cell>
        <Table.Cell>
            {pallet.pallet_num}
        </Table.Cell>
        <Table.Cell>
            {pallet.barcode}
        </Table.Cell>
        <Table.Cell>
            {pallet.amount_boxes}
        </Table.Cell>
        <Table.Cell>
            {session.sentPallets[pallet.barcode]
                ?  <Icon name='checkmark' color="green" />
                :  <Icon name='close' color="red" />
            }
        </Table.Cell>
    </Table.Row>;
}

function renderShipment(order: ShipmentModel | null, pallet: ShipmentPalletModel[], history: ReturnType<typeof useHistory>, session: Session) {
    const addLabel = async (ev: React.KeyboardEvent<HTMLInputElement>) => {
        session.lastError = "";
        const el = ev.currentTarget as HTMLInputElement;
        if (ev.key === "Enter" && el.value.trim() !== "") {
            const bar = el.value.trim();
            el.value = "";

            if (Object.keys(session.sentPallets).some(v => v === bar)) {
                session.lastError = "Паллета с данным штрих-кодом уже отсканирована";
                return;
            }

            if (!pallet.find(p => p.barcode === bar)) {
                session.lastError = "Паллета с данным штрих-кодом не найдена";
                return;
            }

            session.sentPallets[bar] = true;

            if (Object.values(session.sentPallets).filter(v => v).length === pallet.length) {
                try {
                    await session.finishPalletShipment();
                    history.push("/shipment");
                } catch (ex) {
                    console.error(ex);
                    session.lastError = "Произошла ошибка при сохранении данных";
                }
            }
        }
    };

    return <Grid>
        <Grid.Row>
            <Grid.Column width={4}>
                <Header sub>Заказ №:</Header> {order?.order_caption ?? "<ОТСУТСТВУЕТ>"}
            </Grid.Column>
            <Grid.Column width={4}>
                <Header sub>Паллет:</Header> {pallet.length}
            </Grid.Column>
            <Grid.Column width={4}>
                <Header sub>Отгруженно:</Header> {Object.values(session.sentPallets).filter(v => v).length}
            </Grid.Column>
            <Grid.Column width={4}>
                <Header sub>Статус:</Header> не отгруженно
            </Grid.Column>
        </Grid.Row>
        <Grid.Row>
            <Grid.Column width={6}>
                <Header sub>Адрес:</Header> {order?.address ?? "<ОТСУТСТВУЕТ>"}
                <br/>
                <br/>
                <br/>
                <Form error={session.lastError.length > 0}>
                    <Form.Field>
                        <label>Соберите коробку и отсканируйте штрих-код:</label>
                        <Input placeholder='202700030' onKeyPress={addLabel} autoFocus/>
                        <Message error>
                            {session.lastError}
                        </Message>
                    </Form.Field>
                </Form>
            </Grid.Column>
            <Grid.Column width={10}>
                <Table celled singleLine collapsing>
                    <Table.Header>
                        <Table.Row>
                            <Table.HeaderCell width={1}>
                                №
                            </Table.HeaderCell>
                            <Table.HeaderCell width={1}>
                                Паллета
                            </Table.HeaderCell>
                            <Table.HeaderCell width={3}>
                                Штрих-код
                            </Table.HeaderCell>
                            <Table.HeaderCell width={1}>
                                Коробов
                            </Table.HeaderCell>
                            <Table.HeaderCell width={1}>
                                Статус
                            </Table.HeaderCell>
                        </Table.Row>
                    </Table.Header>
                    <Table.Body>
                        {pallet.map(renderLabels.bind(null, session))}
                    </Table.Body>
                    <Table.Footer>
                        <Table.Row>
                            <Table.HeaderCell colSpan={3}>
                                Итого:
                            </Table.HeaderCell>
                            <Table.HeaderCell>
                                {pallet.filter(p => session.sentPallets[p.barcode]).reduce((l, r) => l + r.amount_boxes, 0)}
                            </Table.HeaderCell>
                            <Table.HeaderCell>
                                {pallet.filter(p => session.sentPallets[p.barcode]).length}
                            </Table.HeaderCell>
                        </Table.Row>
                    </Table.Footer>
                </Table>
            </Grid.Column>
        </Grid.Row>
    </Grid>;
}

export function ShipmentPalletPage() {
    const {id} = useParams<{ id: string }>();
    const session = useSession();
    const history = useHistory();

    useEffect(() => {
        runInAction(() => {
            session.curPage = "shipment-pallet";
            session.breadcrumbs = [
                {key: 'shipment', content: 'Комплектование', as: Link, to: "/shipment"},
                {key: 'big', content: `Заказ №${id}`, active: true},
            ];
            session.lastError = "";
            session.sentPallets = {};
            session.currentShipmentId = parseInt(id);
            session.fetchShipmentReady().catch(console.error);
            session.fetchShipmentPallet().catch(console.error);
        });

        return () => {
            session.curPage = "none";
        }
    }, [session, id])

    return <Observer>{() =>
        <Layout>
            <Dimmer inverted active={session.findShipment(parseInt(id)) == null}>
                <Loader/>
            </Dimmer>

            {renderShipment(session.findShipment(parseInt(id)), session.currentShipmentPallet, history, session)}
        </Layout>
    }</Observer>;
}
