import React, {useEffect} from "react";
import {Observer} from "mobx-react";
import {Layout} from "../component/layout";
import {useHistory, useParams} from "react-router-dom";
import {Button, Grid, Header, Icon, Message, Table} from "semantic-ui-react";
import {useSession} from "../app";
import {BigOrdersModel, OrdersModel} from "../../api/orders";

function renderForm(form: BigOrdersModel) {
    return <Table.Row positive={form.total - form.built === 0} key={`${form.form_name}-${form.total}-${form.built}`}>
        <Table.Cell width="10">{form.form_name}</Table.Cell>
        <Table.Cell width="1">{form.total}</Table.Cell>
        <Table.Cell width="1">{form.built}</Table.Cell>
        <Table.Cell width="1">{form.total - form.built}</Table.Cell>
        <Table.Cell width="1" negative={form.total - form.built > 0}>
            {form.total - form.built === 0 ? <Icon name='checkmark' /> : <Icon name='close' />}
        </Table.Cell>
    </Table.Row>;
}

function renderOrder(order: OrdersModel | null, forms: BigOrdersModel[], history: ReturnType<typeof useHistory>) {
    if (order == null) {
        return <Message warning>
            Заказ не найден
        </Message>;
    }

    const createPallet = () => {
        history.push(`/orders/pallet/${order?.id}`);
    }

    return <Grid columns={2} >
        <Grid.Row>
            <Grid.Column>
                <Header sub>Заказ:</Header> {order.order_caption}
            </Grid.Column>
            <Grid.Column>
                <Header sub>Заказчик:</Header> {order.customer}
            </Grid.Column>
        </Grid.Row>
        <Grid.Row>
            <Grid.Column>
                <Header sub>Адрес:</Header> {order.address}
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
                            <Table.HeaderCell />
                            <Table.HeaderCell>Всего</Table.HeaderCell>
                            <Table.HeaderCell>Собр.</Table.HeaderCell>
                            <Table.HeaderCell>Ост.</Table.HeaderCell>
                            <Table.HeaderCell>Статус.</Table.HeaderCell>
                        </Table.Row>
                    </Table.Header>

                    <Table.Body>
                        {forms.map(renderForm)}
                    </Table.Body>
                </Table>
            </Grid.Column>
        </Grid.Row>
    </Grid>;
}

export function OrdersBigPage() {
    const {id} = useParams<{id: string}>();
    const session = useSession();
    const history = useHistory();

    useEffect(() => {
        session.currentOrderId = parseInt(id);
        session.fetchBigOrdersToBuild().catch(console.error);
    }, [session, id])

    return <Observer>{() =>
        <Layout>
            {renderOrder(session.findOrder(parseInt(id)), session.currentBigOrder, history)}
        </Layout>
    }</Observer>;
}
