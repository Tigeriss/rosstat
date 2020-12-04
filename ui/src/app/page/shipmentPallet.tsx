import React from "react";
import {Observer} from "mobx-react";
import {Layout} from "../component/layout";

export function ShipmentPalletPage() {
    return <Observer>{() =>
        <Layout>
            shipment pallet
        </Layout>
    }</Observer>;
}
