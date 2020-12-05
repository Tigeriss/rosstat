import {Observer} from "mobx-react";
import React from "react";
import {Container, Divider, Icon, Segment} from "semantic-ui-react";
import {useSession} from "../app";

interface Props {
    children: any;
}

export function Layout({children}: Props) {
    const session = useSession();
    return <Observer>{() =>
        <Container>
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
