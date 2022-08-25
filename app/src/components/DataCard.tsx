import React from "react";
import {
    IonCard,
    IonCardHeader,
    IonCardSubtitle,
    IonCardTitle,
    IonItem, IonLabel, IonText,
} from "@ionic/react";
import Data from "../model/Data";

interface DataCardProps {
    data: Data;
    title?: string;
}

const DataCard: React.FC<DataCardProps> = ({data, title}) => {
    return (
        <IonCard>
            <IonCardHeader color="warning">
                {
                    title != null ?
                        <>
                            <IonCardSubtitle>{data.getTimestamp()}</IonCardSubtitle>
                            <IonCardTitle>{title}</IonCardTitle></> :
                        <IonCardTitle>{data.getTimestamp()}</IonCardTitle>
                }
            </IonCardHeader>

            <IonItem>
                <IonLabel>Value</IonLabel>
                <IonText>{data.getValueWithSymbol()}</IonText>
            </IonItem>
            <IonItem>
                <IonLabel>Category</IonLabel>
                <IonText>{data.category}</IonText>
            </IonItem>
        </IonCard>
    );
};

export default DataCard;