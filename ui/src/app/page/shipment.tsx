import React from "react";
import {Observer} from "mobx-react";
import {Layout} from "../component/layout";

export function ShipmentPage() {
    return <Observer>{() =>
        <Layout>
            shipment
        </Layout>
    }</Observer>;
}
