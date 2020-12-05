import React, {useEffect} from "react";
import {Observer} from "mobx-react";
import {Layout} from "../component/layout";
import {useHistory, useParams} from "react-router-dom";
import {useSession} from "../app";
import {BigOrdersModel, OrdersModel} from "../../api/orders";
import {Button, Checkbox, Divider, Form, Grid, Header, Icon, Input, List, Message, Table} from "semantic-ui-react";
import {Session} from "../../store/session";

function renderForm(session: Session, form: BigOrdersModel, i: number) {
    return <Table.Row key={`${form.form_name}-${form.total}-${form.built}`}>
        <Table.Cell width="1">{i}</Table.Cell>
        <Table.Cell width="13">{form.form_name}</Table.Cell>
        <Table.Cell width="1">{form.total - form.built}</Table.Cell>
        <Table.Cell width="1">
            <Checkbox onChange={ev => {session.completedBoxes[i] = !session.completedBoxes[i]}} checked={session.completedBoxes[i] ?? false}/>
        </Table.Cell>
    </Table.Row>;
}

function renderOrder(order: OrdersModel | null, forms: BigOrdersModel[], history: ReturnType<typeof useHistory>, session: Session) {
    if (order == null) {
        return <Message warning>
            Заказ не найден
        </Message>;
    }

    const addBox = (ev: React.KeyboardEvent<HTMLInputElement>) => {
        const el = ev.currentTarget as HTMLInputElement;
        if (ev.key === "Enter" && el.value.trim() !== "") {
            if (session.preparedBoxes.some(v => v === el.value.trim())) {
                el.value = "";
                return;
            }
            session.preparedBoxes.push(el.value.trim());
            el.value = "";
        }
    }

    const sendOrder = async () => {
        await session.finishOrders();
        history.push("/orders");
    };

    return <Grid>
        <Grid.Row>
            <Grid.Column width={6}>
                <Header sub>Сборные короба для заказа:</Header> {order.order_caption}
                <br />
                <br />
                <br />
                <Form>
                    <Form.Field>
                        <label>Соберите коробку и отсканируйте штрих-код:</label>
                        <Input placeholder='202700030' onKeyPress={addBox} />
                    </Form.Field>
                </Form>
                <Header sub>Собрано коробов:</Header>
                <List>
                    {session.preparedBoxes.map(i => <List.Item key={i}>{i}</List.Item>)}
                </List>
                <Divider />
                <List>
                    <List.Item>Итого: {session.preparedBoxes.length}</List.Item>
                </List>
            </Grid.Column>
            <Grid.Column width={10}>
                <Header sub>Требуется собрать:</Header>
                <Table celled compact>
                    <Table.Header>
                        <Table.Row>
                            <Table.HeaderCell />
                            <Table.HeaderCell>Товар</Table.HeaderCell>
                            <Table.HeaderCell>К&nbsp;сбору</Table.HeaderCell>
                            <Table.HeaderCell>Собрано</Table.HeaderCell>
                        </Table.Row>
                    </Table.Header>

                    <Table.Body>
                        {forms.map(renderForm.bind(null, session))}
                    </Table.Body>
                </Table>

                {forms.some((f, i) => session.completedBoxes[i] !== true)
                    ? <Button disabled negative>Не все короба укомплектованы</Button>
                    : <Button positive onClick={sendOrder}>Сборные короба полностью укомплектованы</Button>
                }
            </Grid.Column>
        </Grid.Row>

    </Grid>;
}

export function OrdersSmallPage() {
    const {id} = useParams<{id: string}>();
    const session = useSession();
    const history = useHistory();

    useEffect(() => {
        session.preparedBoxes = [];
        session.completedBoxes = {};
        session.currentOrderId = parseInt(id);
        session.fetchSmallOrdersToBuild().catch(console.error);
    }, [session, id])

    return <Observer>{() =>
        <Layout>
            {renderOrder(session.findOrder(parseInt(id)), session.currentSmallOrder, history, session)}
        </Layout>
    }</Observer>;
}
