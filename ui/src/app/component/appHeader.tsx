import React from "react";
import {Header, Icon, Menu} from "semantic-ui-react";
import {Link, useLocation} from "react-router-dom";
import {useSession} from "../app";

const menuStyle: React.CSSProperties = {
    backgroundColor: "#fff",
    border: "1px solid #ddd",
    boxShadow: "0px 3px 5px rgba(0, 0, 0, 0.2)"
}

export function AppHeader() {
    const location = useLocation();
    const session = useSession();

    return <Menu style={menuStyle} stackable icon="labeled">

        <Menu.Item as={Link} to="/" header>
            <Header size="huge"> Росстат </Header>
        </Menu.Item>

        <Menu.Menu position='right'>
            <Menu.Item as={Link} to="/admin" active={location.pathname === "/admin"}>
                <Icon name="users"/> Пользователи
            </Menu.Item>
            <Menu.Item onClick={() => session.currentUser = null}>
                <Icon name="log out"/> Выход
            </Menu.Item>
        </Menu.Menu>
    </Menu>
}
