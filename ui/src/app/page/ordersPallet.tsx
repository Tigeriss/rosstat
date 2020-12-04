import React from "react";
import {Observer} from "mobx-react";
import {Layout} from "../component/layout";

export function OrdersPalletPage() {
    return <Observer>{() =>
        <Layout>
            orders pallet
        </Layout>
    }</Observer>;
}
