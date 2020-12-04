import React, {useEffect} from "react";
import {Observer} from "mobx-react";
import {Layout} from "../component/layout";
import { useParams } from "react-router-dom";
import {Grid, Header} from "semantic-ui-react";

export function OrdersBigPage() {
    const {id} = useParams<{id: string}>();

    return <Observer>{() =>
        <Layout>
            <Header>Обработка</Header>

            <Grid columns={3} divided>
                <Grid.Row>
                    <Grid.Column>
                        asd
                    </Grid.Column>
                    <Grid.Column>
                        asd
                    </Grid.Column>
                    <Grid.Column>
                        asd
                    </Grid.Column>
                </Grid.Row>

                <Grid.Row>
                    <Grid.Column>
                        asd
                    </Grid.Column>
                    <Grid.Column>
                        asd
                    </Grid.Column>
                    <Grid.Column>
                        asd
                    </Grid.Column>
                </Grid.Row>
            </Grid>
            orders big {id}
        </Layout>
    }</Observer>;
}
