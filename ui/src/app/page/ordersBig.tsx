import React, {useEffect} from "react";
import {Observer} from "mobx-react";
import {Layout} from "../component/layout";
import {Link, useHistory, useParams} from "react-router-dom";
import {Button, Dimmer, Grid, Header, Icon, Loader, Message, Table} from "semantic-ui-react";
import {useSession} from "../app";
import {BigOrdersModel, OrdersModel} from "../../api/orders";
import {runInAction} from "mobx";

function renderForm(form: BigOrdersModel) {
    return <Table.Row positive={form.total - form.built === 0} key={`${form.type}`}>
        <Table.Cell width="10">{form.form_name}</Table.Cell>
        <Table.Cell width="1">{form.total}</Table.Cell>
        <Table.Cell width="1">{form.built}</Table.Cell>
        <Table.Cell width="1">{form.total - form.built}</Table.Cell>
        <Table.Cell width="1" negative={form.total - form.built > 0}>
            {form.total - form.built === 0 ? <Icon name='checkmark' color="green"/> : <Icon name='close' color="red"/>}
        </Table.Cell>
    </Table.Row>;
}

function renderOrder(order: OrdersModel | null, forms: BigOrdersModel[], history: ReturnType<typeof useHistory>) {
    const createPallet = () => {
        history.push(`/orders/pallet/${order?.id}`);
    }

    return <Grid columns={2}>
        <Grid.Row>
            <Grid.Column>
                <Header sub>Заказ:</Header> {order?.order_caption ?? "<ОТСУТСТВУЕТ>"}
            </Grid.Column>
            <Grid.Column>
                <Header sub>Заказчик:</Header> {order?.customer ?? "<ОТСУТСТВУЕТ>"}
            </Grid.Column>
        </Grid.Row>
        <Grid.Row>
            <Grid.Column>
                <Header sub>Адрес:</Header> {order?.address ?? "<ОТСУТСТВУЕТ>"}
            </Grid.Column>
            <Grid.Column>
                <Button primary onClick={createPallet}>Создать паллету</Button>
            </Grid.Column>
        </Grid.Row>
        <Grid.Row>
            <Grid.Column width={2}>
                <Table celled singleLine>
                    <Table.Header>
                        <Table.Row>
                            <Table.HeaderCell/>
                            <Table.HeaderCell width={1}>Всего</Table.HeaderCell>
                            <Table.HeaderCell width={1}>Собр.</Table.HeaderCell>
                            <Table.HeaderCell width={1}>Ост.</Table.HeaderCell>
                            <Table.HeaderCell width={1}>Статус.</Table.HeaderCell>
                        </Table.Row>
                    </Table.Header>

                    <Table.Body>
                        {forms?.map(renderForm)}
                    </Table.Body>
                </Table>
            </Grid.Column>
        </Grid.Row>
    </Grid>;
}

export function OrdersBigPage() {
    const {id} = useParams<{ id: string }>();
    const session = useSession();
    const history = useHistory();

    useEffect(() => {
        runInAction(() => {
            session.curPage = "orders-big";
            session.breadcrumbs = [
                {key: 'orders', content: 'Комплектование', as: Link, to: "/orders"},
                {key: 'big', content: `Короба №${id}`, active: true},
            ];
            session.currentOrderId = parseInt(id);
            session.fetchOrdersToBuild().catch(console.error);
            session.fetchBigOrdersToBuild().catch(console.error);
        });

        return () => {
            session.curPage = "none";
        }
    }, [session, id])

    return <Observer>{() =>
        <Layout>
            <Dimmer inverted active={session.findOrder(parseInt(id)) == null}>
                <Loader/>
            </Dimmer>

            {renderOrder(session.findOrder(parseInt(id)), session.currentBigOrder, history)}
        </Layout>
    }</Observer>;
}
