import React from "react";
import {Observer} from "mobx-react";
import {Layout} from "../component/layout";

export function OrdersSmallPage() {
    return <Observer>{() =>
        <Layout>
            orders small
        </Layout>
    }</Observer>;
}
