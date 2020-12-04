import {Observer} from "mobx-react";
import React from "react";
import {Container, Segment} from "semantic-ui-react";
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
                    {session.currentUser?.login} {session.currentDate}
                </Segment>
                <Segment padded style={{minHeight: '100vh'}}>
                    {children}
                </Segment>
            </Segment.Group>
        </Container>}
    </Observer>;
}
