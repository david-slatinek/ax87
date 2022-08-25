import {IonContent, IonItem, IonLabel, IonList, IonPage, IonRadio, IonRadioGroup} from '@ionic/react';
import React, {useState} from "react";
import Header from "../components/Header";

const Type: React.FC = () => {
    const [selected, setSelected] = useState<string>('CARBON_MONOXIDE');

    return (
        <IonPage>
            <Header title="Type"/>
            <IonContent fullscreen color="light">
                <IonList>
                    <IonRadioGroup value={selected} onIonChange={e => {
                        setSelected(e.detail.value);
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
