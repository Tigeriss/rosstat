import React, {useState} from "react";
import {Button, Container, Form, Header, Message, Segment} from "semantic-ui-react";
import {useSession} from "../app";
import {useHistory} from "react-router-dom";

const styles: React.CSSProperties = {
    marginTop: "200px",
};

export function Login() {
    const [login, setLogin] = useState("");
    const [password, setPassword] = useState("");
    const [error, setError] = useState<string | null>(null);
    const session = useSession();
    const history = useHistory();

    const signin = async () => {
        setError(null);
        if (await session.login(login, password)) {
            history.replace("/");
        } else {
            setError("Неправильный логин или пароль");
        }
    }

    return <Container style={styles}>
        <Segment padded="very">
            <Header>Вход</Header>
            <Form error={error != null}>
                <Form.Field>
                    <label>Имя пользователя</label>
                    <input placeholder='Имя пользователя' value={login} onChange={(e) => setLogin(e.target.value)} />
                </Form.Field>
                <Form.Field>
                    <label>Пароль</label>
                    <input placeholder='Пароль' type="password" value={password} onChange={(e) => setPassword(e.target.value)} />
                </Form.Field>
                <Button onClick={signin}>Войти</Button>
                <Message error content={error}/>
            </Form>
        </Segment>
    </Container>;
}
