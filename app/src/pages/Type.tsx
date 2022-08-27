import {IonContent, IonItem, IonLabel, IonList, IonPage, IonRadio, IonRadioGroup} from '@ionic/react';
import React from "react";
import Header from "../components/Header";
import Client from "../service/Client";

const Type: React.FC = () => {
    const client = Client.getInstance();

    return (
        <IonPage>
            <Header title="Type"/>
            <IonContent fullscreen color="light">
                <IonList>
                    <IonRadioGroup value={client.type} onIonChange={e => {
                        client.type = e.detail.value;
                    }}>
                        <IonItem>
                            <IonLabel>Carbon monoxide</IonLabel>
                            <IonRadio slot="start" value="CARBON_MONOXIDE"/>
                        </IonItem>

                        <IonItem>
                            <IonLabel>Air quality</IonLabel>
                            <IonRadio slot="start" value="AIR_QUALITY"/>
                        </IonItem>

                        <IonItem>
                            <IonLabel>Raindrops</IonLabel>
                            <IonRadio slot="start" value="RAINDROPS"/>
                        </IonItem>

                        <IonItem>
                            <IonLabel>Soil moisture</IonLabel>
                            <IonRadio slot="start" value="SOIL_MOISTURE"/>
                        </IonItem>
                    </IonRadioGroup>
                </IonList>
            </IonContent>
        </IonPage>
    );
};

export default Type;
