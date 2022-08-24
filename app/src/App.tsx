import {Redirect, Route} from 'react-router-dom';
import {
    IonApp,
    IonIcon,
    IonLabel,
    IonRouterOutlet,
    IonTabBar,
    IonTabButton,
    IonTabs,
    setupIonicReact
} from '@ionic/react';
import {IonReactRouter} from '@ionic/react-router';
import {code, home, list} from 'ionicons/icons';
import Home from './pages/Home';
import Tab2 from './pages/Tab2';

/* Core CSS required for Ionic components to work properly */
import '@ionic/react/css/core.css';

/* Basic CSS for apps built with Ionic */
import '@ionic/react/css/normalize.css';
import '@ionic/react/css/structure.css';
import '@ionic/react/css/typography.css';

/* Theme variables */
import './theme/variables.css';
import React from "react";
import Tab3 from "./pages/Tab3";

setupIonicReact();

const App: React.FC = () => (
    <IonApp>
        <IonReactRouter>
            <IonTabs>

                <IonRouterOutlet>
                    <Route exact path="/home">
                        <Home/>
                    </Route>

                    <Route exact path="/tab2">
                        <Tab2/>
                    </Route>

                    <Route exact path="/tab3">
                        <Tab3/>
                    </Route>

                    <Route exact path="/">
                        <Redirect to="/tab1"/>
                    </Route>
                </IonRouterOutlet>

                <IonTabBar slot="bottom">
                    <IonTabButton tab="home" href="/home">
                        <IonIcon icon={home}/>
                        <IonLabel>Home</IonLabel>
                    </IonTabButton>

                    <IonTabButton tab="tab2" href="/tab2">
                        <IonIcon icon={list}/>
                        <IonLabel>View</IonLabel>
                    </IonTabButton>

                    <IonTabButton tab="tab3" href="/tab3">
                        <IonIcon icon={code}/>
                        <IonLabel>Type</IonLabel>
                    </IonTabButton>

                </IonTabBar>
            </IonTabs>
        </IonReactRouter>
    </IonApp>
);

export default App;
