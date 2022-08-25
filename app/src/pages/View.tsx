import {IonContent, IonList, IonPage} from '@ionic/react';
import React from "react";
import Header from "../components/Header";
import Data from "../model/Data";
import DataCard from "../components/DataCard";

const View: React.FC = () => {
    let data: Data[] = [Data.build(), Data.build(), Data.build(), Data.build()];

    return (
        <IonPage>
            <Header title="Last 24H"/>
            <IonContent fullscreen color="light">
                <IonList>
                    {
                        data.map((d, index) => (
                            <DataCard key={index} data={d}/>
                        ))
                    }
                </IonList>
            </IonContent>
        </IonPage>
    );
};

export default View;
