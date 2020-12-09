import React, {useCallback, useContext, useEffect} from "react";
import {Observer} from "mobx-react";
import {Layout} from "../component/layout";
import {Link, useHistory, useParams} from "react-router-dom";
import {useSession} from "../app";
import {BigOrdersModel, BigPalletModel, OrdersModel} from "../../api/orders";
import {Button, Dimmer, Divider, Form, Grid, Header, Input, Loader, Message, Table} from "semantic-ui-react";
import {Session} from "../../store/session";
import {runInAction} from "mobx";

// let audio = new Audio();
function renderTypes(type: { type: BigOrdersModel, barcode: string | null }, i: number) {
    return <Table.Row key={i} positive={(type.barcode ?? "").length > 0}>
        <Table.Cell>{type.type.form_name}</Table.Cell>
    </Table.Row>
}

function renderOrder(order: OrdersModel | null, pallet: BigPalletModel, history: ReturnType<typeof useHistory>, session: Session) {
    const idp = `${session.currentOrderId}-${session.currentBigPalletOrder?.pallet_num}`;
    if (session.currentOrderId == null || session.bigPalletOrderMatches[idp] == null) {
        return;
    }

    const addBox = async (ev: React.KeyboardEvent<HTMLInputElement>) => {
        const el = ev.currentTarget as HTMLInputElement;
        if (ev.key === "Enter" && el.value.trim() !== "") {
            const barcode = el.value.trim();

            runInAction(() => {
                session.lastError = "";
                session.lastSuccess = "";
            });

            if (session.bigPalletOrderMatches[idp].some(v => v.barcode === barcode)) {
                el.value = "";
                session.lastError = "Такой штрих-код уже добавлен";
                return;
            }

            try {
                const type = await session.requestPalletType(barcode);
                if (type.success) {
                    if (session.matchPalletBarcode(type.type, barcode)) {
                        session.lastSuccess = `Отсканирован короб ${type.type} ${barcode}`;
                    } else {
                        alert("ОШИБКА! Отсканирован ошибочный короб!!!")
                        session.lastError = "Внимание! Отсканирован ошибочный короб!";
                    }
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
                barcodes: session.bigPalletOrderMatches[idp].filter(m => m.barcode).map(m => m.barcode ?? ""),
            });

            if (resp.success) {
                window.open(`/orders/pallet/${order?.id}/print/${pallet.pallet_num}`);
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
                    <Form error={session.lastError.length > 0} success={session.lastSuccess.length > 0}>
                        <Form.Field>
                            <label>Сканируйте штрих-код коробки:</label>
                            <Input placeholder='202700030' onKeyPress={addBox} type="number" min={1} autoFocus/>
                            <Message error>
                                {session.lastError}
                            </Message>
                            <Message success>
                                {session.lastSuccess}
                            </Message>
                        </Form.Field>
                    </Form>
                </Grid.Column>
                <Grid.Column>
                    <Header sub>Для заказа №:</Header> {order?.order_caption ?? "<ОТСУТСТВУЕТ>"}
                </Grid.Column>
            </Grid.Row>
        </Grid>
        <Table celled singleLine collapsing>
            <Table.Body>
                {session.bigPalletOrderMatches[idp]
                    .filter(v => v.barcode == null)
                    .map(renderTypes)}
            </Table.Body>
            <Table.Footer>
                <Table.Row>
                    <Table.HeaderCell>
                        Итого: {session.bigPalletOrderMatches[idp].filter(f => (f.barcode?.length ?? 0) > 0).length}
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

    const warnOnUnload = useCallback((ev: BeforeUnloadEvent) => {
        ev.preventDefault();
        ev.returnValue = "";
        return "";
    }, []);

    useEffect(() => {
        runInAction(() => {
            session.curPage = "orders-pallet";
            session.breadcrumbs = [
                {key: 'orders', content: 'Комплектование', as: Link, to: "/orders"},
                {key: 'big', content: `Короба №${id}`, as: Link, to: `/orders/big/${id}`},
                {key: 'pallet', content: `Паллета`, active: true},
            ];
            session.lastError = "";
            session.currentOrderId = parseInt(id);
            session.fetchOrdersToBuild().catch(console.error);
            session.fetchBigPallet().then(() => {
                session.clearPalletBarcode();
            }).catch(console.error);
        });

        window.addEventListener("beforeunload", warnOnUnload);
        return () => {
            window.removeEventListener("beforeunload", warnOnUnload);
            runInAction(() => {
                session.curPage = "none";
                session.currentBigPalletOrder = {pallet_num: 0, types: []};
            });
        }
    }, [session, id, warnOnUnload])

    return <Observer>{() =>
        <Layout>
            <Dimmer inverted active={session.findOrder(parseInt(id)) == null}>
                <Loader/>
            </Dimmer>

            {renderOrder(session.findOrder(parseInt(id)), session.currentBigPalletOrder, history, session)}
        </Layout>
    }</Observer>;
}
