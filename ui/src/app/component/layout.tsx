import React from "react";
import {Container, Segment} from "semantic-ui-react";

interface Props {
    children: any;
}

export function Layout({children}: Props) {
    return <Container>
        <Segment padded style={{minHeight: '100vh'}}>
            {children}
        </Segment>
    </Container>;
}
