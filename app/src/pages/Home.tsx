import {IonContent, IonPage} from '@ionic/react';
import React from "react";
import Header from "../components/Header";
import DataCard from "../components/DataCard";
import Data from "../model/Data";

const Home: React.FC = () => {
    return (
        <IonPage>
            <Header title="ax87"/>
            <IonContent fullscreen color="light">
                <DataCard data={Data.build()} title="Latest"/>
                <DataCard data={Data.build()} title="Median"/>
                <DataCard data={Data.build()} title="Max"/>
                <DataCard data={Data.build()} title="Min"/>
            </IonContent>
        </IonPage>
    );
};

export default Home;
