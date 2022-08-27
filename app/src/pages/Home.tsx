import {IonContent, IonPage, useIonViewWillEnter} from '@ionic/react';
import React, {useState} from "react";
import Header from "../components/Header";
import DataCard from "../components/DataCard";
import Client from "../service/Client";

const Home: React.FC = () => {
    const [client] = useState<Client>(Client.getInstance());
    let [count, setCount] = useState(0);

    useIonViewWillEnter(() => {
        setCount(count++);
    });

    return (
        <IonPage>
            <Header title="ax87"/>
            <IonContent fullscreen color="light">
                <DataCard data={client.latest} title="Latest"/>
                <DataCard data={client.median} title="Median"/>
                <DataCard data={client.max} title="Max"/>
                <DataCard data={client.min} title="Min"/>
            </IonContent>
        </IonPage>
    );
};

export default Home;
