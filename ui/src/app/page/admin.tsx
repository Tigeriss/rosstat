import React, {useEffect, useState} from "react";
import {Layout} from "../component/layout";
import {Button, Form, Header, Icon, List, Table} from "semantic-ui-react";
import {useSession} from "../app";
import {Observer} from "mobx-react";
import {User} from "../../api/admin";
import {Session} from "../../store/session";

const options = [
    {key: 'a', text: 'Админ', value: 'admin'},
    {key: 'c', text: 'Коллектор', value: 'collector'},
    {key: 's', text: 'Кладовщик', value: 'storekeeper'},
]

function renderUser(session: Session, user: User) {
    return <Table.Row key={user.login}>
        <Table.Cell>{user.login}</Table.Cell>
        <Table.Cell>{options.find(o => o.value === user.role)?.text}</Table.Cell>
        <Table.Cell>
            <Icon name="close" color="red" style={{cursor: "pointer"}} onClick={() => session.deleteUser(user.login)}/>
        </Table.Cell>
    </Table.Row>
}

export function Admin() {
    const session = useSession();
    const [user, setUser] = useState<User>({login: "", password: "", role: "admin"});

    useEffect(() => {
        session.curPage = "admin";
        session.fetchUsers().catch(console.error);
        return () => {
            session.curPage = "none";
        }
    }, [session]);

    const addUser = async () => {
        await session.addUser(user);
        setUser({login: "", password: "", role: "admin"});
    };

    return <Observer>{() => <Layout>
        <Header>Пользователи</Header>
        <Table celled striped>
            <Table.Header>
                <Table.Row>
                    <Table.HeaderCell>Логин</Table.HeaderCell>
                    <Table.HeaderCell>Роль</Table.HeaderCell>
                    <Table.HeaderCell/>
                </Table.Row>
            </Table.Header>

            <Table.Body>
                {session.users.map(renderUser.bind(null, session))}
            </Table.Body>
        </Table>

        <Form onSubmit={addUser}>
            <Form.Group widths='equal'>
                <Form.Field>
                    <label>Логин</label>
                    <input placeholder='login' value={user.login} onChange={e => setUser({...user, login: e.target.value})}/>
                </Form.Field>
                <Form.Field>
                    <label>Пароль</label>
                    <input type="password" placeholder='password' value={user.password}  onChange={e => setUser({...user, password: e.target.value})}/>
                </Form.Field>
                <Form.Field>
                    <Form.Select
                        fluid
                        label='Роль'
                        options={options}
                        value={user.role}
                        onChange={(ev, t) => setUser({...user, role: t.value as string})}
                    />
                </Form.Field>
            </Form.Group>
            <Form.Button type="submit">Добавить</Form.Button>
        </Form>
    </Layout>
    }</Observer>;
}
