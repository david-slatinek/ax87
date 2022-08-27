import {IonContent, IonList, IonPage, useIonViewWillEnter} from '@ionic/react';
import React, {useState} from "react";
import Header from "../components/Header";
import DataCard from "../components/DataCard";
import Client from "../service/Client";

const View: React.FC = () => {
    const client = Client.getInstance();
    let [count, setCount] = useState(0);

    useIonViewWillEnter(() => {
        setCount(count++);
    });

    return (
        <IonPage>
            <Header title="Last 24H"/>
            <IonContent fullscreen color="light">
                <IonList>
                    {
                        client.today.map((d, index) => (
                            <DataCard key={index} data={d}/>
                        ))
                    }
                </IonList>
            </IonContent>
        </IonPage>
    );
};

export default View;
