import {Observer} from "mobx-react";
import React from "react";
import { Link } from "react-router-dom";
import {Breadcrumb, Container, Divider, Icon, Segment} from "semantic-ui-react";
import {useSession} from "../app";

interface Props {
    children: any;
}
const sections = [
    { key: 'Home', content: 'Home', as: Link, to: "/" },
    { key: 'Store', content: 'Store', link: true },
    { key: 'Shirt', content: 'T-Shirt', active: true },
]

export function Layout({children}: Props) {
    const session = useSession();
    return <Observer>{() =>
        <Container>
            <Breadcrumb icon='right angle' sections={session.breadcrumbs} />
            <Segment.Group>
                <Segment textAlign="right">
                    <Icon name="user" color="blue" />{session.currentUser?.login}
                    <Icon name="time" color="green" style={{marginLeft: 10}} />{session.currentDate}
                </Segment>
                <Segment padded style={{minHeight: '100vh'}}>
                    {children}
                </Segment>
            </Segment.Group>
        </Container>}
    </Observer>;
}
